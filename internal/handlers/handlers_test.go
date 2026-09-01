package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reviewmysocialnetworks/internal/analyzer"
	"reviewmysocialnetworks/internal/config"
	"reviewmysocialnetworks/internal/instagram"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHandler_Health(t *testing.T) {
	cfg := config.AppConfig
	h := NewHandler(cfg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", resp["status"])
	}
}

func TestHandler_AnalyzeDemo(t *testing.T) {
	cfg := config.AppConfig
	h := NewHandler(cfg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tiers := []string{"A", "B", "D", "F"}
	for _, tier := range tiers {
		body, _ := json.Marshal(map[string]string{"tier": tier})
		req := httptest.NewRequest("POST", "/api/instagram/analyze/demo", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Tier %s: expected status 200, got %d", tier, w.Code)
		}

		var report analyzer.AccountReport
		if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
			t.Fatalf("Tier %s: failed to decode report: %v", tier, err)
		}

		if string(report.OverallGrade) != tier {
			t.Errorf("Expected grade %s, got %s (score %d)", tier, report.OverallGrade, report.OverallScore)
		}
	}
}

func TestHandler_AnalyzeTikTokDemo(t *testing.T) {
	h := NewHandler(config.AppConfig)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	for _, tier := range []string{"A", "B", "D", "F"} {
		body, _ := json.Marshal(map[string]string{"tier": tier})
		req := httptest.NewRequest(http.MethodPost, "/api/tiktok/analyze/demo", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("tier %s: status %d body %s", tier, w.Code, w.Body.String())
		}
		var report analyzer.AccountReport
		if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
		if report.Platform != "tiktok" || string(report.OverallGrade) != tier || report.TikTokMetrics == nil {
			t.Fatalf("unexpected report: %+v", report)
		}
	}
}

func TestTikTokAuthURLSetsStateAndScopes(t *testing.T) {
	cfg := &config.Config{TikTokClientKey: "key", TikTokClientSecret: "secret", TikTokRedirectURI: "https://app.example/api/tiktok/auth/callback"}
	h := NewHandler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/api/tiktok/auth/url", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()
	h.handleTikTokAuthURL(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "video.list") || !strings.Contains(w.Body.String(), "user.info.stats") {
		t.Fatalf("missing TikTok scopes: %s", w.Body.String())
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != tiktokStateCookie || !cookies[0].HttpOnly || !cookies[0].Secure {
		t.Fatalf("unexpected OAuth cookie: %#v", cookies)
	}
}

func TestTikTokCallbackRejectsMissingState(t *testing.T) {
	h := NewHandler(&config.Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/tiktok/auth/callback?code=test", nil)
	w := httptest.NewRecorder()
	h.handleTikTokAuthCallback(w, req)
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "platform=tiktok") {
		t.Fatalf("unexpected redirect: %d %s", w.Code, w.Header().Get("Location"))
	}
}

func TestAuthCallbackRejectsMissingState(t *testing.T) {
	h := NewHandler(&config.Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/instagram/auth/callback?code=test", nil)
	w := httptest.NewRecorder()
	h.handleAuthCallback(w, req)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}
	if location := w.Header().Get("Location"); location != "/?error=invalid_oauth_state" {
		t.Fatalf("unexpected redirect location %q", location)
	}
}

func TestAnalyzeManualRejectsOversizedBody(t *testing.T) {
	h := NewHandler(&config.Config{})
	largeBody := strings.NewReader(`{"profile":{},"media":[],"padding":"` + strings.Repeat("x", maxJSONBody) + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/instagram/analyze/manual", largeBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.handleAnalyzeManual(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected request entity too large, got %d", w.Code)
	}
}

func TestAnalyzeDemoRequiresJSONContentType(t *testing.T) {
	h := NewHandler(&config.Config{})
	req := httptest.NewRequest(http.MethodPost, "/api/instagram/analyze/demo", strings.NewReader(`{"tier":"A"}`))
	w := httptest.NewRecorder()
	h.handleAnalyzeDemo(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected unsupported media type, got %d", w.Code)
	}
}

func TestCORSRejectsUnknownPreflightOrigin(t *testing.T) {
	cfg := &config.Config{FrontendURL: "https://app.example"}
	handler := CORS(cfg)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	req := httptest.NewRequest(http.MethodOptions, "/api/instagram/analyze/demo", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden, got %d", w.Code)
	}
}

func TestAuthURLSetsPrivateSecureStateCookie(t *testing.T) {
	cfg := &config.Config{InstagramAppID: "app", InstagramAppSecret: "secret", InstagramRedirectURI: "https://app.example/api/instagram/auth/callback"}
	h := NewHandler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/api/instagram/auth/url", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()
	h.handleAuthURL(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected success, got %d", w.Code)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Fatalf("OAuth state cookie is not hardened: %#v", cookies)
	}
	if strings.Contains(w.Body.String(), `"state"`) {
		t.Fatal("OAuth state must not be exposed in the JSON response")
	}
}

func TestValidateManualInputRejectsInvalidValues(t *testing.T) {
	profile := instagram.UserProfile{FollowersCount: -1}
	media := []instagram.MediaItem{{Timestamp: time.Now(), Caption: strings.Repeat("x", 2201)}}
	if err := validateManualInput(profile, media); err == nil {
		t.Fatal("expected semantic validation error")
	}
}

func TestCompressionProducesGzipResponse(t *testing.T) {
	handler := Compression(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "15")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("Content-Encoding") != "gzip" {
		t.Fatal("expected gzip content encoding")
	}
	if value := w.Header().Get("Content-Length"); value != "" {
		t.Fatalf("compressed response retained Content-Length %q", value)
	}
	reader, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatalf("invalid gzip response: %v", err)
	}
	body, err := io.ReadAll(reader)
	if err != nil || string(body) != `{"status":"ok"}` {
		t.Fatalf("unexpected decompressed response %q: %v", body, err)
	}
}

func TestCompressionDoesNotCompressResponsesWithoutBody(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusResetContent, http.StatusNotModified} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			handler := Compression(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			}))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != status {
				t.Fatalf("expected status %d, got %d", status, w.Code)
			}
			if w.Header().Get("Content-Encoding") != "" {
				t.Fatalf("unexpected content encoding %q", w.Header().Get("Content-Encoding"))
			}
			if w.Body.Len() != 0 {
				t.Fatalf("expected empty response body, got %d bytes", w.Body.Len())
			}
		})
	}
}

func TestClientIPTrustsForwardedHeaderOnlyWhenConfigured(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.RemoteAddr = "192.0.2.10:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.5, 192.0.2.10")
	if got := clientIP(req, false); got != "192.0.2.10" {
		t.Fatalf("unexpected direct client IP %q", got)
	}
	if got := clientIP(req, true); got != "198.51.100.5" {
		t.Fatalf("unexpected proxied client IP %q", got)
	}
}

func TestSecurityHeadersDoNotAllowExternalFonts(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	csp := w.Header().Get("Content-Security-Policy")
	if strings.Contains(csp, "fonts.googleapis.com") || strings.Contains(csp, "fonts.gstatic.com") {
		t.Fatalf("CSP still allows external font providers: %s", csp)
	}
	if !strings.Contains(csp, "object-src 'none'") {
		t.Fatalf("CSP does not disable plugins: %s", csp)
	}
	if !strings.Contains(csp, "style-src-attr 'none'") {
		t.Fatalf("CSP still permits inline style attributes: %s", csp)
	}
}

func TestLoggerAddsRequestIDAndMetrics(t *testing.T) {
	before := HTTPMetricsSnapshot()
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.Header.Set("X-Request-ID", "request_test_123")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if got := w.Header().Get("X-Request-ID"); got != "request_test_123" {
		t.Fatalf("unexpected request ID %q", got)
	}
	after := HTTPMetricsSnapshot()
	if after.Requests-before.Requests != 1 {
		t.Fatalf("request metric was not incremented: before=%+v after=%+v", before, after)
	}
}

func TestMetricsEndpointUsesPrometheusFormat(t *testing.T) {
	h := NewHandler(&config.Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	w := httptest.NewRecorder()
	h.handleMetrics(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "rmsn_instagram_retries_total") {
		t.Fatalf("unexpected metrics response: status=%d body=%q", w.Code, w.Body.String())
	}
}

func TestRateLimitUsesLowerLimitForExternalAnalyses(t *testing.T) {
	handler := RateLimit(&config.Config{})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for requestNumber := 1; requestNumber <= analysisRequestsPerMinute+1; requestNumber++ {
		req := httptest.NewRequest(http.MethodPost, "/api/tiktok/analyze/token", nil)
		req.RemoteAddr = "192.0.2.25:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if requestNumber <= analysisRequestsPerMinute && w.Code != http.StatusNoContent {
			t.Fatalf("request %d was limited early with status %d", requestNumber, w.Code)
		}
		if requestNumber == analysisRequestsPerMinute+1 && w.Code != http.StatusTooManyRequests {
			t.Fatalf("expected request %d to be limited, got status %d", requestNumber, w.Code)
		}
	}
}

func TestAnalysisConcurrencyRejectsExcessWork(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	handler := analysisConcurrency(1)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusNoContent)
	}))

	firstRequest := httptest.NewRequest(http.MethodPost, "/api/instagram/analyze/token", nil)
	firstResponse := httptest.NewRecorder()
	go func() {
		defer close(done)
		handler.ServeHTTP(firstResponse, firstRequest)
	}()
	<-started

	secondRequest := httptest.NewRequest(http.MethodPost, "/api/tiktok/analyze/token", nil)
	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, secondRequest)
	if secondResponse.Code != http.StatusTooManyRequests {
		t.Fatalf("expected excess analysis to be rejected, got %d", secondResponse.Code)
	}
	if secondResponse.Header().Get("Retry-After") == "" {
		t.Fatal("limited response did not include Retry-After")
	}

	close(release)
	<-done
	if firstResponse.Code != http.StatusNoContent {
		t.Fatalf("first analysis did not complete: %d", firstResponse.Code)
	}
}
