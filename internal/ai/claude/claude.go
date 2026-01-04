package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiURL     = "https://api.anthropic.com/v1/messages"
	apiVersion = "2023-06-01"
)

// Client represents a Claude API client
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	maxTokens  int
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents a Claude API request
type Request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

// Response represents a Claude API response
type Response struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// NewClient creates a new Claude API client
func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		maxTokens: 4096,
	}
}

// SendMessage sends a message to Claude and returns the response
func (c *Client) SendMessage(ctx context.Context, systemPrompt string, userMessage string) (*Response, error) {
	request := Request{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    systemPrompt,
		Messages: []Message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}

	return c.sendRequest(ctx, request)
}

// AnalyzeSEO uses Claude to analyze SEO data and provide insights
func (c *Client) AnalyzeSEO(ctx context.Context, seoData map[string]interface{}) (string, error) {
	systemPrompt := `Du bist ein Experte für Suchmaschinenoptimierung (SEO) mit jahrelanger Erfahrung.
	Deine Aufgabe ist es, SEO-Daten zu analysieren und konkrete, umsetzbare Empfehlungen zu geben.

	Analysiere die bereitgestellten SEO-Daten und gebe:
	1. Eine Zusammenfassung der wichtigsten Probleme
	2. Priorisierte Handlungsempfehlungen
	3. Geschätzte Auswirkungen der Verbesserungen
	4. Konkrete Umsetzungsschritte

	Antworte auf Deutsch und sei präzise und handlungsorientiert.`

	dataJSON, err := json.MarshalIndent(seoData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal SEO data: %w", err)
	}

	userMessage := fmt.Sprintf("Bitte analysiere die folgenden SEO-Daten und gebe mir konkrete Empfehlungen:\n\n%s", string(dataJSON))

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return "", err
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return response.Content[0].Text, nil
}

// GenerateContentStrategy generates a content strategy based on keywords
func (c *Client) GenerateContentStrategy(ctx context.Context, keywords []string, niche string) (string, error) {
	systemPrompt := `Du bist ein Content-Strategie-Experte mit Fokus auf SEO.
	Entwickle umfassende Content-Strategien, die auf Keywords und Zielgruppen abgestimmt sind.`

	userMessage := fmt.Sprintf(`Erstelle eine Content-Strategie für die Nische "%s" basierend auf folgenden Keywords:
	%v

	Bitte gebe:
	1. Content-Themen-Vorschläge
	2. Content-Formate (Blog, Video, Infografiken, etc.)
	3. Posting-Frequenz Empfehlung
	4. Interne Verlinkungsstrategie
	5. Content-Cluster Vorschläge`, niche, keywords)

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return "", err
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return response.Content[0].Text, nil
}

// OptimizeMetaTags generates optimized meta tags
func (c *Client) OptimizeMetaTags(ctx context.Context, pageContent string, targetKeywords []string) (map[string]string, error) {
	systemPrompt := `Du bist ein SEO-Experte spezialisiert auf Meta-Tag-Optimierung.
	Erstelle optimierte Meta-Tags (Title und Description) basierend auf dem Seiteninhalt und Ziel-Keywords.

	Wichtig:
	- Title: 50-60 Zeichen
	- Description: 120-160 Zeichen
	- Keywords natürlich einbauen
	- Klick-anregend formulieren`

	userMessage := fmt.Sprintf(`Erstelle optimierte Meta-Tags für folgende Seite:

	Content: %s

	Ziel-Keywords: %v

	Bitte antworte im JSON-Format:
	{
	  "title": "...",
	  "description": "...",
	  "reasoning": "..."
	}`, pageContent, targetKeywords)

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, err
	}

	if len(response.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	// Try to parse JSON response
	var result map[string]string
	if err := json.Unmarshal([]byte(response.Content[0].Text), &result); err != nil {
		// If JSON parsing fails, return raw text
		return map[string]string{
			"raw_response": response.Content[0].Text,
		}, nil
	}

	return result, nil
}

// sendRequest sends a request to the Claude API
func (c *Client) sendRequest(ctx context.Context, request Request) (*Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetTokenUsage returns the token usage from the last request
func (r *Response) GetTokenUsage() (input, output int) {
	return r.Usage.InputTokens, r.Usage.OutputTokens
}
