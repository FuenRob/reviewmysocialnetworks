package analyzer

import (
	"math"
	"reviewmysocialnetworks/internal/instagram"
	"sort"
	"time"
)

var spanishDays = [...]string{
	"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado",
}

func AnalyzeAccount(profile *instagram.UserProfile, media []instagram.MediaItem) *AccountReport {
	report := &AccountReport{
		GeneratedAt: time.Now(),
		Profile:     *profile,
	}

	orderedMedia := append([]instagram.MediaItem(nil), media...)
	sort.Slice(orderedMedia, func(i, j int) bool {
		return orderedMedia[i].Timestamp.After(orderedMedia[j].Timestamp)
	})

	engagementMetrics, mediaAnalysis := calculateEngagement(profile, orderedMedia)
	report.EngagementMetrics = engagementMetrics
	report.MediaAnalysis = mediaAnalysis

	cadenceMetrics := calculateCadence(orderedMedia)
	report.CadenceMetrics = cadenceMetrics

	contentMetrics := calculateContent(orderedMedia, profile.FollowersCount)
	report.ContentMetrics = contentMetrics

	growthMetrics := calculateGrowth(profile, mediaAnalysis)
	report.GrowthMetrics = growthMetrics

	subScores, overallScore, grade, title, color := calculateScores(engagementMetrics, cadenceMetrics, contentMetrics, growthMetrics, profile)
	report.SubScores = subScores
	report.OverallScore = overallScore
	report.OverallGrade = grade
	report.GradeTitle = title
	report.GradeColor = color

	generateQualitativeInsights(report)

	return report
}

func calculateEngagement(profile *instagram.UserProfile, media []instagram.MediaItem) (EngagementMetrics, []MediaAnalysisItem) {
	metrics := EngagementMetrics{
		BenchmarkComparison: "Promedio estándar del sector: 1.5% - 2.5%",
	}

	if len(media) == 0 {
		return metrics, nil
	}

	var totalLikes, totalComments, totalSaves int
	var totalRate float64
	rates := make([]float64, 0, len(media))
	analysisItems := make([]MediaAnalysisItem, 0, len(media))

	topRate := -1.0
	var topPostID string
	topIndex := -1

	for _, item := range media {
		interactions := item.LikeCount + item.CommentsCount
		saves := 0
		if item.Insights != nil {
			saves = item.Insights.Saved
			interactions += saves
		}

		rate := 0.0
		if profile.FollowersCount > 0 {
			rate = (float64(interactions) / float64(profile.FollowersCount)) * 100.0
		}

		if rate > topRate {
			topRate = rate
			topPostID = item.ID
			topIndex = len(analysisItems)
		}

		totalLikes += item.LikeCount
		totalComments += item.CommentsCount
		totalSaves += saves
		totalRate += rate
		rates = append(rates, rate)

		analysisItems = append(analysisItems, MediaAnalysisItem{
			ID:               item.ID,
			Caption:          item.Caption,
			MediaType:        item.MediaType,
			MediaProductType: item.MediaProductType,
			MediaURL:         item.MediaURL,
			ThumbnailURL:     item.ThumbnailURL,
			Permalink:        item.Permalink,
			Timestamp:        item.Timestamp,
			LikeCount:        item.LikeCount,
			CommentsCount:    item.CommentsCount,
			EngagementRate:   math.Round(rate*100) / 100,
			Insights:         item.Insights,
		})
	}

	if topIndex >= 0 {
		analysisItems[topIndex].IsTopPerformer = true
	}

	n := float64(len(media))
	metrics.AverageLikes = math.Round((float64(totalLikes)/n)*10) / 10
	metrics.AverageComments = math.Round((float64(totalComments)/n)*10) / 10
	metrics.AverageSaves = math.Round((float64(totalSaves)/n)*10) / 10
	metrics.AverageEngagementRate = math.Round((totalRate/n)*100) / 100
	metrics.TotalInteractions = totalLikes + totalComments + totalSaves

	if totalLikes > 0 {
		metrics.CommentToLikeRatio = math.Round((float64(totalComments)/float64(totalLikes))*10000) / 100
	}

	metrics.TopEngagingPostID = topPostID
	metrics.TopEngagementRate = math.Round(topRate*100) / 100

	sort.Float64s(rates)
	if len(rates)%2 == 0 {
		metrics.MedianEngagementRate = math.Round(((rates[len(rates)/2-1]+rates[len(rates)/2])/2.0)*100) / 100
	} else {
		metrics.MedianEngagementRate = math.Round(rates[len(rates)/2]*100) / 100
	}

	if metrics.AverageEngagementRate >= 4.0 {
		metrics.BenchmarkComparison = "🌟 Sobresaliente (Muy superior al 2.0% promedio)"
	} else if metrics.AverageEngagementRate >= 2.0 {
		metrics.BenchmarkComparison = "✅ Saludable (Dentro del rango óptimo de 2.0% - 3.5%)"
	} else if metrics.AverageEngagementRate >= 0.8 {
		metrics.BenchmarkComparison = "⚠️ Moderado (Por debajo del promedio óptimo)"
	} else {
		metrics.BenchmarkComparison = "🚨 Crítico (Muy por debajo del estándar del 1.5%)"
	}

	return metrics, analysisItems
}

func calculateCadence(media []instagram.MediaItem) CadenceMetrics {
	metrics := CadenceMetrics{
		DayDistribution:  make(map[string]int),
		HourDistribution: make(map[int]int),
		CadenceStatus:    "Sin publicaciones suficientes",
	}

	if len(media) == 0 {
		return metrics
	}

	for _, day := range spanishDays {
		metrics.DayDistribution[day] = 0
	}

	dayEngagements := make(map[string][]float64)
	hourEngagements := make(map[int][]float64)

	for _, item := range media {
		dayName := spanishDays[item.Timestamp.Weekday()]
		hour := item.Timestamp.Hour()

		metrics.DayDistribution[dayName]++
		metrics.HourDistribution[hour]++

		interactions := float64(item.LikeCount + item.CommentsCount)
		dayEngagements[dayName] = append(dayEngagements[dayName], interactions)
		hourEngagements[hour] = append(hourEngagements[hour], interactions)
	}

	bestDay := "Miércoles"
	bestDayAvg := -1.0
	for day, scores := range dayEngagements {
		if len(scores) > 0 {
			var sum float64
			for _, s := range scores {
				sum += s
			}
			avg := sum / float64(len(scores))
			if avg > bestDayAvg {
				bestDayAvg = avg
				bestDay = day
			}
		}
	}
	metrics.BestPostingDay = bestDay

	bestHour := 18
	bestHourAvg := -1.0
	for hour, scores := range hourEngagements {
		if len(scores) > 0 {
			var sum float64
			for _, s := range scores {
				sum += s
			}
			avg := sum / float64(len(scores))
			if avg > bestHourAvg {
				bestHourAvg = avg
				bestHour = hour
			}
		}
	}
	metrics.BestPostingHour = bestHour

	metrics.DaysSinceLastPost = int(time.Since(media[0].Timestamp).Hours() / 24)
	if metrics.DaysSinceLastPost < 0 {
		metrics.DaysSinceLastPost = 0
	}

	if len(media) >= 2 {
		var totalGapDays float64
		gapsCount := 0
		for i := 0; i < len(media)-1; i++ {
			gap := media[i].Timestamp.Sub(media[i+1].Timestamp).Hours() / 24
			if gap > 0 {
				totalGapDays += gap
				gapsCount++
			}
		}

		if gapsCount > 0 {
			avgDays := totalGapDays / float64(gapsCount)
			metrics.AverageDaysBetweenPosts = math.Round(avgDays*10) / 10
			if avgDays > 0 {
				metrics.EstimatedPostsPerWeek = math.Round((7.0/avgDays)*10) / 10
				metrics.EstimatedPostsPerMonth = math.Round((30.4/avgDays)*10) / 10
			}
		}
	} else {
		metrics.AverageDaysBetweenPosts = float64(metrics.DaysSinceLastPost)
	}

	if metrics.DaysSinceLastPost > 60 {
		metrics.CadenceStatus = "Inactiva (Más de 2 meses sin publicar)"
	} else if metrics.DaysSinceLastPost > 20 {
		metrics.CadenceStatus = "En pausa prolongada"
	} else if metrics.EstimatedPostsPerWeek >= 3.0 {
		metrics.CadenceStatus = "Óptima y muy activa (3+ posts/sem)"
	} else if metrics.EstimatedPostsPerWeek >= 1.0 {
		metrics.CadenceStatus = "Moderada y constante (1-2 posts/sem)"
	} else {
		metrics.CadenceStatus = "Irregular / Baja frecuencia"
	}

	return metrics
}

func calculateContent(media []instagram.MediaItem, followers int) ContentMetrics {
	metrics := ContentMetrics{
		TotalAnalyzedPosts: len(media),
		AverageByFormat:    make(map[string]FormatStats),
	}

	if len(media) == 0 {
		return metrics
	}

	var imgCount, vidCount, carCount int
	var imgLikes, vidLikes, carLikes int
	var imgComms, vidComms, carComms int
	var imgRates, vidRates, carRates float64

	for _, item := range media {
		interactions := item.LikeCount + item.CommentsCount
		rate := 0.0
		if followers > 0 {
			rate = (float64(interactions) / float64(followers)) * 100.0
		}

		switch item.MediaType {
		case "VIDEO":
			vidCount++
			vidLikes += item.LikeCount
			vidComms += item.CommentsCount
			vidRates += rate
		case "CAROUSEL_ALBUM":
			carCount++
			carLikes += item.LikeCount
			carComms += item.CommentsCount
			carRates += rate
		case "IMAGE":
			fallthrough
		default:
			imgCount++
			imgLikes += item.LikeCount
			imgComms += item.CommentsCount
			imgRates += rate
		}
	}

	total := float64(len(media))
	metrics.ImageCount = imgCount
	metrics.VideoCount = vidCount
	metrics.CarouselCount = carCount
	metrics.ImagePercentage = math.Round((float64(imgCount)/total)*1000) / 10
	metrics.VideoPercentage = math.Round((float64(vidCount)/total)*1000) / 10
	metrics.CarouselPercentage = math.Round((float64(carCount)/total)*1000) / 10

	bestType := "IMAGE"
	bestAvgRate := -1.0

	if imgCount > 0 {
		avgRate := imgRates / float64(imgCount)
		metrics.AverageByFormat["IMAGE"] = FormatStats{
			Count:                 imgCount,
			AverageLikes:          math.Round((float64(imgLikes)/float64(imgCount))*10) / 10,
			AverageComments:       math.Round((float64(imgComms)/float64(imgCount))*10) / 10,
			AverageEngagementRate: math.Round(avgRate*100) / 100,
		}
		if avgRate > bestAvgRate {
			bestAvgRate = avgRate
			bestType = "IMAGE"
		}
	}

	if vidCount > 0 {
		avgRate := vidRates / float64(vidCount)
		metrics.AverageByFormat["VIDEO"] = FormatStats{
			Count:                 vidCount,
			AverageLikes:          math.Round((float64(vidLikes)/float64(vidCount))*10) / 10,
			AverageComments:       math.Round((float64(vidComms)/float64(vidCount))*10) / 10,
			AverageEngagementRate: math.Round(avgRate*100) / 100,
		}
		if avgRate > bestAvgRate {
			bestAvgRate = avgRate
			bestType = "VIDEO"
		}
	}

	if carCount > 0 {
		avgRate := carRates / float64(carCount)
		metrics.AverageByFormat["CAROUSEL_ALBUM"] = FormatStats{
			Count:                 carCount,
			AverageLikes:          math.Round((float64(carLikes)/float64(carCount))*10) / 10,
			AverageComments:       math.Round((float64(carComms)/float64(carCount))*10) / 10,
			AverageEngagementRate: math.Round(avgRate*100) / 100,
		}
		if avgRate > bestAvgRate {
			bestAvgRate = avgRate
			bestType = "CAROUSEL_ALBUM"
		}
	}

	metrics.BestPerformingType = bestType
	return metrics
}

func calculateGrowth(profile *instagram.UserProfile, items []MediaAnalysisItem) GrowthMetrics {
	metrics := GrowthMetrics{
		RecentTrendDirection: "Estable ➜",
	}

	if profile.FollowsCount > 0 {
		metrics.FollowerToFollowingRatio = math.Round((float64(profile.FollowersCount)/float64(profile.FollowsCount))*100) / 100
	} else {
		metrics.FollowerToFollowingRatio = float64(profile.FollowersCount)
	}

	if metrics.FollowerToFollowingRatio >= 5.0 {
		metrics.AudienceHealthStatus = "Excelente (Autoridad consolidada)"
	} else if metrics.FollowerToFollowingRatio >= 1.5 {
		metrics.AudienceHealthStatus = "Saludable y natural"
	} else if metrics.FollowerToFollowingRatio >= 0.8 {
		metrics.AudienceHealthStatus = "Equilibrada pero con potencial de optimización"
	} else {
		metrics.AudienceHealthStatus = "Desbalanceada (Sigue a demasiadas cuentas)"
	}

	if len(items) >= 4 {
		half := len(items) / 2
		var recentSum, olderSum float64
		for i := 0; i < half; i++ {
			recentSum += items[i].EngagementRate
		}
		for i := half; i < len(items); i++ {
			olderSum += items[i].EngagementRate
		}
		recentAvg := recentSum / float64(half)
		olderAvg := olderSum / float64(len(items)-half)

		if olderAvg > 0 {
			diffPercent := ((recentAvg - olderAvg) / olderAvg) * 100.0
			metrics.RecentTrendPercentage = math.Round(diffPercent*10) / 10
			if diffPercent > 10.0 {
				metrics.RecentTrendDirection = "Creciente ↗"
			} else if diffPercent < -10.0 {
				metrics.RecentTrendDirection = "Decreciente ↘"
			} else {
				metrics.RecentTrendDirection = "Estable ➜"
			}
		}
	}

	if profile.FollowersCount > 0 {
		metrics.EstimatedReachMultiplier = 1.25
	}

	return metrics
}
