package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EricFreesoul/phoenix-feuer-os/internal/ai/claude"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/ai/openai"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/seo/analyzer"
	"github.com/EricFreesoul/phoenix-feuer-os/internal/seo/crawler"
)

// SEOHandler handles SEO-related API endpoints
type SEOHandler struct {
	crawler     *crawler.Crawler
	claudeClient *claude.Client
	openaiClient *openai.Client
}

// NewSEOHandler creates a new SEO handler
func NewSEOHandler(crawlerInst *crawler.Crawler, claudeClient *claude.Client, openaiClient *openai.Client) *SEOHandler {
	return &SEOHandler{
		crawler:     crawlerInst,
		claudeClient: claudeClient,
		openaiClient: openaiClient,
	}
}

// AnalyzeURLRequest represents a request to analyze a URL
type AnalyzeURLRequest struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords,omitempty"`
	UseAI    bool     `json:"use_ai"`
}

// AnalyzeURLResponse represents the response of URL analysis
type AnalyzeURLResponse struct {
	CrawlResult *crawler.CrawlResult `json:"crawl_result"`
	SEOScore    *analyzer.SEOScore   `json:"seo_score"`
	AIInsights  string               `json:"ai_insights,omitempty"`
	AnalyzedAt  time.Time            `json:"analyzed_at"`
}

// AnalyzeURL handles POST /api/v1/seo/analyze
func (h *SEOHandler) AnalyzeURL(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Crawl the URL
	crawlResult, err := h.crawler.CrawlPage(ctx, req.URL)
	if err != nil {
		http.Error(w, "Failed to crawl URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Analyze SEO
	seoAnalyzer := analyzer.NewAnalyzer(req.Keywords)
	seoScore := seoAnalyzer.Analyze(crawlResult)

	response := AnalyzeURLResponse{
		CrawlResult: crawlResult,
		SEOScore:    seoScore,
		AnalyzedAt:  time.Now(),
	}

	// Generate AI insights if requested
	if req.UseAI && h.claudeClient != nil {
		seoData := map[string]interface{}{
			"url":        req.URL,
			"score":      seoScore,
			"crawl_data": crawlResult,
		}

		insights, err := h.claudeClient.AnalyzeSEO(ctx, seoData)
		if err == nil {
			response.AIInsights = insights
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerateKeywordsRequest represents keyword generation request
type GenerateKeywordsRequest struct {
	Topic string `json:"topic"`
	Count int    `json:"count"`
}

// GenerateKeywords handles POST /api/v1/seo/keywords/generate
func (h *SEOHandler) GenerateKeywords(w http.ResponseWriter, r *http.Request) {
	var req GenerateKeywordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}

	if req.Count == 0 {
		req.Count = 10
	}

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	keywords, err := h.openaiClient.GenerateKeywords(ctx, req.Topic, req.Count)
	if err != nil {
		http.Error(w, "Failed to generate keywords: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"topic":    req.Topic,
		"keywords": keywords,
		"count":    len(keywords),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OptimizeMetaRequest represents meta tag optimization request
type OptimizeMetaRequest struct {
	Content  string   `json:"content"`
	Keywords []string `json:"keywords"`
}

// OptimizeMeta handles POST /api/v1/seo/meta/optimize
func (h *SEOHandler) OptimizeMeta(w http.ResponseWriter, r *http.Request) {
	var req OptimizeMetaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	metaTags, err := h.claudeClient.OptimizeMetaTags(ctx, req.Content, req.Keywords)
	if err != nil {
		http.Error(w, "Failed to optimize meta tags: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metaTags)
}

// HealthCheck handles GET /api/v1/health
func (h *SEOHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "phoenix-seo-api",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
