package handlers

import (
	"compress/gzip"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"reviewmysocialnetworks/internal/config"
	"strings"
	"sync"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("[%s] %s %d (%v)", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
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
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self' https://www.instagram.com https://api.instagram.com")
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

		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		writer, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer writer.Close()
		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, writer: writer}, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(body []byte) (int, error) {
	w.Header().Del("Content-Length")
	return w.writer.Write(body)
}

func (w *gzipResponseWriter) Flush() {
	_ = w.writer.Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *gzipResponseWriter) ReadFrom(reader io.Reader) (int64, error) {
	w.Header().Del("Content-Length")
	return io.Copy(w.writer, reader)
}

func acceptsGzip(value string) bool {
	for _, encoding := range strings.Split(value, ",") {
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
				log.Printf("PANIC recovered in %s %s: %v", r.Method, r.URL.Path, err)
				http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
