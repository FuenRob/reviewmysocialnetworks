package handlers

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"reviewmysocialnetworks/internal/config"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		h = middleware(h)
	}
	return h
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := validRequestID(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID, _ = randomToken(12)
		}
		w.Header().Set("X-Request-ID", requestID)
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, requestID))

		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		httpMetrics.inFlight.Add(1)
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		httpMetrics.inFlight.Add(-1)
		httpMetrics.requests.Add(1)
		httpMetrics.durationNanos.Add(uint64(duration))
		if rw.statusCode >= http.StatusBadRequest {
			httpMetrics.errors.Add(1)
		}
		slog.Info("http_request",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"bytes", rw.bytesWritten,
			"duration_ms", float64(duration.Microseconds())/1000,
		)
	})
}

type requestIDKey struct{}

func RequestID(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func validRequestID(value string) string {
	if len(value) < 8 || len(value) > 64 {
		return ""
	}
	for _, ch := range value {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && ch != '-' && ch != '_' {
			return ""
		}
	}
	return value
}

var httpMetrics struct {
	requests      atomic.Uint64
	errors        atomic.Uint64
	durationNanos atomic.Uint64
	inFlight      atomic.Int64
}

type HTTPMetrics struct {
	Requests        uint64
	Errors          uint64
	DurationSeconds float64
	InFlight        int64
}

func HTTPMetricsSnapshot() HTTPMetrics {
	return HTTPMetrics{
		Requests:        httpMetrics.requests.Load(),
		Errors:          httpMetrics.errors.Load(),
		DurationSeconds: float64(httpMetrics.durationNanos.Load()) / float64(time.Second),
		InFlight:        httpMetrics.inFlight.Load(),
	}
}

func CORS(cfg *config.Config) Middleware {
	allowedOrigin := strings.TrimSuffix(cfg.GetFrontendURL(), "/")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := strings.TrimSuffix(r.Header.Get("Origin"), "/")
			if origin != "" && origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Add("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				if origin != "" && origin != allowedOrigin {
					http.Error(w, "origin not allowed", http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; style-src-attr 'none'; font-src 'self'; img-src 'self' data: https:; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self' https://www.instagram.com https://api.instagram.com https://www.tiktok.com")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
		}
		if requestIsHTTPS(r) {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

type visitor struct {
	window time.Time
	count  int
}

var visitors = struct {
	sync.Mutex
	m map[string]visitor
}{m: make(map[string]visitor)}

func RateLimit(cfg *config.Config) Middleware {
	trustProxy := cfg.IsTrustedProxy()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api/health" {
				next.ServeHTTP(w, r)
				return
			}
			host := clientIP(r, trustProxy)
			now := time.Now()
			visitors.Lock()
			v := visitors.m[host]
			if now.Sub(v.window) >= time.Minute {
				v = visitor{window: now}
			}
			v.count++
			visitors.m[host] = v
			limited := v.count > 60
			if len(visitors.m) > 10000 {
				for key, value := range visitors.m {
					if now.Sub(value.window) > 2*time.Minute {
						delete(visitors.m, key)
					}
				}
			}
			visitors.Unlock()
			if limited {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				http.Error(w, `{"error":"Demasiadas solicitudes"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if value := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); net.ParseIP(value) != nil {
			return value
		}
		if value := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); net.ParseIP(value) != nil {
			return value
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead || r.Method == http.MethodOptions || r.Header.Get("Range") != "" || !acceptsGzip(r.Header.Get("Accept-Encoding")) || isCompressedAsset(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		writer := &gzipResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(writer, r)
		if err := writer.close(); err != nil {
			slog.Error("failed to close gzip writer", "error", err)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer          *gzip.Writer
	statusCode      int
	wroteHeader     bool
	responseStarted bool
}

func (w *gzipResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.statusCode = statusCode
}

func (w *gzipResponseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if !statusAllowsBody(w.statusCode) || len(body) == 0 {
		w.commitHeader()
		return w.ResponseWriter.Write(body)
	}
	if err := w.startGzip(); err != nil {
		return w.ResponseWriter.Write(body)
	}
	return w.writer.Write(body)
}

func (w *gzipResponseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if statusAllowsBody(w.statusCode) {
		if err := w.startGzip(); err != nil {
			return
		}
		_ = w.writer.Flush()
	} else {
		w.commitHeader()
	}
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *gzipResponseWriter) ReadFrom(reader io.Reader) (int64, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if !statusAllowsBody(w.statusCode) {
		w.commitHeader()
		return io.Copy(w.ResponseWriter, reader)
	}
	if err := w.startGzip(); err != nil {
		return io.Copy(w.ResponseWriter, reader)
	}
	return io.Copy(w.writer, reader)
}

func (w *gzipResponseWriter) startGzip() error {
	if w.writer != nil {
		return nil
	}
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Del("Content-Length")
	writer, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
	if err != nil {
		w.Header().Del("Content-Encoding")
		w.commitHeader()
		return err
	}
	w.writer = writer
	w.commitHeader()
	return nil
}

func (w *gzipResponseWriter) commitHeader() {
	if w.responseStarted {
		return
	}
	w.responseStarted = true
	w.ResponseWriter.WriteHeader(w.statusCode)
}

func (w *gzipResponseWriter) close() error {
	if w.writer != nil {
		return w.writer.Close()
	}
	w.commitHeader()
	return nil
}

func statusAllowsBody(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode != http.StatusNoContent && statusCode != http.StatusResetContent && statusCode != http.StatusNotModified
}

func acceptsGzip(value string) bool {
	for encoding := range strings.SplitSeq(value, ",") {
		parts := strings.Split(strings.TrimSpace(encoding), ";")
		if strings.EqualFold(parts[0], "gzip") {
			return len(parts) == 1 || !strings.Contains(strings.Join(parts[1:], ";"), "q=0")
		}
	}
	return false
}

func isCompressedAsset(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".woff", ".woff2", ".zip", ".gz", ".br", ".mp4", ".webm":
		return true
	default:
		return false
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("http_panic", "request_id", RequestID(r.Context()), "method", r.Method, "path", r.URL.Path, "error", fmt.Sprint(err))
				http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	wroteHeader  bool
	bytesWritten int
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	written, err := rw.ResponseWriter.Write(body)
	rw.bytesWritten += written
	return written, err
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
