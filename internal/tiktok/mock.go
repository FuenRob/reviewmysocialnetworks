package tiktok

import "time"

func GetMockAccount(tier string) (*UserProfile, []Video) {
	now := time.Now()
	type preset struct {
		username, name, bio                    string
		followers, following, likes, videos    int
		views, likesPerVideo, comments, shares int
		gapDays, ageDays, duration             int
		verified                               bool
	}
	presets := map[string]preset{
		"A": {"cienciaen60s", "Clara | Ciencia en 60s", "Experimentos, espacio y curiosidades explicadas sin humo. Nuevo vídeo cada día.", 120000, 420, 6800000, 438, 260000, 21000, 1800, 2600, 1, 1, 32, true},
		"B": {"cocinaconmarta", "Marta cocina fácil", "Recetas sencillas, económicas y listas en menos de 30 minutos.", 52000, 900, 740000, 192, 43000, 1500, 130, 180, 3, 2, 41, false},
		"D": {"entrenaencasa", "Entrena en casa", "Rutinas y motivación fitness.", 28000, 5100, 120000, 86, 7000, 115, 18, 10, 7, 9, 24, false},
		"F": {"ofertas_flash_ya", "OFERTAS FLASH", "Sígueme para más.", 46000, 70000, 90000, 41, 3500, 14, 2, 1, 35, 74, 12, false},
	}
	p, ok := presets[tier]
	if !ok {
		p = presets["F"]
	}
	profile := &UserProfile{
		OpenID: "mock_tiktok_" + tier, Username: p.username, DisplayName: p.name,
		BioDescription: p.bio, AvatarURL: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=300&auto=format&fit=crop&q=80",
		ProfileDeepLink: "https://www.tiktok.com/@" + p.username, IsVerified: p.verified,
		FollowerCount: p.followers, FollowingCount: p.following, LikesCount: p.likes, VideoCount: p.videos,
	}
	count := 8
	items := make([]Video, 0, count)
	for i := range count {
		factor := 1.0 + float64((i%3)-1)*0.12
		created := now.Add(-time.Duration(p.ageDays+i*p.gapDays) * 24 * time.Hour)
		items = append(items, Video{
			ID: "tt_" + tier + string(rune('a'+i)), CreateTime: created.Unix(),
			CoverImageURL:    "https://images.unsplash.com/photo-1531297484001-80022131f5a1?w=800&auto=format&fit=crop&q=80",
			ShareURL:         "https://www.tiktok.com/@" + p.username + "/video/mock" + string(rune('1'+i)),
			VideoDescription: []string{"El error que casi todo el mundo comete y cómo corregirlo", "Tres pasos prácticos que puedes probar hoy", "La explicación sencilla que me habría gustado conocer antes", "¿Mito o realidad? Lo comprobamos", "Guarda esta guía para consultarla después", "Parte 2 de la serie: el detalle que cambia el resultado", "Respondo a la pregunta más repetida", "El resultado final y lo que aprendimos"}[i],
			Duration:         p.duration + (i%3-1)*5, ViewCount: int64(float64(p.views) * factor),
			LikeCount: int(float64(p.likesPerVideo) * factor), CommentCount: int(float64(p.comments) * factor), ShareCount: int(float64(p.shares) * factor),
		})
	}
	return profile, items
}
