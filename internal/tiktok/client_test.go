package tiktok

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
)

func TestAuthorizationURLIncludesRequiredScopes(t *testing.T) {
	client := NewClient("client-key", "secret")
	parsed, err := url.Parse(client.GetAuthorizationURL("https://app.example/tiktok/callback", "state-123"))
	if err != nil {
		t.Fatal(err)
	}
	for _, scope := range []string{"user.info.basic", "user.info.profile", "user.info.stats", "video.list"} {
		if !strings.Contains(parsed.Query().Get("scope"), scope) {
			t.Fatalf("missing scope %s", scope)
		}
	}
	if parsed.Query().Get("state") != "state-123" {
		t.Fatal("missing OAuth state")
	}
}

func TestExchangeCodeForToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Form.Get("client_key") != "key" || r.Form.Get("client_secret") != "secret" || r.Form.Get("code") != "code" || r.Form.Get("redirect_uri") != "https://app.example/callback" {
			t.Fatalf("unexpected token form: %v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "access", "open_id": "creator", "expires_in": 86400})
	}))
	defer server.Close()
	client := NewClient("key", "secret")
	client.tokenURL = server.URL
	client.httpClient = server.Client()
	token, err := client.ExchangeCodeForToken(context.Background(), "code", "https://app.example/callback")
	if err != nil || token.AccessToken != "access" || token.OpenID != "creator" {
		t.Fatalf("token=%+v err=%v", token, err)
	}
}

func TestFetchProfileRetriesTransientFailure(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"user": map[string]any{"open_id": "id", "username": "creator"}}, "error": map[string]any{"code": "ok"}})
	}))
	defer server.Close()
	client := NewClient("key", "secret")
	client.apiBaseURL = server.URL
	client.httpClient = server.Client()
	if _, err := client.FetchProfile(context.Background(), "token"); err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("attempts=%d, want 2", attempts.Load())
	}
}

func TestFetchProfileAndPaginatedVideos(t *testing.T) {
	videoCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("missing bearer token")
		}
		switch r.URL.Path {
		case "/v2/user/info/":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"user": map[string]any{"open_id": "id", "username": "creator", "follower_count": 1000}}, "error": map[string]any{"code": "ok"}})
		case "/v2/video/list/":
			videoCalls++
			hasMore := videoCalls == 1
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"videos": []map[string]any{{"id": strings.Repeat("v", videoCalls), "view_count": 100}}, "cursor": 123, "has_more": hasMore}, "error": map[string]any{"code": "ok"}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := NewClient("key", "secret")
	client.apiBaseURL = server.URL
	client.httpClient = server.Client()
	profile, err := client.FetchProfile(context.Background(), "token")
	if err != nil || profile.Username != "creator" {
		t.Fatalf("profile=%+v err=%v", profile, err)
	}
	videos, err := client.FetchVideos(context.Background(), "token", 30)
	if err != nil || len(videos) != 2 || videoCalls != 2 {
		t.Fatalf("videos=%+v calls=%d err=%v", videos, videoCalls, err)
	}
}
