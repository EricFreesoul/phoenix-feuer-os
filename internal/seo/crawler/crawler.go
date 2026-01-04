package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// CrawlResult represents the result of crawling a single page
type CrawlResult struct {
	URL              string
	StatusCode       int
	Title            string
	MetaDescription  string
	H1Tags           []string
	H2Tags           []string
	Links            []string
	Images           []Image
	WordCount        int
	LoadTimeMs       int64
	MobileFriendly   bool
	HasHTTPS         bool
	CanonicalURL     string
	Errors           []string
	Headers          map[string]string
	ResponseSize     int64
}

// Image represents an image found on the page
type Image struct {
	Src    string
	Alt    string
	Title  string
	Width  int
	Height int
}

// Crawler handles website crawling
type Crawler struct {
	userAgent          string
	timeout            time.Duration
	maxDepth           int
	respectRobotsTxt   bool
	crawlDelay         time.Duration
	maxConcurrent      int
	client             *http.Client
	visitedURLs        sync.Map
	lastRequestTime    sync.Map
}

// NewCrawler creates a new crawler instance
func NewCrawler(userAgent string, timeout time.Duration, maxDepth int) *Crawler {
	return &Crawler{
		userAgent:        userAgent,
		timeout:          timeout,
		maxDepth:         maxDepth,
		respectRobotsTxt: true,
		crawlDelay:       1 * time.Second,
		maxConcurrent:    5,
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// CrawlPage crawls a single page and returns the result
func (c *Crawler) CrawlPage(ctx context.Context, urlStr string) (*CrawlResult, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Enforce crawl delay
	c.enforceCrawlDelay(parsedURL.Host)

	// Start timing
	startTime := time.Now()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "de-DE,de;q=0.9,en;q=0.8")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	loadTime := time.Since(startTime).Milliseconds()

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract data
	result := &CrawlResult{
		URL:          urlStr,
		StatusCode:   resp.StatusCode,
		LoadTimeMs:   loadTime,
		HasHTTPS:     parsedURL.Scheme == "https",
		Headers:      make(map[string]string),
		Errors:       []string{},
	}

	// Extract headers
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}

	// Parse document
	c.parseNode(doc, result, parsedURL)

	// Check mobile-friendly (simplified check)
	result.MobileFriendly = c.checkMobileFriendly(result)

	// Mark as visited
	c.visitedURLs.Store(urlStr, true)

	return result, nil
}

// parseNode recursively parses HTML nodes
func (c *Crawler) parseNode(n *html.Node, result *CrawlResult, baseURL *url.URL) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				result.Title = strings.TrimSpace(n.FirstChild.Data)
			}
		case "meta":
			c.parseMeta(n, result)
		case "h1":
			text := c.extractText(n)
			if text != "" {
				result.H1Tags = append(result.H1Tags, text)
			}
		case "h2":
			text := c.extractText(n)
			if text != "" {
				result.H2Tags = append(result.H2Tags, text)
			}
		case "a":
			href := c.getAttr(n, "href")
			if href != "" {
				if absURL := c.makeAbsolute(href, baseURL); absURL != "" {
					result.Links = append(result.Links, absURL)
				}
			}
		case "img":
			img := Image{
				Src:   c.getAttr(n, "src"),
				Alt:   c.getAttr(n, "alt"),
				Title: c.getAttr(n, "title"),
			}
			if img.Src != "" {
				img.Src = c.makeAbsolute(img.Src, baseURL)
				result.Images = append(result.Images, img)
			}
		case "link":
			rel := c.getAttr(n, "rel")
			if rel == "canonical" {
				result.CanonicalURL = c.getAttr(n, "href")
			}
		}
	}

	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			words := strings.Fields(text)
			result.WordCount += len(words)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.parseNode(child, result, baseURL)
	}
}

// parseMeta extracts meta tag information
func (c *Crawler) parseMeta(n *html.Node, result *CrawlResult) {
	name := c.getAttr(n, "name")
	property := c.getAttr(n, "property")
	content := c.getAttr(n, "content")

	if name == "description" || property == "og:description" {
		if result.MetaDescription == "" {
			result.MetaDescription = content
		}
	}
}

// extractText extracts all text content from a node
func (c *Crawler) extractText(n *html.Node) string {
	var buf strings.Builder
	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
		}
	}
	extract(n)
	return strings.TrimSpace(buf.String())
}

// getAttr gets an attribute value from a node
func (c *Crawler) getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// makeAbsolute converts relative URLs to absolute
func (c *Crawler) makeAbsolute(href string, baseURL *url.URL) string {
	parsed, err := url.Parse(href)
	if err != nil {
		return ""
	}

	if parsed.IsAbs() {
		return href
	}

	return baseURL.ResolveReference(parsed).String()
}

// checkMobileFriendly performs a basic mobile-friendly check
func (c *Crawler) checkMobileFriendly(result *CrawlResult) bool {
	// Check for viewport meta tag in headers
	viewport := result.Headers["viewport"]
	if viewport != "" {
		return strings.Contains(strings.ToLower(viewport), "width=device-width")
	}
	return false
}

// enforceCrawlDelay ensures we respect crawl delay between requests to the same host
func (c *Crawler) enforceCrawlDelay(host string) {
	if lastTime, ok := c.lastRequestTime.Load(host); ok {
		if last, ok := lastTime.(time.Time); ok {
			elapsed := time.Since(last)
			if elapsed < c.crawlDelay {
				time.Sleep(c.crawlDelay - elapsed)
			}
		}
	}
	c.lastRequestTime.Store(host, time.Now())
}

// IsVisited checks if a URL has been visited
func (c *Crawler) IsVisited(urlStr string) bool {
	_, visited := c.visitedURLs.Load(urlStr)
	return visited
}

// CrawlSite crawls an entire site (multiple pages) - simplified version
func (c *Crawler) CrawlSite(ctx context.Context, startURL string, maxPages int) ([]*CrawlResult, error) {
	results := make([]*CrawlResult, 0)
	queue := []string{startURL}
	visited := make(map[string]bool)

	baseURL, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	for len(queue) > 0 && len(results) < maxPages {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		result, err := c.CrawlPage(ctx, current)
		if err != nil {
			continue
		}

		results = append(results, result)
		visited[current] = true

		// Add internal links to queue
		if len(results) < maxPages {
			for _, link := range result.Links {
				linkURL, err := url.Parse(link)
				if err != nil {
					continue
				}

				// Only follow internal links
				if linkURL.Host == baseURL.Host && !visited[link] {
					queue = append(queue, link)
				}
			}
		}
	}

	return results, nil
}
