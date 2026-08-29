package tiktok

type UserProfile struct {
	OpenID          string `json:"open_id"`
	UnionID         string `json:"union_id,omitempty"`
	Username        string `json:"username"`
	DisplayName     string `json:"display_name"`
	BioDescription  string `json:"bio_description,omitempty"`
	AvatarURL       string `json:"avatar_url,omitempty"`
	ProfileDeepLink string `json:"profile_deep_link,omitempty"`
	IsVerified      bool   `json:"is_verified"`
	FollowerCount   int    `json:"follower_count"`
	FollowingCount  int    `json:"following_count"`
	LikesCount      int    `json:"likes_count"`
	VideoCount      int    `json:"video_count"`
}

type Video struct {
	ID               string `json:"id"`
	CreateTime       int64  `json:"create_time"`
	CoverImageURL    string `json:"cover_image_url,omitempty"`
	ShareURL         string `json:"share_url,omitempty"`
	VideoDescription string `json:"video_description,omitempty"`
	Title            string `json:"title,omitempty"`
	Duration         int    `json:"duration,omitempty"`
	LikeCount        int    `json:"like_count"`
	CommentCount     int    `json:"comment_count"`
	ShareCount       int    `json:"share_count"`
	ViewCount        int64  `json:"view_count"`
	IsAIGC           bool   `json:"is_aigc,omitempty"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	OpenID           string `json:"open_id"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LogID   string `json:"log_id,omitempty"`
}
