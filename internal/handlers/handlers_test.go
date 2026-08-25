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
		req := httptest.NewRequest("POST", "/api/analyze/demo", bytes.NewReader(body))
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

func TestAuthCallbackRejectsMissingState(t *testing.T) {
	h := NewHandler(&config.Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback?code=test", nil)
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
	req := httptest.NewRequest(http.MethodPost, "/api/analyze/manual", largeBody)
	w := httptest.NewRecorder()
	h.handleAnalyzeManual(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", w.Code)
	}
}

func TestCORSRejectsUnknownPreflightOrigin(t *testing.T) {
	cfg := &config.Config{FrontendURL: "https://app.example"}
	handler := CORS(cfg)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	req := httptest.NewRequest(http.MethodOptions, "/api/analyze/demo", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden, got %d", w.Code)
	}
}

func TestAuthURLSetsPrivateSecureStateCookie(t *testing.T) {
	cfg := &config.Config{InstagramAppID: "app", InstagramAppSecret: "secret", InstagramRedirectURI: "https://app.example/api/auth/callback"}
	h := NewHandler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/api/auth/url", nil)
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
