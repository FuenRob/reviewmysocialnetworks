package tiktok

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	OAuthAuthorizeURL = "https://www.tiktok.com/v2/auth/authorize/"
	OAuthTokenURL     = "https://open.tiktokapis.com/v2/oauth/token/"
	APIBaseURL        = "https://open.tiktokapis.com"
	maxResponseBytes  = 8 << 20
)

type Client struct {
	httpClient   *http.Client
	clientKey    string
	clientSecret string
	apiBaseURL   string
	tokenURL     string
}

func NewClient(clientKey, clientSecret string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		clientKey:  clientKey, clientSecret: clientSecret,
		apiBaseURL: APIBaseURL, tokenURL: OAuthTokenURL,
	}
}

func (c *Client) GetAuthorizationURL(redirectURI, state string) string {
	params := url.Values{}
	params.Set("client_key", c.clientKey)
	params.Set("response_type", "code")
	params.Set("scope", "user.info.basic,user.info.profile,user.info.stats,video.list")
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)
	return OAuthAuthorizeURL + "?" + params.Encode()
}

func (c *Client) ExchangeCodeForToken(ctx context.Context, code, redirectURI string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("client_key", c.clientKey)
	form.Set("client_secret", c.clientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", redirectURI)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var token TokenResponse
	if err := c.doJSON(req, &token, false); err != nil {
		return nil, fmt.Errorf("tiktok token exchange: %w", err)
	}
	if token.AccessToken == "" {
		return nil, errors.New("tiktok returned an empty access token")
	}
	return &token, nil
}

func (c *Client) FetchProfile(ctx context.Context, accessToken string) (*UserProfile, error) {
	fields := "open_id,union_id,avatar_url,display_name,bio_description,profile_deep_link,is_verified,username,follower_count,following_count,likes_count,video_count"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiBaseURL+"/v2/user/info/?fields="+fields, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	var response struct {
		Data struct {
			User UserProfile `json:"user"`
		} `json:"data"`
		Error APIError `json:"error"`
	}
	if err := c.doJSON(req, &response, true); err != nil {
		return nil, fmt.Errorf("tiktok profile: %w", err)
	}
	if response.Error.Code != "" && response.Error.Code != "ok" {
		return nil, fmt.Errorf("tiktok profile: %s", response.Error.Code)
	}
	if response.Data.User.OpenID == "" {
		return nil, errors.New("tiktok profile response did not include open_id")
	}
	return &response.Data.User, nil
}

func (c *Client) FetchVideos(ctx context.Context, accessToken string, limit int) ([]Video, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	fields := "id,create_time,cover_image_url,share_url,video_description,duration,title,like_count,comment_count,share_count,view_count,is_aigc"
	videos := make([]Video, 0, limit)
	var cursor int64
	for len(videos) < limit {
		pageSize := min(limit-len(videos), 20)
		body := struct {
			Cursor   int64 `json:"cursor,omitempty"`
			MaxCount int   `json:"max_count"`
		}{Cursor: cursor, MaxCount: pageSize}
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiBaseURL+"/v2/video/list/?fields="+fields, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")
		var response struct {
			Data struct {
				Videos  []Video `json:"videos"`
				Cursor  int64   `json:"cursor"`
				HasMore bool    `json:"has_more"`
			} `json:"data"`
			Error APIError `json:"error"`
		}
		if err := c.doJSON(req, &response, true); err != nil {
			return nil, fmt.Errorf("tiktok videos: %w", err)
		}
		if response.Error.Code != "" && response.Error.Code != "ok" {
			return nil, fmt.Errorf("tiktok videos: %s", response.Error.Code)
		}
		videos = append(videos, response.Data.Videos...)
		if !response.Data.HasMore || len(response.Data.Videos) == 0 {
			break
		}
		cursor = response.Data.Cursor
	}
	if len(videos) > limit {
		videos = videos[:limit]
	}
	return videos, nil
}

func (c *Client) doJSON(req *http.Request, target any, retryable bool) error {
	attempts := 1
	if retryable {
		attempts = 3
	}
	for attempt := 0; attempt < attempts; attempt++ {
		request := req.Clone(req.Context())
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return err
			}
			request.Body = body
		}
		resp, err := c.httpClient.Do(request)
		if err != nil {
			if retryable && attempt+1 < attempts {
				if err := waitForRetry(req.Context(), attempt); err != nil {
					return err
				}
				continue
			}
			return err
		}
		limited := io.LimitReader(resp.Body, maxResponseBytes+1)
		body, readErr := io.ReadAll(limited)
		_ = resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		if len(body) > maxResponseBytes {
			return errors.New("tiktok response exceeded size limit")
		}
		if retryable && attempt+1 < attempts && isRetryableTikTokStatus(resp.StatusCode) {
			if err := waitForRetry(req.Context(), attempt); err != nil {
				return err
			}
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("status %d", resp.StatusCode)
		}
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("invalid JSON response: %w", err)
		}
		return nil
	}
	return errors.New("tiktok request retries exhausted")
}

func isRetryableTikTokStatus(status int) bool {
	return status == http.StatusTooManyRequests || status == http.StatusBadGateway || status == http.StatusServiceUnavailable || status == http.StatusGatewayTimeout
}

func waitForRetry(ctx context.Context, attempt int) error {
	timer := time.NewTimer(time.Duration(100*(1<<attempt)) * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
