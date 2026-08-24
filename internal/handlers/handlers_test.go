package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reviewmysocialnetworks/internal/analyzer"
	"reviewmysocialnetworks/internal/config"
	"strings"
	"testing"
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
