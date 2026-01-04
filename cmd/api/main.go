package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EricFreesoul/phoenix-feuer-os/internal/ai/claude"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/ai/openai"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/api/handlers"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/api/routes"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/seo/crawler"
	"github.com/EricFreesoul/phoenix-feuer-os/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	log.Println("🚀 Starting PHOENIX SEO Platform API...")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Port: %s", cfg.Server.Port)

	// Initialize crawler
	crawlerInst := crawler.NewCrawler(
		cfg.SEO.UserAgent,
		cfg.SEO.CrawlTimeout,
		cfg.SEO.MaxCrawlDepth,
	)

	// Initialize AI clients
	var claudeClient *claude.Client
	var openaiClient *openai.Client

	if cfg.AI.ClaudeAPIKey != "" {
		claudeClient = claude.NewClient(cfg.AI.ClaudeAPIKey, cfg.AI.ClaudeModel)
		log.Println("✅ Claude AI client initialized")
	} else {
		log.Println("⚠️  Claude API key not configured")
	}

	if cfg.AI.OpenAIAPIKey != "" {
		openaiClient = openai.NewClient(cfg.AI.OpenAIAPIKey, cfg.AI.OpenAIModel)
		log.Println("✅ OpenAI client initialized")
	} else {
		log.Println("⚠️  OpenAI API key not configured")
	}

	// Initialize handlers
	seoHandler := handlers.NewSEOHandler(crawlerInst, claudeClient, openaiClient)

	// Setup routes
	handler := routes.Setup(seoHandler, cfg.Server.AllowedOrigins)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Server listening on http://localhost:%s", cfg.Server.Port)
		log.Println("📊 API Documentation: http://localhost:" + cfg.Server.Port + "/api/v1/health")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
