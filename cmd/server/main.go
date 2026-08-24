package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reviewmysocialnetworks/internal/config"
	"reviewmysocialnetworks/internal/handlers"
	"syscall"
	"time"
)

func main() {
	config.LoadEnvFile(".env")

	cfg := config.AppConfig
	appID, _, redirectURI, port := cfg.Get()

	mux := http.NewServeMux()

	h := handlers.NewHandler(cfg)
	h.RegisterRoutes(mux)

	webDist := filepath.Join(".", "web", "dist")
	mux.HandleFunc("/", handlers.SPAHandler(webDist))

	handler := handlers.Chain(mux, handlers.Recovery, handlers.Logger, handlers.SecurityHeaders, handlers.CORS(cfg), handlers.RateLimit)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		fmt.Println("================================================================")
		fmt.Println("🚀 ReviewMySocialNetworks - Instagram Analytics & Audit Platform")
		fmt.Println("================================================================")
		fmt.Printf("📡 Servidor HTTP Go iniciado en: http://localhost:%s\n", port)
		fmt.Printf("📷 Instagram App ID:           %s\n", maskString(appID))
		fmt.Printf("🔗 Instagram Redirect URI:     %s\n", redirectURI)
		fmt.Printf("📁 Sirviendo frontend desde:   %s\n", webDist)
		fmt.Println("================================================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error crítico en el servidor HTTP: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Apagando el servidor ReviewMySocialNetworks de forma segura...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Error durante el apagado del servidor: %v", err)
	} else {
		log.Println("✅ Servidor detenido correctamente.")
	}
}

func maskString(s string) string {
	if s == "" {
		return "(no configurado)"
	}
	if len(s) <= 6 {
		return "******"
	}
	return s[:3] + "..." + s[len(s)-3:]
}
