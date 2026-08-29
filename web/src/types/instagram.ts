export type Grade = 'A' | 'B' | 'D' | 'F';

export interface UserProfile {
  id: string;
  username: string;
  name?: string;
  biography?: string;
  profile_picture_url?: string;
  followers_count: number;
  follows_count: number;
  media_count: number;
  website?: string;
  account_type?: string;
  likes_count?: number;
  is_verified?: boolean;
  profile_url?: string;
}

export interface MediaInsights {
  impressions?: number;
  reach?: number;
  saved?: number;
  engagement?: number;
  video_views?: number;
  shares?: number;
}

export interface MediaAnalysisItem {
  id: string;
  caption: string;
  media_type: string;
  media_product_type?: string;
  media_url: string;
  thumbnail_url?: string;
  permalink: string;
  timestamp: string;
  like_count: number;
  comments_count: number;
  engagement_rate: number;
  view_engagement_rate?: number;
  insights?: MediaInsights;
  is_top_performer?: boolean;
  view_count?: number;
  share_count?: number;
  duration_seconds?: number;
  is_aigc?: boolean;
}

export interface SubScores {
  engagement_score: number;
  consistency_score: number;
  content_mix_score: number;
  audience_health_score: number;
}

export interface EngagementMetrics {
  average_likes: number;
  average_comments: number;
  average_saves: number;
  average_engagement_rate: number;
  median_engagement_rate: number;
  total_interactions: number;
  comment_to_like_ratio: number;
  top_engaging_post_id: string;
  top_engagement_rate: number;
  benchmark_comparison: string;
  average_shares?: number;
  average_views?: number;
  view_engagement_rate?: number;
}

export interface TikTokMetrics {
  total_views: number;
  average_views: number;
  median_views: number;
  median_view_engagement: number;
  total_shares: number;
  average_shares: number;
  share_rate: number;
  views_per_follower: number;
  average_duration_seconds: number;
  viral_videos_count: number;
  top_video_id: string;
  top_view_count: number;
}

export interface DataCoverage {
  analyzed_posts: number;
  max_posts: number;
  available: string[];
  unavailable?: string[];
}

export interface CadenceMetrics {
  average_days_between_posts: number;
  estimated_posts_per_week: number;
  estimated_posts_per_month: number;
  days_since_last_post: number;
  cadence_status: string;
  best_posting_day: string;
  best_posting_hour: number;
  day_distribution: Record<string, number>;
  hour_distribution: Record<number, number>;
}

export interface FormatStats {
  count: number;
  average_likes: number;
  average_comments: number;
  average_engagement_rate: number;
}

export interface ContentMetrics {
  total_analyzed_posts: number;
  image_count: number;
  video_count: number;
  carousel_count: number;
  image_percentage: number;
  video_percentage: number;
  carousel_percentage: number;
  best_performing_type: string;
  average_by_format: Record<string, FormatStats>;
}

export interface GrowthMetrics {
  follower_to_following_ratio: number;
  audience_health_status: string;
  recent_trend_direction: string;
  recent_trend_percentage: number;
  estimated_reach_multiplier: number;
}

export interface Recommendation {
  category: string;
  priority: 'Alta' | 'Media' | 'Baja';
  title: string;
  action: string;
  impact: string;
}

export interface AccountReport {
  generated_at: string;
  platform: 'instagram' | 'tiktok';
  profile: UserProfile;
  overall_grade: Grade;
  overall_score: number;
  grade_title: string;
  grade_color: string;
  executive_summary: string;
  sub_scores: SubScores;
  engagement_metrics: EngagementMetrics;
  cadence_metrics: CadenceMetrics;
  content_metrics: ContentMetrics;
  growth_metrics: GrowthMetrics;
  strengths: string[];
  weaknesses: string[];
  recommendations: Recommendation[];
  media_analysis: MediaAnalysisItem[];
  tiktok_metrics?: TikTokMetrics;
  data_coverage: DataCoverage;
}
