package main

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	config.LoadEnvFile(".env")

	cfg := config.AppConfig
	if err := cfg.Validate(); err != nil {
		slog.Error("invalid_configuration", "error", err)
		os.Exit(1)
	}
	appID, _, redirectURI, port := cfg.Get()

	mux := http.NewServeMux()

	h := handlers.NewHandler(cfg)
	h.RegisterRoutes(mux)

	webDist := filepath.Join(".", "web", "dist")
	mux.HandleFunc("/", handlers.SPAHandler(webDist))

	handler := handlers.Chain(mux, handlers.Compression, handlers.Logger, handlers.Recovery, handlers.SecurityHeaders, handlers.CORS(cfg), handlers.RateLimit(cfg))

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
			slog.Error("http_server_failure", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("http_server_shutdown_started")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("http_server_shutdown_failure", "error", err)
	} else {
		slog.Info("http_server_shutdown_complete")
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
