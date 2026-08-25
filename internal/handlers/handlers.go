package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reviewmysocialnetworks/internal/analyzer"
	"reviewmysocialnetworks/internal/config"
	"reviewmysocialnetworks/internal/instagram"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	maxJSONBody       = 1 << 20
	maxManualPosts    = 100
	oauthStateCookie  = "rmsn_oauth_state"
	authSessionCookie = "rmsn_auth_session"
)

type authSession struct {
	accessToken string
	userID      string
	expiresAt   time.Time
}

type Handler struct {
	cfg      *config.Config
	sessions sync.Map
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", h.handleHealth)
	mux.HandleFunc("GET /api/auth/url", h.handleAuthURL)
	mux.HandleFunc("GET /api/auth/callback", h.handleAuthCallback)
	mux.HandleFunc("GET /api/auth/result", h.handleAuthResult)
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
	now := time.Now()
	h.sessions.Range(func(key, value any) bool {
		if session, ok := value.(authSession); ok && now.After(session.expiresAt) {
			h.sessions.Delete(key)
		}
		return true
	})

	appID, appSecret, redirectURI, _ := h.cfg.Get()
	if appID == "" || appSecret == "" {
		respondError(w, http.StatusBadRequest, "Instagram App ID y Secret deben estar configurados en el archivo .env")
		return
	}

	state, err := randomToken(32)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo iniciar la autenticación")
		return
	}
	setPrivateCookie(w, r, oauthStateCookie, state, 10*time.Minute)

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
	})
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie(oauthStateCookie)
	state := r.URL.Query().Get("state")
	if err != nil || state == "" || subtle.ConstantTimeCompare([]byte(stateCookie.Value), []byte(state)) != 1 {
		clearCookie(w, r, oauthStateCookie)
		redirectWithError(w, r, "invalid_oauth_state")
		return
	}
	clearCookie(w, r, oauthStateCookie)

	code := r.URL.Query().Get("code")
	if code == "" {
		errorParam := r.URL.Query().Get("error")
		errorReason := r.URL.Query().Get("error_reason")
		errorDesc := r.URL.Query().Get("error_description")
		redirectWithError(w, r, firstNonEmpty(errorParam, errorDesc, "oauth_denied"))
		_ = errorReason
		return
	}

	appID, appSecret, redirectURI, _ := h.cfg.Get()
	client := instagram.NewClient(appID, appSecret)

	code = strings.TrimSuffix(code, "#_")

	tokenResp, err := client.ExchangeCodeForToken(r.Context(), code, redirectURI)
	if err != nil {
		redirectWithError(w, r, "token_exchange_failed")
		return
	}

	longLivedToken, err := client.ExchangeForLongLivedToken(r.Context(), tokenResp.AccessToken)
	tokenToUse := tokenResp.AccessToken
	if err == nil && longLivedToken != nil && longLivedToken.AccessToken != "" {
		tokenToUse = longLivedToken.AccessToken
	}

	userID := tokenResp.GetUserID()
	sessionID, err := randomToken(32)
	if err != nil {
		redirectWithError(w, r, "session_creation_failed")
		return
	}
	h.sessions.Store(sessionID, authSession{accessToken: tokenToUse, userID: userID, expiresAt: time.Now().Add(5 * time.Minute)})
	setPrivateCookie(w, r, authSessionCookie, sessionID, 5*time.Minute)
	http.Redirect(w, r, "/#auth_success", http.StatusSeeOther)
}

func (h *Handler) handleAuthResult(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authSessionCookie)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Sesión de autenticación no disponible")
		return
	}
	clearCookie(w, r, authSessionCookie)
	value, ok := h.sessions.LoadAndDelete(cookie.Value)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Sesión de autenticación inválida o ya utilizada")
		return
	}
	session := value.(authSession)
	if time.Now().After(session.expiresAt) {
		respondError(w, http.StatusUnauthorized, "Sesión de autenticación caducada")
		return
	}
	h.analyzeToken(w, r, session.accessToken, session.userID)
}

func (h *Handler) handleAnalyzeToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"access_token"`
		AccountID   string `json:"account_id,omitempty"`
	}

	if err := decodeJSON(w, r, &req); err != nil || strings.TrimSpace(req.AccessToken) == "" {
		respondError(w, http.StatusBadRequest, "El parámetro access_token es obligatorio")
		return
	}

	accountID := strings.TrimSpace(req.AccountID)
	if accountID == "" {
		accountID = "me"
	}

	if accountID != "me" && !isDigits(accountID) {
		respondError(w, http.StatusBadRequest, "account_id no es válido")
		return
	}
	h.analyzeToken(w, r, strings.TrimSpace(req.AccessToken), accountID)
}

func (h *Handler) analyzeToken(w http.ResponseWriter, r *http.Request, accessToken, accountID string) {
	appID, appSecret, _, _ := h.cfg.Get()
	client := instagram.NewClient(appID, appSecret)
	if accountID == "" {
		accountID = "me"
	}
	profile, err := client.FetchProfile(r.Context(), accountID, accessToken)
	if err != nil {
		respondError(w, http.StatusBadGateway, "No se pudo obtener el perfil de Instagram")
		return
	}

	mediaList, err := client.FetchMediaList(r.Context(), accountID, accessToken, 50)
	if err != nil {
		respondError(w, http.StatusBadGateway, "No se pudieron obtener las publicaciones de Instagram")
		return
	}

	report := analyzer.AnalyzeAccount(profile, mediaList)

	respondJSON(w, http.StatusOK, report)
}

func (h *Handler) handleAnalyzeDemo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tier string `json:"tier"`
	}

	if err := decodeJSON(w, r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud demo inválida")
		return
	}
	if req.Tier != "A" && req.Tier != "B" && req.Tier != "D" && req.Tier != "F" {
		respondError(w, http.StatusBadRequest, "El nivel demo debe ser A, B, D o F")
		return
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

	if err := decodeJSON(w, r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Datos de perfil y media inválidos")
		return
	}
	if len(req.Media) > maxManualPosts {
		respondError(w, http.StatusRequestEntityTooLarge, "Demasiadas publicaciones; el máximo es 100")
		return
	}
	if err := validateManualInput(req.Profile, req.Media); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	report := analyzer.AnalyzeAccount(&req.Profile, req.Media)
	respondJSON(w, http.StatusOK, report)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("solo se permite un objeto JSON")
	}
	return nil
}

func randomToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func setPrivateCookie(w http.ResponseWriter, r *http.Request, name, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/api/auth", HttpOnly: true, Secure: requestIsHTTPS(r), SameSite: http.SameSiteLaxMode, MaxAge: int(ttl.Seconds())})
}

func clearCookie(w http.ResponseWriter, r *http.Request, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Path: "/api/auth", HttpOnly: true, Secure: requestIsHTTPS(r), SameSite: http.SameSiteLaxMode, MaxAge: -1})
}

func requestIsHTTPS(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func redirectWithError(w http.ResponseWriter, r *http.Request, code string) {
	params := url.Values{"error": {code}}
	http.Redirect(w, r, "/?"+params.Encode(), http.StatusSeeOther)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "oauth_error"
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func validateManualInput(profile instagram.UserProfile, media []instagram.MediaItem) error {
	if err := maxLengths(map[string]struct {
		value string
		max   int
	}{
		"profile.id":           {profile.ID, 128},
		"profile.username":     {profile.Username, 64},
		"profile.name":         {profile.Name, 256},
		"profile.biography":    {profile.Biography, 2200},
		"profile.account_type": {profile.AccountType, 64},
	}); err != nil {
		return err
	}
	if err := validateURLField("profile.profile_picture_url", profile.ProfilePictureURL); err != nil {
		return err
	}
	if err := validateURLField("profile.website", profile.Website); err != nil {
		return err
	}
	if !validCount(profile.FollowersCount) || !validCount(profile.FollowsCount) || !validCount(profile.MediaCount) {
		return errors.New("los contadores del perfil deben estar entre 0 y 1.000.000.000")
	}

	for i, item := range media {
		prefix := fmt.Sprintf("media[%d]", i)
		if err := maxLengths(map[string]struct {
			value string
			max   int
		}{
			prefix + ".id":                 {item.ID, 128},
			prefix + ".caption":            {item.Caption, 2200},
			prefix + ".media_type":         {item.MediaType, 32},
			prefix + ".media_product_type": {item.MediaProductType, 32},
		}); err != nil {
			return err
		}
		for name, value := range map[string]string{
			prefix + ".media_url":     item.MediaURL,
			prefix + ".thumbnail_url": item.ThumbnailURL,
			prefix + ".permalink":     item.Permalink,
		} {
			if err := validateURLField(name, value); err != nil {
				return err
			}
		}
		if !validCount(item.LikeCount) || !validCount(item.CommentsCount) {
			return fmt.Errorf("%s contiene contadores fuera del rango permitido", prefix)
		}
		if item.Timestamp.IsZero() || item.Timestamp.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) || item.Timestamp.After(time.Now().Add(24*time.Hour)) {
			return fmt.Errorf("%s.timestamp está fuera del rango permitido", prefix)
		}
		if len(item.Children) > 20 {
			return fmt.Errorf("%s contiene más de 20 elementos hijos", prefix)
		}
		for childIndex, child := range item.Children {
			childPrefix := fmt.Sprintf("%s.children[%d]", prefix, childIndex)
			if err := maxLengths(map[string]struct {
				value string
				max   int
			}{
				childPrefix + ".id":         {child.ID, 128},
				childPrefix + ".media_type": {child.MediaType, 32},
			}); err != nil {
				return err
			}
			if err := validateURLField(childPrefix+".media_url", child.MediaURL); err != nil {
				return err
			}
		}
		if item.Insights != nil && (!validCount(item.Insights.Impressions) || !validCount(item.Insights.Reach) || !validCount(item.Insights.Saved) || !validCount(item.Insights.Engagement) || !validCount(item.Insights.VideoViews) || !validCount(item.Insights.Shares)) {
			return fmt.Errorf("%s.insights contiene contadores fuera del rango permitido", prefix)
		}
	}
	return nil
}

func maxLengths(fields map[string]struct {
	value string
	max   int
}) error {
	for name, field := range fields {
		if utf8.RuneCountInString(field.value) > field.max {
			return fmt.Errorf("%s supera el máximo de %d caracteres", name, field.max)
		}
	}
	return nil
}

func validateURLField(name, value string) error {
	if value == "" {
		return nil
	}
	if utf8.RuneCountInString(value) > 2048 {
		return fmt.Errorf("%s supera el máximo de 2048 caracteres", name)
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("%s debe ser una URL http:// o https:// válida", name)
	}
	return nil
}

func validCount(value int) bool {
	return value >= 0 && value <= 1_000_000_000
}

func SPAHandler(distDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(filepath.Clean(r.URL.Path), string(filepath.Separator))
		path := filepath.Join(distDir, cleanPath)
		rel, relErr := filepath.Rel(distDir, path)
		if relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			http.NotFound(w, r)
			return
		}
		info, err := os.Stat(path)

		if errors.Is(err, os.ErrNotExist) || info.IsDir() {
			indexPath := filepath.Join(distDir, "index.html")
			http.ServeFile(w, r, indexPath)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
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
