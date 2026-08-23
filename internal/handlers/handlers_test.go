package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reviewmysocialnetworks/internal/analyzer"
	"reviewmysocialnetworks/internal/config"
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
