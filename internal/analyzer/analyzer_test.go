package analyzer

import (
	"reviewmysocialnetworks/internal/instagram"
	"testing"
)

func TestAnalyzeAccount_Tiers(t *testing.T) {
	tests := []struct {
		tier         string
		expectedMin  int
		expectedMax  int
		expectedGrad Grade
	}{
		{tier: "A", expectedMin: 85, expectedMax: 100, expectedGrad: GradeA},
		{tier: "B", expectedMin: 70, expectedMax: 84, expectedGrad: GradeB},
		{tier: "D", expectedMin: 50, expectedMax: 69, expectedGrad: GradeD},
		{tier: "F", expectedMin: 0, expectedMax: 49, expectedGrad: GradeF},
	}

	for _, tt := range tests {
		t.Run("Tier_"+tt.tier, func(t *testing.T) {
			profile, media := instagram.GetMockAccount(tt.tier)
			report := AnalyzeAccount(profile, media)

			if report == nil {
				t.Fatalf("Expected report, got nil")
			}

			if report.OverallGrade != tt.expectedGrad {
				t.Errorf("Tier %s: expected grade %v, got %v (score %d)", tt.tier, tt.expectedGrad, report.OverallGrade, report.OverallScore)
			}

			if report.OverallScore < tt.expectedMin || report.OverallScore > tt.expectedMax {
				t.Errorf("Tier %s: expected score between %d and %d, got %d", tt.tier, tt.expectedMin, tt.expectedMax, report.OverallScore)
			}

			if len(report.Strengths) == 0 {
				t.Errorf("Tier %s: expected at least 1 strength", tt.tier)
			}

			if len(report.Recommendations) == 0 && tt.expectedGrad != GradeA {
				t.Errorf("Tier %s: expected recommendations", tt.tier)
			}

			if report.ExecutiveSummary == "" {
				t.Errorf("Tier %s: expected non-empty executive summary", tt.tier)
			}
		})
	}
}

func TestCadenceAndMediaMixCalculations(t *testing.T) {
	profile, media := instagram.GetMockAccount("A")
	report := AnalyzeAccount(profile, media)

	if report.ContentMetrics.TotalAnalyzedPosts != len(media) {
		t.Errorf("Expected total posts %d, got %d", len(media), report.ContentMetrics.TotalAnalyzedPosts)
	}

	if report.CadenceMetrics.EstimatedPostsPerWeek <= 0 {
		t.Errorf("Expected positive posts per week estimate, got %f", report.CadenceMetrics.EstimatedPostsPerWeek)
	}

	if report.EngagementMetrics.AverageEngagementRate <= 0 {
		t.Errorf("Expected positive engagement rate, got %f", report.EngagementMetrics.AverageEngagementRate)
	}
}
