package analyzer

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/EricFreesoul/phoenix-feuer-os/internal/seo/crawler"
)

// SEOScore represents the overall SEO score and breakdown
type SEOScore struct {
	Overall      float64            `json:"overall"`
	Technical    float64            `json:"technical"`
	Content      float64            `json:"content"`
	OnPage       float64            `json:"on_page"`
	Performance  float64            `json:"performance"`
	Issues       []Issue            `json:"issues"`
	Opportunities []Opportunity     `json:"opportunities"`
	Breakdown    map[string]float64 `json:"breakdown"`
}

// Issue represents an SEO issue found
type Issue struct {
	Severity    string `json:"severity"` // critical, high, medium, low
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	HowToFix    string `json:"how_to_fix"`
}

// Opportunity represents an SEO improvement opportunity
type Opportunity struct {
	Priority    string  `json:"priority"` // high, medium, low
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Effort      string  `json:"effort"` // low, medium, high
	Potential   float64 `json:"potential"` // potential score improvement
}

// Analyzer performs SEO analysis on crawl results
type Analyzer struct {
	targetKeywords []string
}

// NewAnalyzer creates a new SEO analyzer
func NewAnalyzer(keywords []string) *Analyzer {
	return &Analyzer{
		targetKeywords: keywords,
	}
}

// Analyze performs comprehensive SEO analysis
func (a *Analyzer) Analyze(result *crawler.CrawlResult) *SEOScore {
	score := &SEOScore{
		Issues:        make([]Issue, 0),
		Opportunities: make([]Opportunity, 0),
		Breakdown:     make(map[string]float64),
	}

	// Analyze different aspects
	score.Technical = a.analyzeTechnical(result, score)
	score.Content = a.analyzeContent(result, score)
	score.OnPage = a.analyzeOnPage(result, score)
	score.Performance = a.analyzePerformance(result, score)

	// Calculate overall score (weighted average)
	score.Overall = (score.Technical*0.25 +
		score.Content*0.35 +
		score.OnPage*0.25 +
		score.Performance*0.15)

	return score
}

// analyzeTechnical analyzes technical SEO aspects
func (a *Analyzer) analyzeTechnical(result *crawler.CrawlResult, score *SEOScore) float64 {
	points := 100.0

	// HTTPS check
	if !result.HasHTTPS {
		points -= 15
		score.Issues = append(score.Issues, Issue{
			Severity:    "critical",
			Category:    "security",
			Title:       "Missing HTTPS",
			Description: "Website is not using HTTPS encryption",
			Impact:      "Negative ranking factor and security risk",
			HowToFix:    "Install SSL certificate and redirect all HTTP traffic to HTTPS",
		})
	} else {
		score.Breakdown["https"] = 15
	}

	// Status code check
	if result.StatusCode != 200 {
		points -= 20
		score.Issues = append(score.Issues, Issue{
			Severity:    "critical",
			Category:    "technical",
			Title:       fmt.Sprintf("Non-200 Status Code: %d", result.StatusCode),
			Description: "Page returns an error status code",
			Impact:      "Search engines may not index this page",
			HowToFix:    "Fix server configuration or broken links",
		})
	} else {
		score.Breakdown["status_code"] = 20
	}

	// Canonical URL check
	if result.CanonicalURL == "" {
		points -= 5
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "medium",
			Category:    "technical",
			Title:       "Missing Canonical URL",
			Description: "No canonical link tag found",
			Impact:      "May cause duplicate content issues",
			Effort:      "low",
			Potential:   5,
		})
	} else {
		score.Breakdown["canonical"] = 5
	}

	// Mobile-friendly check
	if !result.MobileFriendly {
		points -= 10
		score.Issues = append(score.Issues, Issue{
			Severity:    "high",
			Category:    "mobile",
			Title:       "Not Mobile-Friendly",
			Description: "Missing or incorrect viewport meta tag",
			Impact:      "Poor mobile experience and ranking penalty",
			HowToFix:    "Add <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">",
		})
	} else {
		score.Breakdown["mobile"] = 10
	}

	return math.Max(0, points)
}

// analyzeContent analyzes content quality
func (a *Analyzer) analyzeContent(result *crawler.CrawlResult, score *SEOScore) float64 {
	points := 100.0

	// Word count check
	if result.WordCount < 300 {
		points -= 20
		score.Issues = append(score.Issues, Issue{
			Severity:    "high",
			Category:    "content",
			Title:       "Thin Content",
			Description: fmt.Sprintf("Only %d words found (recommended: 300+)", result.WordCount),
			Impact:      "May be considered low-quality by search engines",
			HowToFix:    "Add more valuable, relevant content to the page",
		})
	} else if result.WordCount < 600 {
		points -= 10
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "medium",
			Category:    "content",
			Title:       "Expand Content",
			Description: fmt.Sprintf("Page has %d words (recommended: 600+ for better rankings)", result.WordCount),
			Impact:      "More comprehensive content tends to rank better",
			Effort:      "medium",
			Potential:   10,
		})
	} else {
		score.Breakdown["word_count"] = 20
	}

	// H1 tags check
	if len(result.H1Tags) == 0 {
		points -= 15
		score.Issues = append(score.Issues, Issue{
			Severity:    "high",
			Category:    "content",
			Title:       "Missing H1 Tag",
			Description: "No H1 heading found on page",
			Impact:      "H1 is important for SEO and accessibility",
			HowToFix:    "Add a clear, keyword-rich H1 heading",
		})
	} else if len(result.H1Tags) > 1 {
		points -= 5
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "low",
			Category:    "content",
			Title:       "Multiple H1 Tags",
			Description: fmt.Sprintf("Found %d H1 tags (best practice: 1)", len(result.H1Tags)),
			Impact:      "May dilute SEO impact",
			Effort:      "low",
			Potential:   5,
		})
	} else {
		score.Breakdown["h1"] = 15
	}

	// H2 structure
	if len(result.H2Tags) == 0 && result.WordCount > 300 {
		points -= 10
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "medium",
			Category:    "content",
			Title:       "No H2 Headings",
			Description: "Page lacks subheadings for content structure",
			Impact:      "Improves readability and SEO",
			Effort:      "low",
			Potential:   10,
		})
	} else if len(result.H2Tags) > 0 {
		score.Breakdown["h2"] = 10
	}

	// Keyword optimization (if keywords provided)
	if len(a.targetKeywords) > 0 {
		keywordScore := a.analyzeKeywordUsage(result)
		points = points*0.8 + keywordScore*0.2
		score.Breakdown["keywords"] = keywordScore
	}

	return math.Max(0, points)
}

// analyzeOnPage analyzes on-page SEO elements
func (a *Analyzer) analyzeOnPage(result *crawler.CrawlResult, score *SEOScore) float64 {
	points := 100.0

	// Title check
	if result.Title == "" {
		points -= 30
		score.Issues = append(score.Issues, Issue{
			Severity:    "critical",
			Category:    "on_page",
			Title:       "Missing Title Tag",
			Description: "No title tag found",
			Impact:      "Critical for SEO and CTR",
			HowToFix:    "Add a unique, descriptive title tag (50-60 characters)",
		})
	} else {
		titleLen := len(result.Title)
		if titleLen < 30 {
			points -= 10
			score.Issues = append(score.Issues, Issue{
				Severity:    "medium",
				Category:    "on_page",
				Title:       "Title Too Short",
				Description: fmt.Sprintf("Title is %d characters (recommended: 50-60)", titleLen),
				Impact:      "Not utilizing full SERP space",
				HowToFix:    "Expand title to include more relevant keywords",
			})
		} else if titleLen > 60 {
			points -= 5
			score.Opportunities = append(score.Opportunities, Opportunity{
				Priority:    "low",
				Category:    "on_page",
				Title:       "Title Too Long",
				Description: fmt.Sprintf("Title is %d characters (may be truncated)", titleLen),
				Impact:      "May be cut off in search results",
				Effort:      "low",
				Potential:   5,
			})
		} else {
			score.Breakdown["title"] = 30
		}
	}

	// Meta description check
	if result.MetaDescription == "" {
		points -= 20
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "high",
			Category:    "on_page",
			Title:       "Missing Meta Description",
			Description: "No meta description found",
			Impact:      "Missed opportunity to improve CTR",
			Effort:      "low",
			Potential:   20,
		})
	} else {
		descLen := len(result.MetaDescription)
		if descLen < 120 || descLen > 160 {
			points -= 5
			score.Opportunities = append(score.Opportunities, Opportunity{
				Priority:    "medium",
				Category:    "on_page",
				Title:       "Meta Description Length",
				Description: fmt.Sprintf("Description is %d characters (optimal: 120-160)", descLen),
				Impact:      "May be truncated or too short",
				Effort:      "low",
				Potential:   5,
			})
		} else {
			score.Breakdown["meta_description"] = 20
		}
	}

	// Images ALT text check
	missingAlt := 0
	for _, img := range result.Images {
		if img.Alt == "" {
			missingAlt++
		}
	}
	if missingAlt > 0 {
		points -= float64(math.Min(10, float64(missingAlt)))
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "medium",
			Category:    "on_page",
			Title:       "Missing Image ALT Text",
			Description: fmt.Sprintf("%d images without ALT attributes", missingAlt),
			Impact:      "Accessibility and image SEO",
			Effort:      "low",
			Potential:   10,
		})
	} else if len(result.Images) > 0 {
		score.Breakdown["image_alt"] = 10
	}

	return math.Max(0, points)
}

// analyzePerformance analyzes page performance
func (a *Analyzer) analyzePerformance(result *crawler.CrawlResult, score *SEOScore) float64 {
	points := 100.0

	// Load time check
	if result.LoadTimeMs > 3000 {
		points -= 30
		score.Issues = append(score.Issues, Issue{
			Severity:    "high",
			Category:    "performance",
			Title:       "Slow Page Load",
			Description: fmt.Sprintf("Page loads in %dms (target: <3000ms)", result.LoadTimeMs),
			Impact:      "Negative ranking factor and user experience",
			HowToFix:    "Optimize images, enable caching, use CDN, minimize CSS/JS",
		})
	} else if result.LoadTimeMs > 2000 {
		points -= 15
		score.Opportunities = append(score.Opportunities, Opportunity{
			Priority:    "medium",
			Category:    "performance",
			Title:       "Improve Load Time",
			Description: fmt.Sprintf("Page loads in %dms (good but can be better)", result.LoadTimeMs),
			Impact:      "Faster is always better for UX and SEO",
			Effort:      "medium",
			Potential:   15,
		})
	} else {
		score.Breakdown["load_time"] = 30
	}

	return math.Max(0, points)
}

// analyzeKeywordUsage checks keyword presence and density
func (a *Analyzer) analyzeKeywordUsage(result *crawler.CrawlResult) float64 {
	if len(a.targetKeywords) == 0 {
		return 100
	}

	score := 0.0
	maxScore := 100.0
	pointsPerKeyword := maxScore / float64(len(a.targetKeywords))

	content := strings.ToLower(result.Title + " " + result.MetaDescription + " " +
		strings.Join(result.H1Tags, " ") + " " + strings.Join(result.H2Tags, " "))

	for _, keyword := range a.targetKeywords {
		keyword = strings.ToLower(keyword)

		// Check presence in important places
		if strings.Contains(strings.ToLower(result.Title), keyword) {
			score += pointsPerKeyword * 0.4
		}
		if strings.Contains(strings.ToLower(result.MetaDescription), keyword) {
			score += pointsPerKeyword * 0.2
		}
		if a.containsInSlice(result.H1Tags, keyword) {
			score += pointsPerKeyword * 0.3
		}
		if a.containsInSlice(result.H2Tags, keyword) {
			score += pointsPerKeyword * 0.1
		}
	}

	return math.Min(maxScore, score)
}

// containsInSlice checks if any string in slice contains the search term
func (a *Analyzer) containsInSlice(slice []string, search string) bool {
	search = strings.ToLower(search)
	for _, item := range slice {
		if strings.Contains(strings.ToLower(item), search) {
			return true
		}
	}
	return false
}

// CalculateReadability calculates content readability score (Flesch Reading Ease approximation)
func (a *Analyzer) CalculateReadability(text string) float64 {
	sentences := regexp.MustCompile(`[.!?]+`).Split(text, -1)
	words := strings.Fields(text)

	if len(sentences) == 0 || len(words) == 0 {
		return 0
	}

	syllables := a.countSyllables(text)

	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))
	avgSyllablesPerWord := float64(syllables) / float64(len(words))

	// Flesch Reading Ease formula (simplified)
	score := 206.835 - 1.015*avgWordsPerSentence - 84.6*avgSyllablesPerWord

	// Normalize to 0-100
	return math.Max(0, math.Min(100, score))
}

// countSyllables approximates syllable count
func (a *Analyzer) countSyllables(text string) int {
	// Simplified syllable counting (not 100% accurate but good enough)
	vowels := regexp.MustCompile(`[aeiouyäöü]+`)
	matches := vowels.FindAllString(strings.ToLower(text), -1)
	return len(matches)
}
