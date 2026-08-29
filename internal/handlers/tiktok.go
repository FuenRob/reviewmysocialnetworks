package handlers

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"net/url"
	"reviewmysocialnetworks/internal/analyzer"
	"reviewmysocialnetworks/internal/tiktok"
	"strings"
	"time"
)

func (h *Handler) handleTikTokAuthURL(w http.ResponseWriter, r *http.Request) {
	clientKey, clientSecret, redirectURI := h.cfg.GetTikTok()
	if clientKey == "" || clientSecret == "" || redirectURI == "" {
		respondError(w, http.StatusBadRequest, "TikTok Client Key, Client Secret y Redirect URI deben estar configurados")
		return
	}
	state, err := randomToken(32)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo iniciar la autenticación de TikTok")
		return
	}
	setPrivateCookie(w, r, tiktokStateCookie, state, 10*time.Minute)
	client := tiktok.NewClient(clientKey, clientSecret)
	respondJSON(w, http.StatusOK, map[string]string{"auth_url": client.GetAuthorizationURL(redirectURI, state)})
}

func (h *Handler) handleTikTokAuthCallback(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(tiktokStateCookie)
	state := r.URL.Query().Get("state")
	if err != nil || state == "" || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state)) != 1 {
		clearCookie(w, r, tiktokStateCookie)
		redirectTikTokError(w, r, "invalid_oauth_state")
		return
	}
	clearCookie(w, r, tiktokStateCookie)
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		redirectTikTokError(w, r, firstNonEmpty(r.URL.Query().Get("error"), r.URL.Query().Get("error_description"), "oauth_denied"))
		return
	}

	clientKey, clientSecret, redirectURI := h.cfg.GetTikTok()
	client := tiktok.NewClient(clientKey, clientSecret)
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	token, err := client.ExchangeCodeForToken(ctx, code, redirectURI)
	if err != nil {
		redirectTikTokError(w, r, "token_exchange_failed")
		return
	}
	sessionID, err := randomToken(32)
	if err != nil {
		redirectTikTokError(w, r, "session_creation_failed")
		return
	}
	h.sessions.Store(sessionID, authSession{accessToken: token.AccessToken, userID: token.OpenID, expiresAt: time.Now().Add(5 * time.Minute)})
	setPrivateCookie(w, r, tiktokSessionCookie, sessionID, 5*time.Minute)
	http.Redirect(w, r, "/#tiktok_auth_success", http.StatusSeeOther)
}

func (h *Handler) handleTikTokAuthResult(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(tiktokSessionCookie)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Sesión de TikTok no disponible")
		return
	}
	clearCookie(w, r, tiktokSessionCookie)
	value, ok := h.sessions.LoadAndDelete(cookie.Value)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Sesión de TikTok inválida o ya utilizada")
		return
	}
	session := value.(authSession)
	if time.Now().After(session.expiresAt) {
		respondError(w, http.StatusUnauthorized, "Sesión de TikTok caducada")
		return
	}
	h.analyzeTikTokToken(w, r, session.accessToken)
}

func (h *Handler) handleTikTokAnalyzeToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"access_token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		respondDecodeError(w, err, "Solicitud de token de TikTok inválida")
		return
	}
	if strings.TrimSpace(req.AccessToken) == "" {
		respondError(w, http.StatusBadRequest, "El parámetro access_token es obligatorio")
		return
	}
	h.analyzeTikTokToken(w, r, strings.TrimSpace(req.AccessToken))
}

func (h *Handler) analyzeTikTokToken(w http.ResponseWriter, r *http.Request, accessToken string) {
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	clientKey, clientSecret, _ := h.cfg.GetTikTok()
	client := tiktok.NewClient(clientKey, clientSecret)
	profile, err := client.FetchProfile(ctx, accessToken)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "TikTok tardó demasiado en responder")
			return
		}
		respondError(w, http.StatusBadGateway, "No se pudo obtener el perfil de TikTok; revisa el token y los scopes autorizados")
		return
	}
	videos, err := client.FetchVideos(ctx, accessToken, 50)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "TikTok tardó demasiado en responder")
			return
		}
		respondError(w, http.StatusBadGateway, "No se pudieron obtener los vídeos públicos de TikTok")
		return
	}
	respondJSON(w, http.StatusOK, analyzer.AnalyzeTikTokAccount(profile, videos))
}

func (h *Handler) handleTikTokAnalyzeDemo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tier string `json:"tier"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		respondDecodeError(w, err, "Solicitud demo de TikTok inválida")
		return
	}
	if req.Tier != "A" && req.Tier != "B" && req.Tier != "D" && req.Tier != "F" {
		respondError(w, http.StatusBadRequest, "El nivel demo debe ser A, B, D o F")
		return
	}
	profile, videos := tiktok.GetMockAccount(req.Tier)
	respondJSON(w, http.StatusOK, analyzer.AnalyzeTikTokAccount(profile, videos))
}

func redirectTikTokError(w http.ResponseWriter, r *http.Request, code string) {
	params := url.Values{"error": {code}, "platform": {"tiktok"}}
	http.Redirect(w, r, "/?"+params.Encode(), http.StatusSeeOther)
}
