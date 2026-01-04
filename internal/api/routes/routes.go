package routes

import (
	"net/http"

	"github.com/EricFreesoul/phoenix-feuer-os/internal/api/handlers"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/api/middleware"
)

// Setup sets up all API routes
func Setup(seoHandler *handlers.SEOHandler, allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/v1/health", seoHandler.HealthCheck)

	// SEO Analysis endpoints
	mux.HandleFunc("POST /api/v1/seo/analyze", seoHandler.AnalyzeURL)
	mux.HandleFunc("POST /api/v1/seo/keywords/generate", seoHandler.GenerateKeywords)
	mux.HandleFunc("POST /api/v1/seo/meta/optimize", seoHandler.OptimizeMeta)

	// Apply middleware
	handler := middleware.Logger(mux)
	handler = middleware.CORS(allowedOrigins)(handler)

	return handler
}
