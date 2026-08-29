package analyzer

import (
	"fmt"
	"math"
	"strings"
)

func calculateScores(
	eng EngagementMetrics,
	cad CadenceMetrics,
	cnt ContentMetrics,
	gro GrowthMetrics,
	profile *Profile,
) (SubScores, int, Grade, string, string) {
	var engScore float64
	rate := eng.AverageEngagementRate
	if rate >= 5.0 {
		engScore = 95.0 + math.Min(5.0, (rate-5.0)*2.0)
	} else if rate >= 3.0 {
		engScore = 85.0 + ((rate-3.0)/2.0)*10.0
	} else if rate >= 1.8 {
		engScore = 70.0 + ((rate-1.8)/1.2)*15.0
	} else if rate >= 0.8 {
		engScore = 50.0 + ((rate-0.8)/1.0)*20.0
	} else {
		engScore = math.Max(10.0, rate*60.0)
	}

	var cadScore float64
	postsPerWeek := cad.EstimatedPostsPerWeek
	daysSinceLast := cad.DaysSinceLastPost

	if postsPerWeek >= 3.0 && daysSinceLast <= 5 {
		cadScore = 95.0
	} else if postsPerWeek >= 2.0 && daysSinceLast <= 7 {
		cadScore = 85.0
	} else if postsPerWeek >= 1.0 && daysSinceLast <= 12 {
		cadScore = 72.0
	} else if postsPerWeek >= 0.5 && daysSinceLast <= 21 {
		cadScore = 55.0
	} else if daysSinceLast > 60 {
		cadScore = 15.0
	} else {
		cadScore = 35.0
	}

	var cntScore float64
	hasReels := cnt.VideoCount > 0
	hasCarousels := cnt.CarouselCount > 0
	hasImages := cnt.ImageCount > 0

	if hasReels && hasCarousels && hasImages {
		cntScore = 95.0
	} else if (hasReels && hasCarousels) || (hasReels && hasImages) || (hasCarousels && hasImages) {
		cntScore = 80.0
	} else if hasCarousels || hasReels {
		cntScore = 65.0
	} else {
		cntScore = 45.0
	}

	var audScore float64
	ratio := gro.FollowerToFollowingRatio
	if ratio >= 10.0 {
		audScore = 95.0
	} else if ratio >= 3.0 {
		audScore = 85.0
	} else if ratio >= 1.2 {
		audScore = 75.0
	} else if ratio >= 0.7 {
		audScore = 55.0
	} else {
		audScore = 30.0
	}

	sEng := clamp(int(math.Round(engScore)), 0, 100)
	sCad := clamp(int(math.Round(cadScore)), 0, 100)
	sCnt := clamp(int(math.Round(cntScore)), 0, 100)
	sAud := clamp(int(math.Round(audScore)), 0, 100)

	subScores := SubScores{
		EngagementScore:     sEng,
		ConsistencyScore:    sCad,
		ContentMixScore:     sCnt,
		AudienceHealthScore: sAud,
	}

	weighted := float64(sEng)*0.35 + float64(sCad)*0.25 + float64(sCnt)*0.20 + float64(sAud)*0.20
	overallScore := clamp(int(math.Round(weighted)), 0, 100)

	var grade Grade
	var title string
	var color string

	if overallScore >= 85 {
		grade = GradeA
		title = "Cuenta Excelente / Perfecta"
		color = "#10B981"
	} else if overallScore >= 70 {
		grade = GradeB
		title = "Cuenta Buena / Sólida"
		color = "#3B82F6"
	} else if overallScore >= 50 {
		grade = GradeD
		title = "Cuenta Decente / Mucho por Mejorar"
		color = "#F59E0B"
	} else {
		grade = GradeF
		title = "Nivel Muy Bajo / Crítico"
		color = "#EF4444"
	}

	return subScores, overallScore, grade, title, color
}

func generateQualitativeInsights(report *AccountReport) {
	var strengths []string
	var weaknesses []string
	var recs []Recommendation

	if report.EngagementMetrics.AverageEngagementRate >= 3.5 {
		strengths = append(strengths, fmt.Sprintf("Excelente tasa de interacción media (%.2f%%), muy por encima del 2.0%% del mercado.", report.EngagementMetrics.AverageEngagementRate))
	} else if report.EngagementMetrics.AverageEngagementRate < 1.2 {
		weaknesses = append(weaknesses, fmt.Sprintf("Baja tasa de engagement (%.2f%%). La audiencia no está interactuando lo suficiente con tus publicaciones.", report.EngagementMetrics.AverageEngagementRate))
		recs = append(recs, Recommendation{
			Category: "Engagement",
			Priority: "Alta",
			Title:    "Incrementar llamados a la acción (CTAs)",
			Action:   "Incluye preguntas abiertas al final de los copies y solicita guardar o compartir la publicación en los primeros 3 segundos de los videos.",
			Impact:   "+35% en comentarios y guardados en 30 días.",
		})
	}

	if report.CadenceMetrics.DaysSinceLastPost <= 3 && report.CadenceMetrics.EstimatedPostsPerWeek >= 2.5 {
		strengths = append(strengths, fmt.Sprintf("Cadencia constante de publicación (~%.1f posts/semana), manteniendo a la comunidad enganchada.", report.CadenceMetrics.EstimatedPostsPerWeek))
	} else if report.CadenceMetrics.DaysSinceLastPost > 20 {
		weaknesses = append(weaknesses, fmt.Sprintf("Llevas %d días sin publicar. Las pausas largas reducen drásticamente la distribución orgánica del algoritmo.", report.CadenceMetrics.DaysSinceLastPost))
		recs = append(recs, Recommendation{
			Category: "Frecuencia",
			Priority: "Alta",
			Title:    "Establecer un calendario editorial mínimo",
			Action:   fmt.Sprintf("Programa al menos 2 a 3 publicaciones semanales en días clave (especialmente los %s a las %02d:00h).", pluralizeDay(report.CadenceMetrics.BestPostingDay), report.CadenceMetrics.BestPostingHour),
			Impact:   "+50% en reactivación de alcance y visibilidad.",
		})
	} else if report.CadenceMetrics.EstimatedPostsPerWeek < 1.0 {
		weaknesses = append(weaknesses, "Frecuencia de publicación baja (menos de 1 post por semana en promedio).")
		recs = append(recs, Recommendation{
			Category: "Frecuencia",
			Priority: "Media",
			Title:    "Aumentar volumen de contenido semanal",
			Action:   "Añade 1 Reel rápido o 1 carrusel informativo adicional por semana para mantener presencia en el feed.",
			Impact:   "+25% en crecimiento de seguidores.",
		})
	}

	if report.ContentMetrics.CarouselCount > 0 && report.ContentMetrics.VideoCount > 0 {
		strengths = append(strengths, "Excelente diversificación de formatos combinando Reels, Carruseles y Fotos fijas.")
	} else if report.ContentMetrics.VideoCount == 0 {
		weaknesses = append(weaknesses, "No estás aprovechando los Reels (0 videos detectados). Actualmente es el formato #1 para alcance no-seguidores en Instagram.")
		recs = append(recs, Recommendation{
			Category: "Formatos",
			Priority: "Alta",
			Title:    "Integrar Reels en tu estrategia de contenido",
			Action:   "Crea clips cortos de 15 a 45 segundos abordando dudas frecuentes, micro-tutoriales o momentos detrás de escena.",
			Impact:   "+200% de alcance hacia cuentas que aún no te siguen.",
		})
	}

	if report.ContentMetrics.CarouselCount == 0 && report.ContentMetrics.TotalAnalyzedPosts > 0 {
		weaknesses = append(weaknesses, "No estás utilizando Carruseles. Los carruseles tienen el mayor ratio de guardados y retención en el feed.")
		recs = append(recs, Recommendation{
			Category: "Formatos",
			Priority: "Media",
			Title:    "Crear Carruseles Educativos / Visuales",
			Action:   "Diseña carruseles de 4 a 7 diapositivas tipo guía paso a paso, listas de recursos o storytelling.",
			Impact:   "+60% en publicaciones guardadas (Saves).",
		})
	}

	if report.GrowthMetrics.FollowerToFollowingRatio >= 3.0 {
		strengths = append(strengths, fmt.Sprintf("Ratio de autoridad saludable (%.1f seguidores por cada cuenta seguida).", report.GrowthMetrics.FollowerToFollowingRatio))
	} else if report.GrowthMetrics.FollowerToFollowingRatio < 0.8 {
		weaknesses = append(weaknesses, fmt.Sprintf("Ratio desbalanceado (sigues a %d cuentas vs %d seguidores). Puede proyectar falta de enfoque o tácticas masivas.", report.Profile.FollowsCount, report.Profile.FollowersCount))
		recs = append(recs, Recommendation{
			Category: "Audiencia",
			Priority: "Baja",
			Title:    "Limpiar cuentas seguidas inactivas",
			Action:   "Revisa tu lista de seguidos y deja de seguir perfiles inactivos o no afines a tu nicho para mejorar tu perfil de marca.",
			Impact:   "Mejora inmediata de la percepción de autoridad del perfil.",
		})
	}

	var summary string
	switch report.OverallGrade {
	case GradeA:
		summary = fmt.Sprintf(
			"¡Enhorabuena! Tu cuenta @%s presenta una salud digital sobresaliente (Puntuación: %d/100). Cuentas con un engagement del %.2f%% y una cadencia ejemplar. Tu formato con mayor tracción es %s. Mantén este ritmo y optimiza los mejores horarios de publicación.",
			report.Profile.Username, report.OverallScore, report.EngagementMetrics.AverageEngagementRate, translateFormat(report.ContentMetrics.BestPerformingType),
		)
	case GradeB:
		summary = fmt.Sprintf(
			"Tu cuenta @%s tiene una base sólida y un rendimiento positivo (Puntuación: %d/100). Tu tasa de interacción media es del %.2f%%. Para dar el salto a la máxima calificación (A), te recomendamos potenciar el formato %s y mejorar la regularidad de publicación.",
			report.Profile.Username, report.OverallScore, report.EngagementMetrics.AverageEngagementRate, translateFormat(report.ContentMetrics.BestPerformingType),
		)
	case GradeD:
		summary = fmt.Sprintf(
			"Tu cuenta @%s se encuentra en un estado decente pero tiene margen considerable de optimización (Puntuación: %d/100). Tu engagement actual (%.2f%%) o tu frecuencia de publicación necesitan una intervención estratégica para evitar el estancamiento de tu comunidad.",
			report.Profile.Username, report.OverallScore, report.EngagementMetrics.AverageEngagementRate,
		)
	case GradeF:
		summary = fmt.Sprintf(
			"Tu cuenta @%s se encuentra en un nivel crítico (Puntuación: %d/100). La baja interacción (%.2f%%) o la inactividad prolongada (%d días) indican que la audiencia está desconectada. Sigue el plan de rescate recomendado a continuación.",
			report.Profile.Username, report.OverallScore, report.EngagementMetrics.AverageEngagementRate, report.CadenceMetrics.DaysSinceLastPost,
		)
	}

	if len(strengths) == 0 {
		strengths = append(strengths, "La cuenta cuenta con perfil e historial analizable listo para ser optimizado.")
	}
	if len(weaknesses) == 0 {
		weaknesses = append(weaknesses, "No se han detectado anomalías graves en las publicaciones analizadas.")
	}
	if len(recs) == 0 {
		recs = append(recs, Recommendation{
			Category: "Crecimiento",
			Priority: "Media",
			Title:    "Explorar colaboraciones y directos",
			Action:   "Realiza publicaciones conjuntas (Co-author posts) con cuentas afines de tu sector.",
			Impact:   "+20% de expansión hacia nuevas audiencias cualificadas.",
		})
	}

	report.ExecutiveSummary = summary
	report.Strengths = strengths
	report.Weaknesses = weaknesses
	report.Recommendations = recs
}

func translateFormat(format string) string {
	switch format {
	case "VIDEO":
		return "Reels / Vídeos"
	case "CAROUSEL_ALBUM":
		return "Carruseles"
	case "IMAGE":
		return "Fotografías individuales"
	default:
		return "Publicaciones"
	}
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func pluralizeDay(day string) string {
	day = strings.TrimSpace(day)
	if day == "" {
		return "días recomendados"
	}
	if strings.HasSuffix(strings.ToLower(day), "s") {
		return day
	}
	return day + "s"
}
