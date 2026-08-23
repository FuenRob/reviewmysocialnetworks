package instagram

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

type UserProfile struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name,omitempty"`
	Biography         string `json:"biography,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
	FollowersCount    int    `json:"followers_count"`
	FollowsCount      int    `json:"follows_count"`
	MediaCount        int    `json:"media_count"`
	Website           string `json:"website,omitempty"`
	AccountType       string `json:"account_type,omitempty"`
}

type MediaItem struct {
	ID               string           `json:"id"`
	Caption          string           `json:"caption,omitempty"`
	MediaType        string           `json:"media_type"`
	MediaProductType string           `json:"media_product_type,omitempty"`
	MediaURL         string           `json:"media_url,omitempty"`
	ThumbnailURL     string           `json:"thumbnail_url,omitempty"`
	Permalink        string           `json:"permalink,omitempty"`
	Timestamp        time.Time        `json:"timestamp"`
	LikeCount        int              `json:"like_count"`
	CommentsCount    int              `json:"comments_count"`
	Insights         *MediaInsights   `json:"insights,omitempty"`
	Children         []MediaChildItem `json:"children,omitempty"`
}

type MediaChildItem struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url,omitempty"`
}

type MediaInsights struct {
	Impressions int `json:"impressions,omitempty"`
	Reach       int `json:"reach,omitempty"`
	Saved       int `json:"saved,omitempty"`
	Engagement  int `json:"engagement,omitempty"`
	VideoViews  int `json:"video_views,omitempty"`
	Shares      int `json:"shares,omitempty"`
}

type TokenResponse struct {
	AccessToken string          `json:"access_token"`
	TokenType   string          `json:"token_type,omitempty"`
	ExpiresIn   int64           `json:"expires_in,omitempty"`
	UserIDRaw   json.RawMessage `json:"user_id,omitempty"`
	Permissions []string        `json:"permissions,omitempty"`
}

func (t *TokenResponse) GetUserID() string {
	if t == nil || len(t.UserIDRaw) == 0 {
		return ""
	}
	s := strings.Trim(string(bytes.TrimSpace(t.UserIDRaw)), `"`)
	if s == "null" {
		return ""
	}
	return s
}

type GraphAPIError struct {
	Error struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int    `json:"code"`
		ErrorSubcode int    `json:"error_subcode"`
		FBTraceID    string `json:"fbtrace_id"`
	} `json:"error"`
}
