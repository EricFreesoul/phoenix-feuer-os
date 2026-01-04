package openai

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
	apiURL = "https://api.openai.com/v1/chat/completions"
)

// Client represents an OpenAI API client
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents an OpenAI API request
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Response represents an OpenAI API response
type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewClient creates a new OpenAI API client
func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// SendMessage sends a message to OpenAI and returns the response
func (c *Client) SendMessage(ctx context.Context, systemPrompt, userMessage string) (*Response, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}

	request := Request{
		Model:       c.model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	return c.sendRequest(ctx, request)
}

// GenerateKeywords generates keyword suggestions
func (c *Client) GenerateKeywords(ctx context.Context, topic string, count int) ([]string, error) {
	systemPrompt := `Du bist ein SEO-Keyword-Experte. Generiere relevante Keywords für gegebene Themen.
	Fokussiere auf Keywords mit gutem Suchvolumen und moderater Konkurrenz.`

	userMessage := fmt.Sprintf(`Generiere %d relevante SEO-Keywords für das Thema: "%s"

	Bitte gebe die Keywords als JSON-Array zurück:
	["keyword1", "keyword2", "keyword3", ...]`, count, topic)

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	content := response.Choices[0].Message.Content

	// Try to parse JSON
	var keywords []string
	if err := json.Unmarshal([]byte(content), &keywords); err != nil {
		// If not valid JSON, split by newlines and clean up
		lines := splitAndClean(content)
		return lines, nil
	}

	return keywords, nil
}

// GenerateMetaDescription generates an optimized meta description
func (c *Client) GenerateMetaDescription(ctx context.Context, pageContent string, keywords []string) (string, error) {
	systemPrompt := `Du bist ein SEO-Copywriting-Experte. Erstelle optimierte Meta-Descriptions.

	Anforderungen:
	- 120-160 Zeichen
	- Keywords natürlich einbauen
	- Call-to-Action einbeziehen
	- Klick-anregend formulieren`

	userMessage := fmt.Sprintf(`Erstelle eine optimierte Meta-Description für:

	Content: %s
	Keywords: %v

	Gebe NUR die Meta-Description zurück, ohne weitere Erklärungen.`, pageContent, keywords)

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

// GenerateContentIdeas generates content ideas based on keywords
func (c *Client) GenerateContentIdeas(ctx context.Context, keywords []string, niche string) ([]string, error) {
	systemPrompt := `Du bist ein Content-Marketing-Experte. Entwickle kreative Content-Ideen für SEO.`

	userMessage := fmt.Sprintf(`Erstelle 10 konkrete Content-Ideen (Blog-Titel) für die Niche "%s" basierend auf diesen Keywords:
	%v

	Gebe die Ideen als JSON-Array zurück:
	["Idee 1", "Idee 2", ...]`, niche, keywords)

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	content := response.Choices[0].Message.Content

	// Try to parse JSON
	var ideas []string
	if err := json.Unmarshal([]byte(content), &ideas); err != nil {
		// If not valid JSON, split by newlines
		lines := splitAndClean(content)
		return lines, nil
	}

	return ideas, nil
}

// AnalyzeCompetitor analyzes competitor content
func (c *Client) AnalyzeCompetitor(ctx context.Context, competitorURL, yourURL string) (string, error) {
	systemPrompt := `Du bist ein Competitive-Analysis-Experte im SEO-Bereich.
	Analysiere Wettbewerber und gebe konkrete Empfehlungen.`

	userMessage := fmt.Sprintf(`Vergleiche folgende URLs und gebe Empfehlungen:

	Konkurrent: %s
	Eigene Seite: %s

	Was macht der Konkurrent besser? Was können wir verbessern?`, competitorURL, yourURL)

	response, err := c.SendMessage(ctx, systemPrompt, userMessage)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

// sendRequest sends a request to OpenAI API
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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

// splitAndClean splits content by newlines and cleans up
func splitAndClean(content string) []string {
	lines := []string{}
	for _, line := range bytes.Split([]byte(content), []byte("\n")) {
		cleaned := bytes.TrimSpace(line)
		if len(cleaned) > 0 {
			lines = append(lines, string(cleaned))
		}
	}
	return lines
}

// GetTokenUsage returns token usage from response
func (r *Response) GetTokenUsage() (prompt, completion, total int) {
	return r.Usage.PromptTokens, r.Usage.CompletionTokens, r.Usage.TotalTokens
}
