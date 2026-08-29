package analyzer

import (
	"reviewmysocialnetworks/internal/tiktok"
	"strings"
	"testing"
)

func TestAnalyzeTikTokAccountTiers(t *testing.T) {
	tests := []struct {
		tier     string
		grade    Grade
		min, max int
	}{{"A", GradeA, 85, 100}, {"B", GradeB, 70, 84}, {"D", GradeD, 50, 69}, {"F", GradeF, 0, 49}}
	for _, test := range tests {
		t.Run(test.tier, func(t *testing.T) {
			profile, videos := tiktok.GetMockAccount(test.tier)
			report := AnalyzeTikTokAccount(profile, videos)
			if report.Platform != "tiktok" || report.OverallGrade != test.grade || report.OverallScore < test.min || report.OverallScore > test.max {
				t.Fatalf("tier %s: grade=%s score=%d", test.tier, report.OverallGrade, report.OverallScore)
			}
			if report.TikTokMetrics == nil || report.TikTokMetrics.AverageViews == 0 || len(report.Recommendations) < 3 {
				t.Fatalf("incomplete TikTok report: %+v", report)
			}
			if len(report.DataCoverage.Unavailable) == 0 {
				t.Fatal("expected API coverage limitations")
			}
		})
	}
}

func TestAnalyzeTikTokDoesNotMutateInput(t *testing.T) {
	profile, videos := tiktok.GetMockAccount("A")
	first := videos[0].ID
	videos[0], videos[1] = videos[1], videos[0]
	AnalyzeTikTokAccount(profile, videos)
	if videos[1].ID != first {
		t.Fatal("analyzer mutated caller video order")
	}
}

func TestAnalyzeTikTokSeparatesFollowerAndViewEngagement(t *testing.T) {
	profile := &tiktok.UserProfile{OpenID: "id", Username: "creator", FollowerCount: 1000, FollowingCount: 100}
	videos := []tiktok.Video{{ID: "video", CreateTime: 1_700_000_000, ViewCount: 500, LikeCount: 40, CommentCount: 5, ShareCount: 5}}
	report := AnalyzeTikTokAccount(profile, videos)
	if report.EngagementMetrics.AverageEngagementRate != 5 {
		t.Fatalf("follower engagement = %v, want 5", report.EngagementMetrics.AverageEngagementRate)
	}
	if report.EngagementMetrics.ViewEngagementRate != 10 || report.MediaAnalysis[0].ViewEngagementRate != 10 {
		t.Fatalf("view engagement not preserved: report=%v item=%v", report.EngagementMetrics.ViewEngagementRate, report.MediaAnalysis[0].ViewEngagementRate)
	}
}

func TestAnalyzeTikTokWithoutVideosDoesNotInventPostingTime(t *testing.T) {
	profile := &tiktok.UserProfile{OpenID: "id", Username: "new_creator", FollowerCount: 10}
	report := AnalyzeTikTokAccount(profile, nil)
	if report.OverallGrade != GradeF || len(report.Recommendations) == 0 {
		t.Fatalf("unexpected empty-account report: grade=%s recommendations=%d", report.OverallGrade, len(report.Recommendations))
	}
	for _, recommendation := range report.Recommendations {
		if recommendation.Category == "Frecuencia" && recommendation.Action == "" {
			t.Fatal("empty cadence recommendation")
		}
		if recommendation.Category == "Frecuencia" && (strings.Contains(recommendation.Action, "00:00") || strings.Contains(recommendation.Action, "días recomendados")) {
			t.Fatalf("invented posting time: %s", recommendation.Action)
		}
	}
}
