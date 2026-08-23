package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reviewmysocialnetworks/internal/analyzer"
	"reviewmysocialnetworks/internal/config"
	"reviewmysocialnetworks/internal/instagram"
	"strings"
	"time"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", h.handleHealth)
	mux.HandleFunc("GET /api/auth/url", h.handleAuthURL)
	mux.HandleFunc("GET /api/auth/callback", h.handleAuthCallback)
	mux.HandleFunc("POST /api/analyze/token", h.handleAnalyzeToken)
	mux.HandleFunc("POST /api/analyze/demo", h.handleAnalyzeDemo)
	mux.HandleFunc("POST /api/analyze/manual", h.handleAnalyzeManual)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"service":   "ReviewMySocialNetworks API (Instagram Direct)",
	})
}

func (h *Handler) handleAuthURL(w http.ResponseWriter, r *http.Request) {
	appID, appSecret, redirectURI, _ := h.cfg.Get()
	if appID == "" || appSecret == "" {
		respondError(w, http.StatusBadRequest, "Instagram App ID y Secret deben estar configurados en el archivo .env")
		return
	}

	stateBytes := make([]byte, 16)
	_, _ = rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	client := instagram.NewClient(appID, appSecret)
	mode := r.URL.Query().Get("mode")

	var authURL string
	if mode == "basic" {
		authURL = client.GetBasicAuthorizationURL(redirectURI, state)
	} else {
		authURL = client.GetAuthorizationURL(redirectURI, state)
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorParam := r.URL.Query().Get("error")
		errorReason := r.URL.Query().Get("error_reason")
		errorDesc := r.URL.Query().Get("error_description")
		http.Redirect(w, r, fmt.Sprintf("/?error=%s&desc=%s", errorParam, errorDesc), http.StatusTemporaryRedirect)
		_ = errorReason
		return
	}

	appID, appSecret, redirectURI, _ := h.cfg.Get()
	client := instagram.NewClient(appID, appSecret)

	code = strings.TrimSuffix(code, "#_")

	tokenResp, err := client.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/?error=token_exchange_failed&desc=%s", err.Error()), http.StatusTemporaryRedirect)
		return
	}

	longLivedToken, err := client.ExchangeForLongLivedToken(tokenResp.AccessToken)
	tokenToUse := tokenResp.AccessToken
	if err == nil && longLivedToken != nil && longLivedToken.AccessToken != "" {
		tokenToUse = longLivedToken.AccessToken
	}

	userID := tokenResp.GetUserID()
	redirectTarget := fmt.Sprintf("/?access_token=%s", url.QueryEscape(tokenToUse))
	if userID != "" {
		redirectTarget += fmt.Sprintf("&user_id=%s", url.QueryEscape(userID))
	}
	redirectTarget += "#auth_success"
	http.Redirect(w, r, redirectTarget, http.StatusTemporaryRedirect)
}

func (h *Handler) handleAnalyzeToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"access_token"`
		AccountID   string `json:"account_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.AccessToken) == "" {
		respondError(w, http.StatusBadRequest, "El parámetro access_token es obligatorio")
		return
	}

	appID, appSecret, _, _ := h.cfg.Get()
	client := instagram.NewClient(appID, appSecret)

	accountID := strings.TrimSpace(req.AccountID)
	if accountID == "" {
		accountID = "me"
	}

	profile, err := client.FetchProfile(accountID, req.AccessToken)
	if err != nil {
		respondError(w, http.StatusBadGateway, fmt.Sprintf("Error obteniendo perfil de Instagram: %v", err))
		return
	}

	mediaList, err := client.FetchMediaList(accountID, req.AccessToken, 50)
	if err != nil {
		mediaList = []instagram.MediaItem{}
	}

	report := analyzer.AnalyzeAccount(profile, mediaList)

	respondJSON(w, http.StatusOK, report)
}

func (h *Handler) handleAnalyzeDemo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tier string `json:"tier"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Tier == "" {
		req.Tier = "A"
	}

	profile, mediaList := instagram.GetMockAccount(req.Tier)
	report := analyzer.AnalyzeAccount(profile, mediaList)

	respondJSON(w, http.StatusOK, report)
}

func (h *Handler) handleAnalyzeManual(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Profile instagram.UserProfile `json:"profile"`
		Media   []instagram.MediaItem `json:"media"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Datos de perfil y media inválidos")
		return
	}

	report := analyzer.AnalyzeAccount(&req.Profile, req.Media)
	respondJSON(w, http.StatusOK, report)
}

func SPAHandler(distDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		path := filepath.Join(distDir, filepath.Clean(r.URL.Path))
		info, err := os.Stat(path)

		if errors.Is(err, os.ErrNotExist) || info.IsDir() {
			indexPath := filepath.Join(distDir, "index.html")
			http.ServeFile(w, r, indexPath)
			return
		}

		http.ServeFile(w, r, path)
	}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}
