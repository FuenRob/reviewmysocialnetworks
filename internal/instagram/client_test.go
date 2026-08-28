package instagram

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoRequestRetriesTransientResponses(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{httpClient: &http.Client{Timeout: 2 * time.Second}}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	before := MetricsSnapshot()
	resp, err := client.doRequest(context.Background(), req, true)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || attempts.Load() != 3 {
		t.Fatalf("status=%d attempts=%d", resp.StatusCode, attempts.Load())
	}
	after := MetricsSnapshot()
	if after.Retries-before.Retries != 2 || after.Requests-before.Requests != 3 {
		t.Fatalf("unexpected metrics delta: before=%+v after=%+v", before, after)
	}
}

func TestDoRequestDoesNotRetryClientErrors(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client()}
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, err := client.doRequest(context.Background(), req, true)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if attempts.Load() != 1 {
		t.Fatalf("client error was retried %d times", attempts.Load())
	}
}
