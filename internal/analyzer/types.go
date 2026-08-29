package analyzer

import (
	"reviewmysocialnetworks/internal/instagram"
	"time"
)

type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

type AccountReport struct {
	GeneratedAt       time.Time           `json:"generated_at"`
	Platform          string              `json:"platform"`
	Profile           Profile             `json:"profile"`
	OverallGrade      Grade               `json:"overall_grade"`
	OverallScore      int                 `json:"overall_score"`
	GradeTitle        string              `json:"grade_title"`
	GradeColor        string              `json:"grade_color"`
	ExecutiveSummary  string              `json:"executive_summary"`
	SubScores         SubScores           `json:"sub_scores"`
	EngagementMetrics EngagementMetrics   `json:"engagement_metrics"`
	CadenceMetrics    CadenceMetrics      `json:"cadence_metrics"`
	ContentMetrics    ContentMetrics      `json:"content_metrics"`
	GrowthMetrics     GrowthMetrics       `json:"growth_metrics"`
	Strengths         []string            `json:"strengths"`
	Weaknesses        []string            `json:"weaknesses"`
	Recommendations   []Recommendation    `json:"recommendations"`
	MediaAnalysis     []MediaAnalysisItem `json:"media_analysis"`
	TikTokMetrics     *TikTokMetrics      `json:"tiktok_metrics,omitempty"`
	DataCoverage      DataCoverage        `json:"data_coverage"`
}

type Profile struct {
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
	LikesCount        int    `json:"likes_count,omitempty"`
	IsVerified        bool   `json:"is_verified,omitempty"`
	ProfileURL        string `json:"profile_url,omitempty"`
}

type SubScores struct {
	EngagementScore     int `json:"engagement_score"`
	ConsistencyScore    int `json:"consistency_score"`
	ContentMixScore     int `json:"content_mix_score"`
	AudienceHealthScore int `json:"audience_health_score"`
}

type EngagementMetrics struct {
	AverageLikes          float64 `json:"average_likes"`
	AverageComments       float64 `json:"average_comments"`
	AverageSaves          float64 `json:"average_saves"`
	AverageEngagementRate float64 `json:"average_engagement_rate"`
	MedianEngagementRate  float64 `json:"median_engagement_rate"`
	TotalInteractions     int     `json:"total_interactions"`
	CommentToLikeRatio    float64 `json:"comment_to_like_ratio"`
	TopEngagingPostID     string  `json:"top_engaging_post_id"`
	TopEngagementRate     float64 `json:"top_engagement_rate"`
	BenchmarkComparison   string  `json:"benchmark_comparison"`
	AverageShares         float64 `json:"average_shares,omitempty"`
	AverageViews          float64 `json:"average_views,omitempty"`
	ViewEngagementRate    float64 `json:"view_engagement_rate,omitempty"`
}

type TikTokMetrics struct {
	TotalViews             int64   `json:"total_views"`
	AverageViews           float64 `json:"average_views"`
	MedianViews            float64 `json:"median_views"`
	MedianViewEngagement   float64 `json:"median_view_engagement"`
	TotalShares            int     `json:"total_shares"`
	AverageShares          float64 `json:"average_shares"`
	ShareRate              float64 `json:"share_rate"`
	ViewsPerFollower       float64 `json:"views_per_follower"`
	AverageDurationSeconds float64 `json:"average_duration_seconds"`
	ViralVideosCount       int     `json:"viral_videos_count"`
	TopVideoID             string  `json:"top_video_id"`
	TopViewCount           int64   `json:"top_view_count"`
}

type DataCoverage struct {
	AnalyzedPosts int      `json:"analyzed_posts"`
	MaxPosts      int      `json:"max_posts"`
	Available     []string `json:"available"`
	Unavailable   []string `json:"unavailable,omitempty"`
}

type CadenceMetrics struct {
	AverageDaysBetweenPosts float64        `json:"average_days_between_posts"`
	EstimatedPostsPerWeek   float64        `json:"estimated_posts_per_week"`
	EstimatedPostsPerMonth  float64        `json:"estimated_posts_per_month"`
	DaysSinceLastPost       int            `json:"days_since_last_post"`
	CadenceStatus           string         `json:"cadence_status"`
	BestPostingDay          string         `json:"best_posting_day"`
	BestPostingHour         int            `json:"best_posting_hour"`
	DayDistribution         map[string]int `json:"day_distribution"`
	HourDistribution        map[int]int    `json:"hour_distribution"`
}

type ContentMetrics struct {
	TotalAnalyzedPosts int                    `json:"total_analyzed_posts"`
	ImageCount         int                    `json:"image_count"`
	VideoCount         int                    `json:"video_count"`
	CarouselCount      int                    `json:"carousel_count"`
	ImagePercentage    float64                `json:"image_percentage"`
	VideoPercentage    float64                `json:"video_percentage"`
	CarouselPercentage float64                `json:"carousel_percentage"`
	BestPerformingType string                 `json:"best_performing_type"`
	AverageByFormat    map[string]FormatStats `json:"average_by_format"`
}

type FormatStats struct {
	Count                 int     `json:"count"`
	AverageLikes          float64 `json:"average_likes"`
	AverageComments       float64 `json:"average_comments"`
	AverageEngagementRate float64 `json:"average_engagement_rate"`
}

type GrowthMetrics struct {
	FollowerToFollowingRatio float64 `json:"follower_to_following_ratio"`
	AudienceHealthStatus     string  `json:"audience_health_status"`
	RecentTrendDirection     string  `json:"recent_trend_direction"`
	RecentTrendPercentage    float64 `json:"recent_trend_percentage"`
	EstimatedReachMultiplier float64 `json:"estimated_reach_multiplier"`
}

type Recommendation struct {
	Category string `json:"category"`
	Priority string `json:"priority"`
	Title    string `json:"title"`
	Action   string `json:"action"`
	Impact   string `json:"impact"`
}

type MediaAnalysisItem struct {
	ID                 string                   `json:"id"`
	Caption            string                   `json:"caption"`
	MediaType          string                   `json:"media_type"`
	MediaProductType   string                   `json:"media_product_type"`
	MediaURL           string                   `json:"media_url"`
	ThumbnailURL       string                   `json:"thumbnail_url"`
	Permalink          string                   `json:"permalink"`
	Timestamp          time.Time                `json:"timestamp"`
	LikeCount          int                      `json:"like_count"`
	CommentsCount      int                      `json:"comments_count"`
	EngagementRate     float64                  `json:"engagement_rate"`
	ViewEngagementRate float64                  `json:"view_engagement_rate,omitempty"`
	Insights           *instagram.MediaInsights `json:"insights,omitempty"`
	IsTopPerformer     bool                     `json:"is_top_performer"`
	ViewCount          int64                    `json:"view_count,omitempty"`
	ShareCount         int                      `json:"share_count,omitempty"`
	DurationSeconds    int                      `json:"duration_seconds,omitempty"`
	IsAIGC             bool                     `json:"is_aigc,omitempty"`
}
