package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	AI       AIConfig
	SEO      SEOConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	Environment     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// AIConfig holds AI service configuration
type AIConfig struct {
	ClaudeAPIKey      string
	OpenAIAPIKey      string
	ClaudeModel       string
	OpenAIModel       string
	MaxTokensPerMonth int
	EnableCaching     bool
}

// SEOConfig holds SEO-specific configuration
type SEOConfig struct {
	MaxCrawlDepth      int
	CrawlTimeout       time.Duration
	UserAgent          string
	RespectRobotsTxt   bool
	MaxConcurrentCrawls int
	CrawlDelay         time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			AllowedOrigins:  []string{getEnv("ALLOWED_ORIGINS", "*")},
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "phoenix_seo"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		AI: AIConfig{
			ClaudeAPIKey:      getEnv("CLAUDE_API_KEY", ""),
			OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
			ClaudeModel:       getEnv("CLAUDE_MODEL", "claude-3-5-sonnet-20241022"),
			OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-4-turbo-preview"),
			MaxTokensPerMonth: getIntEnv("AI_MAX_TOKENS_PER_MONTH", 1000000),
			EnableCaching:     getBoolEnv("AI_ENABLE_CACHING", true),
		},
		SEO: SEOConfig{
			MaxCrawlDepth:      getIntEnv("SEO_MAX_CRAWL_DEPTH", 3),
			CrawlTimeout:       getDurationEnv("SEO_CRAWL_TIMEOUT", 30*time.Second),
			UserAgent:          getEnv("SEO_USER_AGENT", "PhoenixSEO/1.0 (+https://phoenix-seo.com/bot)"),
			RespectRobotsTxt:   getBoolEnv("SEO_RESPECT_ROBOTS_TXT", true),
			MaxConcurrentCrawls: getIntEnv("SEO_MAX_CONCURRENT_CRAWLS", 5),
			CrawlDelay:         getDurationEnv("SEO_CRAWL_DELAY", 1*time.Second),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
