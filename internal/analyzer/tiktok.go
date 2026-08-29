package analyzer

import (
	"fmt"
	"math"
	"reviewmysocialnetworks/internal/instagram"
	"reviewmysocialnetworks/internal/tiktok"
	"slices"
	"sort"
	"time"
)

func AnalyzeTikTokAccount(profile *tiktok.UserProfile, videos []tiktok.Video) *AccountReport {
	ordered := append([]tiktok.Video(nil), videos...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].CreateTime > ordered[j].CreateTime })

	report := &AccountReport{
		GeneratedAt: time.Now(), Platform: "tiktok",
		Profile: Profile{
			ID: profile.OpenID, Username: profile.Username, Name: profile.DisplayName,
			Biography: profile.BioDescription, ProfilePictureURL: profile.AvatarURL,
			FollowersCount: profile.FollowerCount, FollowsCount: profile.FollowingCount,
			MediaCount: profile.VideoCount, LikesCount: profile.LikesCount,
			IsVerified: profile.IsVerified, ProfileURL: profile.ProfileDeepLink,
			AccountType: "TIKTOK",
		},
		DataCoverage: DataCoverage{
			AnalyzedPosts: len(ordered), MaxPosts: 50,
			Available:   []string{"perfil", "seguidores", "vídeos públicos", "visualizaciones", "likes", "comentarios", "compartidos", "duración"},
			Unavailable: []string{"retención", "tiempo medio de reproducción", "favoritos", "fuentes de tráfico", "crecimiento histórico de seguidores"},
		},
	}

	report.EngagementMetrics, report.MediaAnalysis, report.TikTokMetrics = calculateTikTokEngagement(&report.Profile, ordered)
	report.CadenceMetrics = calculateTikTokCadence(ordered)
	report.ContentMetrics = calculateTikTokContent(report.MediaAnalysis)
	trendItems := append([]MediaAnalysisItem(nil), report.MediaAnalysis...)
	for i := range trendItems {
		trendItems[i].EngagementRate = trendItems[i].ViewEngagementRate
	}
	report.GrowthMetrics = calculateGrowth(&report.Profile, trendItems)
	report.SubScores, report.OverallScore, report.OverallGrade, report.GradeTitle, report.GradeColor = calculateTikTokScores(report)
	generateTikTokInsights(report)
	return report
}

func calculateTikTokEngagement(profile *Profile, videos []tiktok.Video) (EngagementMetrics, []MediaAnalysisItem, *TikTokMetrics) {
	eng := EngagementMetrics{BenchmarkComparison: "Sin vídeos suficientes para comparar"}
	tm := &TikTokMetrics{}
	if len(videos) == 0 {
		return eng, nil, tm
	}

	items := make([]MediaAnalysisItem, 0, len(videos))
	followerRates := make([]float64, 0, len(videos))
	viewRates := make([]float64, 0, len(videos))
	views := make([]int64, 0, len(videos))
	var totalLikes, totalComments, totalShares int
	var totalViews int64
	var totalFollowerRate, totalViewRate float64
	var totalDuration int
	topViewIndex := -1
	topEngagementRate := -1.0

	for _, video := range videos {
		interactions := video.LikeCount + video.CommentCount + video.ShareCount
		followerRate := 0.0
		if profile.FollowersCount > 0 {
			followerRate = float64(interactions) / float64(profile.FollowersCount) * 100
		}
		viewRate := 0.0
		if video.ViewCount > 0 {
			viewRate = float64(interactions) / float64(video.ViewCount) * 100
		}
		caption := video.VideoDescription
		if caption == "" {
			caption = video.Title
		}
		items = append(items, MediaAnalysisItem{
			ID: video.ID, Caption: caption, MediaType: "VIDEO", MediaProductType: "TIKTOK",
			ThumbnailURL: video.CoverImageURL, Permalink: video.ShareURL,
			Timestamp: time.Unix(video.CreateTime, 0), LikeCount: video.LikeCount,
			CommentsCount: video.CommentCount, ShareCount: video.ShareCount,
			ViewCount: video.ViewCount, DurationSeconds: video.Duration,
			EngagementRate: round2(followerRate), ViewEngagementRate: round2(viewRate), IsAIGC: video.IsAIGC,
		})
		if video.ViewCount > tm.TopViewCount {
			tm.TopViewCount, tm.TopVideoID, topViewIndex = video.ViewCount, video.ID, len(items)-1
		}
		if viewRate > topEngagementRate {
			topEngagementRate = viewRate
			eng.TopEngagingPostID = video.ID
		}
		if profile.FollowersCount > 0 && float64(video.ViewCount) >= float64(profile.FollowersCount)*1.5 {
			tm.ViralVideosCount++
		}
		totalLikes += video.LikeCount
		totalComments += video.CommentCount
		totalShares += video.ShareCount
		totalViews += video.ViewCount
		totalDuration += video.Duration
		totalFollowerRate += followerRate
		totalViewRate += viewRate
		followerRates = append(followerRates, followerRate)
		viewRates = append(viewRates, viewRate)
		views = append(views, video.ViewCount)
	}
	if topViewIndex >= 0 {
		items[topViewIndex].IsTopPerformer = true
	}

	n := float64(len(videos))
	eng.AverageLikes = round1(float64(totalLikes) / n)
	eng.AverageComments = round1(float64(totalComments) / n)
	eng.AverageShares = round1(float64(totalShares) / n)
	eng.AverageViews = round1(float64(totalViews) / n)
	eng.AverageEngagementRate = round2(totalFollowerRate / n)
	eng.ViewEngagementRate = round2(totalViewRate / n)
	eng.TotalInteractions = totalLikes + totalComments + totalShares
	if totalLikes > 0 {
		eng.CommentToLikeRatio = round2(float64(totalComments) / float64(totalLikes) * 100)
	}
	eng.MedianEngagementRate = medianFloat(followerRates)
	eng.TopEngagementRate = round2(topEngagementRate)
	eng.BenchmarkComparison = tikTokBenchmark(eng.ViewEngagementRate)

	tm.TotalViews = totalViews
	tm.AverageViews = eng.AverageViews
	tm.MedianViews = medianInt64(views)
	tm.MedianViewEngagement = medianFloat(viewRates)
	tm.TotalShares = totalShares
	tm.AverageShares = eng.AverageShares
	if totalViews > 0 {
		tm.ShareRate = round2(float64(totalShares) / float64(totalViews) * 100)
	}
	if profile.FollowersCount > 0 {
		tm.ViewsPerFollower = round2(tm.AverageViews / float64(profile.FollowersCount))
	}
	tm.AverageDurationSeconds = round1(float64(totalDuration) / n)
	return eng, items, tm
}

func calculateTikTokCadence(videos []tiktok.Video) CadenceMetrics {
	converted := make([]instagram.MediaItem, 0, len(videos))
	for _, video := range videos {
		converted = append(converted, instagram.MediaItem{Timestamp: time.Unix(video.CreateTime, 0), LikeCount: video.LikeCount, CommentsCount: video.CommentCount + video.ShareCount})
	}
	return calculateCadence(converted)
}

func calculateTikTokContent(items []MediaAnalysisItem) ContentMetrics {
	metrics := ContentMetrics{TotalAnalyzedPosts: len(items), VideoCount: len(items), AverageByFormat: map[string]FormatStats{}, BestPerformingType: "VIDEO"}
	if len(items) == 0 {
		return metrics
	}
	metrics.VideoPercentage = 100
	var likes, comments int
	var engagement float64
	for _, item := range items {
		likes += item.LikeCount
		comments += item.CommentsCount
		engagement += item.ViewEngagementRate
	}
	n := float64(len(items))
	metrics.AverageByFormat["VIDEO"] = FormatStats{Count: len(items), AverageLikes: round1(float64(likes) / n), AverageComments: round1(float64(comments) / n), AverageEngagementRate: round2(engagement / n)}
	return metrics
}

func calculateTikTokScores(report *AccountReport) (SubScores, int, Grade, string, string) {
	engScore := scaleScore(report.EngagementMetrics.ViewEngagementRate, []scorePoint{{0, 10}, {1.5, 50}, {3, 70}, {5, 85}, {8, 95}, {12, 100}})
	posts := report.CadenceMetrics.EstimatedPostsPerWeek
	days := report.CadenceMetrics.DaysSinceLastPost
	cadScore := 35.0
	switch {
	case posts >= 5 && days <= 3:
		cadScore = 95
	case posts >= 3 && days <= 5:
		cadScore = 85
	case posts >= 2 && days <= 7:
		cadScore = 72
	case posts >= 1 && days <= 14:
		cadScore = 55
	case days > 60:
		cadScore = 15
	}
	reachScore := scaleScore(report.TikTokMetrics.ViewsPerFollower, []scorePoint{{0, 10}, {.2, 45}, {.5, 70}, {1, 85}, {2, 95}, {4, 100}})
	if report.TikTokMetrics.ShareRate >= 1 {
		reachScore = math.Min(100, reachScore+5)
	}
	ratio := report.GrowthMetrics.FollowerToFollowingRatio
	audScore := scaleScore(ratio, []scorePoint{{0, 20}, {.7, 55}, {1.2, 75}, {3, 85}, {10, 95}, {20, 100}})
	sub := SubScores{EngagementScore: int(math.Round(engScore)), ConsistencyScore: int(math.Round(cadScore)), ContentMixScore: int(math.Round(reachScore)), AudienceHealthScore: int(math.Round(audScore))}
	overall := clamp(int(math.Round(engScore*.35+cadScore*.25+reachScore*.20+audScore*.20)), 0, 100)
	grade, title, color := gradeForScore(overall)
	return sub, overall, grade, title, color
}

func generateTikTokInsights(report *AccountReport) {
	tm := report.TikTokMetrics
	eng := report.EngagementMetrics
	cad := report.CadenceMetrics
	var strengths, weaknesses []string
	var recs []Recommendation

	if eng.ViewEngagementRate >= 5 {
		strengths = append(strengths, fmt.Sprintf("La audiencia responde con fuerza: %.2f%% de interacciones por visualización.", eng.ViewEngagementRate))
	} else if eng.ViewEngagementRate < 2 {
		weaknesses = append(weaknesses, fmt.Sprintf("La conversión de visualizaciones en interacción es baja (%.2f%%).", eng.ViewEngagementRate))
		recs = append(recs, Recommendation{"Engagement", "Alta", "Reforzar el gancho y el CTA", "Abre con el resultado o conflicto en el primer segundo y termina con una pregunta concreta que invite a comentar o compartir.", "Más retención inicial e interacción cualificada."})
	}

	if tm.ViewsPerFollower >= 1 {
		strengths = append(strengths, fmt.Sprintf("El alcance medio equivale a %.2fx tu base de seguidores: el contenido ya rompe la burbuja de la audiencia propia.", tm.ViewsPerFollower))
	} else if tm.ViewsPerFollower < .35 {
		weaknesses = append(weaknesses, fmt.Sprintf("Cada vídeo alcanza de media solo %.2fx tu base de seguidores.", tm.ViewsPerFollower))
		recs = append(recs, Recommendation{"Alcance", "Alta", "Crear series repetibles", "Convierte el tema con mejor rendimiento en una serie de 5 vídeos con promesa, formato y portada reconocibles.", "Aumenta las señales de retorno y el alcance en Para ti."})
	}

	if tm.ShareRate >= .7 {
		strengths = append(strengths, fmt.Sprintf("Buen potencial de distribución: %.2f%% de las visualizaciones terminan en compartido.", tm.ShareRate))
	} else {
		weaknesses = append(weaknesses, fmt.Sprintf("La tasa de compartidos (%.2f%%) tiene margen de mejora.", tm.ShareRate))
		recs = append(recs, Recommendation{"Viralidad", "Media", "Diseñar contenido compartible", "Prioriza listas útiles, errores frecuentes, opiniones contrastadas y plantillas que la audiencia quiera enviar a otra persona.", "Más distribución orgánica fuera de tus seguidores."})
	}

	if cad.EstimatedPostsPerWeek >= 3 && cad.DaysSinceLastPost <= 5 {
		strengths = append(strengths, fmt.Sprintf("Ritmo competitivo de %.1f vídeos por semana y actividad reciente.", cad.EstimatedPostsPerWeek))
	} else {
		weaknesses = append(weaknesses, fmt.Sprintf("La cadencia actual (%.1f vídeos/semana) limita el aprendizaje del formato.", cad.EstimatedPostsPerWeek))
		action := "Publica 3 vídeos iniciales en días alternos para crear una base de comparación y descubrir qué temas generan respuesta."
		if len(report.MediaAnalysis) > 0 {
			action = fmt.Sprintf("Planifica 3 a 5 vídeos por semana; empieza los %s alrededor de las %02d:00, el mejor patrón observado.", pluralizeDay(cad.BestPostingDay), cad.BestPostingHour)
		}
		recs = append(recs, Recommendation{"Frecuencia", "Alta", "Publicar con una cadencia sostenible", action, "Más oportunidades de validación y crecimiento estable."})
	}

	if tm.AverageDurationSeconds > 45 {
		recs = append(recs, Recommendation{"Retención", "Media", "Comparar versiones cortas y largas", fmt.Sprintf("Tu duración media es de %.0f s. Prueba versiones de 15–30 s del mismo concepto y compara interacción por visualización.", tm.AverageDurationSeconds), "Identifica la duración que mejor convierte sin asumir datos de retención no disponibles."})
	} else {
		recs = append(recs, Recommendation{"Retención", "Media", "Optimizar los primeros tres segundos", "Elimina saludos e introducciones; muestra primero el resultado, añade texto en pantalla y cambia el estímulo visual cada 2–4 segundos.", "Mejora la probabilidad de completar y volver a ver el vídeo."})
	}

	if report.Profile.Biography == "" || report.Profile.ProfileURL == "" {
		weaknesses = append(weaknesses, "El perfil no comunica con suficiente claridad la propuesta o el siguiente paso.")
		recs = append(recs, Recommendation{"Perfil", "Media", "Convertir visitas de perfil", "Define en una frase para quién creas contenido, qué resultado ofreces y una llamada a la acción única en la bio.", "Mejor conversión de visitas en seguidores."})
	}
	if len(report.MediaAnalysis) >= 4 {
		recs = append(recs, Recommendation{"Experimentación", "Baja", "Duplicar patrones ganadores", fmt.Sprintf("Desglosa el vídeo con más visualizaciones (%s) en gancho, tema, duración y CTA; reutiliza dos de esos elementos en tres nuevas piezas.", tm.TopVideoID), "Aprendizaje basado en el historial real de la cuenta."})
	}
	recs = append(recs, Recommendation{"Descubrimiento", "Baja", "Reforzar el contexto semántico", "Nombra el tema principal en voz, texto en pantalla y descripción; utiliza de 3 a 5 hashtags específicos en lugar de etiquetas genéricas.", "Mejora la clasificación temática y la llegada a audiencias interesadas."})

	if len(strengths) == 0 {
		strengths = append(strengths, "La cuenta dispone de datos públicos suficientes para establecer una línea base medible.")
	}
	if len(weaknesses) == 0 {
		weaknesses = append(weaknesses, "No se detectan bloqueos graves; el siguiente salto depende de sistematizar experimentos y formatos ganadores.")
	}
	username := report.Profile.Username
	if username == "" {
		username = report.Profile.Name
	}
	report.ExecutiveSummary = fmt.Sprintf("La cuenta @%s obtiene %d/100 (grado %s). Sus últimos %d vídeos promedian %.0f visualizaciones, %.2f%% de interacción sobre visualizaciones y %.2fx visualizaciones por seguidor. El plan prioriza las palancas con mayor margen: gancho, compartidos, cadencia y repetición de patrones ganadores.", username, report.OverallScore, report.OverallGrade, len(report.MediaAnalysis), tm.AverageViews, eng.ViewEngagementRate, tm.ViewsPerFollower)
	report.Strengths, report.Weaknesses, report.Recommendations = strengths, weaknesses, recs
}

type scorePoint struct{ x, y float64 }

func scaleScore(value float64, points []scorePoint) float64 {
	if value <= points[0].x {
		return points[0].y
	}
	for i := 1; i < len(points); i++ {
		if value <= points[i].x {
			p, n := points[i-1], points[i]
			return p.y + (value-p.x)/(n.x-p.x)*(n.y-p.y)
		}
	}
	return points[len(points)-1].y
}
func gradeForScore(score int) (Grade, string, string) {
	if score >= 85 {
		return GradeA, "Cuenta Excelente / Perfecta", "#10B981"
	}
	if score >= 70 {
		return GradeB, "Cuenta Buena / Sólida", "#3B82F6"
	}
	if score >= 50 {
		return GradeD, "Cuenta Decente / Mucho por Mejorar", "#F59E0B"
	}
	return GradeF, "Nivel Muy Bajo / Crítico", "#EF4444"
}
func tikTokBenchmark(rate float64) string {
	if rate >= 8 {
		return "🌟 Sobresaliente por visualización"
	}
	if rate >= 5 {
		return "✅ Fuerte por visualización"
	}
	if rate >= 3 {
		return "➜ Saludable, con margen de optimización"
	}
	if rate >= 1.5 {
		return "⚠️ Moderado por visualización"
	}
	return "🚨 Bajo: pocas visualizaciones se convierten en interacción"
}
func round1(v float64) float64 { return math.Round(v*10) / 10 }
func round2(v float64) float64 { return math.Round(v*100) / 100 }
func medianFloat(values []float64) float64 {
	sort.Float64s(values)
	n := len(values)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return round2((values[n/2-1] + values[n/2]) / 2)
	}
	return round2(values[n/2])
}
func medianInt64(values []int64) float64 {
	slices.Sort(values)
	n := len(values)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return float64(values[n/2-1]+values[n/2]) / 2
	}
	return float64(values[n/2])
}
