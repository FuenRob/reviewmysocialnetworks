package instagram

import "time"

type MockAccountData struct {
	Profile UserProfile `json:"profile"`
	Media   []MediaItem `json:"media"`
}

func GetMockAccount(tier string) (*UserProfile, []MediaItem) {
	now := time.Now()

	switch tier {
	case "A", "perfect", "excelente":
		profile := UserProfile{
			ID:                "mock_user_a",
			Username:          "codeandcoffee.dev",
			Name:              "Alex Rivera | Code & Cloud 🚀",
			Biography:         "💻 Full-stack tips, Go & React mastery\n🚀 Ayudo a devs a pasar de Jr a Sr\n📦 100+ recursos gratis en el link 👇",
			ProfilePictureURL: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=300&auto=format&fit=crop&q=80",
			FollowersCount:    48500,
			FollowsCount:      340,
			MediaCount:        284,
			Website:           "https://codeandcoffee.dev/resources",
			AccountType:       "BUSINESS",
		}

		media := []MediaItem{
			{
				ID:               "media_a_1",
				Caption:          "5 Patrones de Concurrencia en Go que TODO desarrollador debería dominar en 2025 ⚡ ¿Cuál usas más? 👇 Guarda este post!",
				MediaType:        "CAROUSEL_ALBUM",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_a_1",
				Timestamp:        now.Add(-24 * time.Hour),
				LikeCount:        3420,
				CommentsCount:    284,
				Insights:         &MediaInsights{Impressions: 48200, Reach: 39500, Saved: 1840, Engagement: 3704},
			},
			{
				ID:               "media_a_2",
				Caption:          "Cómo optimizar tus queries SQL de 12s a 45ms en 3 pasos simples 🔥🚀 #backend #golang #sql",
				MediaType:        "VIDEO",
				MediaProductType: "REELS",
				MediaURL:         "https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=800&auto=format&fit=crop&q=80",
				ThumbnailURL:     "https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_a_2",
				Timestamp:        now.Add(-3 * 24 * time.Hour),
				LikeCount:        4210,
				CommentsCount:    390,
				Insights:         &MediaInsights{Impressions: 62000, Reach: 51200, Saved: 2900, Engagement: 4600, VideoViews: 45000},
			},
			{
				ID:               "media_a_3",
				Caption:          "Mi setup minimalista para programar 8 horas al día sin dolor de espalda 💻🎧 ¿Qué opinan de las pantallas verticales?",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1593062096033-9a26b09da705?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_a_3",
				Timestamp:        now.Add(-5 * 24 * time.Hour),
				LikeCount:        2980,
				CommentsCount:    195,
				Insights:         &MediaInsights{Impressions: 38000, Reach: 31000, Saved: 820, Engagement: 3175},
			},
			{
				ID:               "media_a_4",
				Caption:          "Roadmap completo para convertirte en Backend Engineer de alto impacto en 2025 🗺️ Etiqueta a tu compa de código!",
				MediaType:        "CAROUSEL_ALBUM",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_a_4",
				Timestamp:        now.Add(-7 * 24 * time.Hour),
				LikeCount:        3890,
				CommentsCount:    310,
				Insights:         &MediaInsights{Impressions: 54000, Reach: 44000, Saved: 2450, Engagement: 4200},
			},
			{
				ID:               "media_a_5",
				Caption:          "Microservicios vs Monolitos modulares: La verdad que nadie te dice en las conferencias 🎯 #softwareengineer",
				MediaType:        "VIDEO",
				MediaProductType: "REELS",
				MediaURL:         "https://images.unsplash.com/photo-1504639725590-34d0984388bd?w=800&auto=format&fit=crop&q=80",
				ThumbnailURL:     "https://images.unsplash.com/photo-1504639725590-34d0984388bd?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_a_5",
				Timestamp:        now.Add(-10 * 24 * time.Hour),
				LikeCount:        4800,
				CommentsCount:    440,
				Insights:         &MediaInsights{Impressions: 71000, Reach: 59000, Saved: 3600, Engagement: 5240, VideoViews: 58000},
			},
			{
				ID:               "media_a_6",
				Caption:          "Top 7 Extensiones de VS Code que triplican tu velocidad de desarrollo 🔥",
				MediaType:        "CAROUSEL_ALBUM",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_a_6",
				Timestamp:        now.Add(-12 * 24 * time.Hour),
				LikeCount:        3650,
				CommentsCount:    260,
				Insights:         &MediaInsights{Impressions: 49000, Reach: 41000, Saved: 2100, Engagement: 3910},
			},
		}
		return &profile, media

	case "B", "good", "buena":
		profile := UserProfile{
			ID:                "mock_user_b",
			Username:          "atelier.cafe.madrid",
			Name:              "Atelier Café de Especialidad ☕",
			Biography:         "📍 Calle del Sol 14, Madrid\n🌿 Granos de origen único tostados artesanalmente\n🥐 Brunch diario 9:00 - 16:00",
			ProfilePictureURL: "https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=300&auto=format&fit=crop&q=80",
			FollowersCount:    19200,
			FollowsCount:      1420,
			MediaCount:        156,
			Website:           "https://ateliercafemadrid.es",
			AccountType:       "BUSINESS",
		}

		media := []MediaItem{
			{
				ID:               "media_b_1",
				Caption:          "Llegó nueva cosecha de Etiopía Yirgacheffe 🌸 Notas florales y a jazmín que enamoran. Ven a probarlo en filtro V60.",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_b_1",
				Timestamp:        now.Add(-6 * 24 * time.Hour),
				LikeCount:        420,
				CommentsCount:    18,
				Insights:         &MediaInsights{Impressions: 5900, Reach: 4800, Saved: 65, Engagement: 438},
			},
			{
				ID:               "media_b_2",
				Caption:          "El secreto detrás de nuestro crujiente croissant de pistacho 🥐✨ Horneado cada mañana a las 7am.",
				MediaType:        "CAROUSEL_ALBUM",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1555396273-367ea4eb4db5?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_b_2",
				Timestamp:        now.Add(-14 * 24 * time.Hour),
				LikeCount:        580,
				CommentsCount:    32,
				Insights:         &MediaInsights{Impressions: 7500, Reach: 6200, Saved: 140, Engagement: 612},
			},
			{
				ID:               "media_b_3",
				Caption:          "Tardes de lluvia y café caliente en nuestro rincón favorito 🌧️☕",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_b_3",
				Timestamp:        now.Add(-21 * 24 * time.Hour),
				LikeCount:        380,
				CommentsCount:    12,
				Insights:         &MediaInsights{Impressions: 4800, Reach: 3900, Saved: 45, Engagement: 392},
			},
			{
				ID:               "media_b_4",
				Caption:          "Workshop de cata de cafés este sábado: últimas 3 plazas disponibles 🎟️ Enlace en bio para reservar.",
				MediaType:        "CAROUSEL_ALBUM",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1447933601403-0c6688de566e?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_b_4",
				Timestamp:        now.Add(-30 * 24 * time.Hour),
				LikeCount:        340,
				CommentsCount:    15,
				Insights:         &MediaInsights{Impressions: 4200, Reach: 3400, Saved: 50, Engagement: 355},
			},
		}
		return &profile, media

	case "D", "decent", "decente":
		profile := UserProfile{
			ID:                "mock_user_d",
			Username:          "fitness_routine_daily",
			Name:              "Fitness Routine & Motivation",
			Biography:         "Fitness quotes & workout motivation 💪\nDM for promo | Follow for more",
			ProfilePictureURL: "https://images.unsplash.com/photo-1517838277536-f5f99be501cd?w=300&auto=format&fit=crop&q=80",
			FollowersCount:    14500,
			FollowsCount:      3100,
			MediaCount:        74,
			Website:           "",
			AccountType:       "PERSONAL",
		}

		media := []MediaItem{
			{
				ID:               "media_d_1",
				Caption:          "No pain no gain 🔥 #fitness #gym #workout #motivation #legday",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_d_1",
				Timestamp:        now.Add(-6 * 24 * time.Hour),
				LikeCount:        130,
				CommentsCount:    4,
				Insights:         &MediaInsights{Impressions: 1900, Reach: 1500, Saved: 12, Engagement: 134},
			},
			{
				ID:               "media_d_2",
				Caption:          "Monday mood 🏋️‍♂️ Mantén el enfoque en tus metas.",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1581009146145-b5ef050c2e1e?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_d_2",
				Timestamp:        now.Add(-21 * 24 * time.Hour),
				LikeCount:        115,
				CommentsCount:    2,
				Insights:         &MediaInsights{Impressions: 1600, Reach: 1200, Saved: 8, Engagement: 117},
			},
			{
				ID:               "media_d_3",
				Caption:          "Suplementos básicos: Creatina y Proteína. ¿Cuál tomas tú?",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1574680096145-d05b474e2155?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_d_3",
				Timestamp:        now.Add(-48 * 24 * time.Hour),
				LikeCount:        98,
				CommentsCount:    3,
				Insights:         &MediaInsights{Impressions: 1400, Reach: 1050, Saved: 5, Engagement: 101},
			},
		}
		return &profile, media

	case "F", "poor", "bajo", "critica":
		fallthrough
	default:
		profile := UserProfile{
			ID:                "mock_user_f",
			Username:          "crypto_signals_daily_official",
			Name:              "CRYPTO SIGNALS 📈💰 100x GEMS",
			Biography:         "🚀 98% Winrate Daily VIP Signals\n💎 Join Telegram VIP (link)\n🚫 No financial advice",
			ProfilePictureURL: "https://images.unsplash.com/photo-1622979135225-d2ba269bc1df?w=300&auto=format&fit=crop&q=80",
			FollowersCount:    32400,
			FollowsCount:      7450,
			MediaCount:        38,
			Website:           "https://t.me/fake_signals_channel",
			AccountType:       "PERSONAL",
		}

		media := []MediaItem{
			{
				ID:               "media_f_1",
				Caption:          "NEW 100x GEM ALERT 🚀🚀 Check link in bio now before pump! #crypto #bitcoin #altcoins",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1621416894569-0f39ed31d247?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_f_1",
				Timestamp:        now.Add(-42 * 24 * time.Hour),
				LikeCount:        24,
				CommentsCount:    1,
				Insights:         &MediaInsights{Impressions: 400, Reach: 310, Saved: 1, Engagement: 25},
			},
			{
				ID:               "media_f_2",
				Caption:          "Bitcoin to $150k soon? Drop your thoughts below 👇",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1518770660439-4636190af475?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_f_2",
				Timestamp:        now.Add(-115 * 24 * time.Hour),
				LikeCount:        18,
				CommentsCount:    0,
				Insights:         &MediaInsights{Impressions: 280, Reach: 210, Saved: 0, Engagement: 18},
			},
			{
				ID:               "media_f_3",
				Caption:          "VIP group made +450% today! Don't miss out.",
				MediaType:        "IMAGE",
				MediaProductType: "FEED",
				MediaURL:         "https://images.unsplash.com/photo-1642543492481-44e81e3914a7?w=800&auto=format&fit=crop&q=80",
				Permalink:        "https://instagram.com/p/mock_f_3",
				Timestamp:        now.Add(-180 * 24 * time.Hour),
				LikeCount:        12,
				CommentsCount:    1,
				Insights:         &MediaInsights{Impressions: 220, Reach: 180, Saved: 0, Engagement: 13},
			},
		}
		return &profile, media
	}
}
