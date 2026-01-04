package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Client represents a customer/client
type Client struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Name             string     `json:"name" db:"name"`
	Email            string     `json:"email" db:"email"`
	Company          string     `json:"company" db:"company"`
	SubscriptionTier string     `json:"subscription_tier" db:"subscription_tier"`
	Status           string     `json:"status" db:"status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Project represents a domain/website project
type Project struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ClientID     uuid.UUID  `json:"client_id" db:"client_id"`
	Domain       string     `json:"domain" db:"domain"`
	Name         string     `json:"name" db:"name"`
	Status       string     `json:"status" db:"status"`
	LastCrawlAt  *time.Time `json:"last_crawl_at,omitempty" db:"last_crawl_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// SEOAudit represents an SEO audit result
type SEOAudit struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
	AuditType   string     `json:"audit_type" db:"audit_type"`
	Status      string     `json:"status" db:"status"`
	Score       *float64   `json:"score,omitempty" db:"score"`
	Data        JSONB      `json:"data" db:"data"`
	AIInsights  JSONB      `json:"ai_insights,omitempty" db:"ai_insights"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// Keyword represents a tracked keyword
type Keyword struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ProjectID      uuid.UUID  `json:"project_id" db:"project_id"`
	Keyword        string     `json:"keyword" db:"keyword"`
	SearchVolume   *int       `json:"search_volume,omitempty" db:"search_volume"`
	Difficulty     *float64   `json:"difficulty,omitempty" db:"difficulty"`
	CurrentRanking *int       `json:"current_ranking,omitempty" db:"current_ranking"`
	TargetRanking  *int       `json:"target_ranking,omitempty" db:"target_ranking"`
	TrackedSince   time.Time  `json:"tracked_since" db:"tracked_since"`
	LastCheckedAt  *time.Time `json:"last_checked_at,omitempty" db:"last_checked_at"`
}

// Report represents a generated report
type Report struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	ProjectID  uuid.UUID  `json:"project_id" db:"project_id"`
	ReportType string     `json:"report_type" db:"report_type"`
	Title      string     `json:"title" db:"title"`
	Data       JSONB      `json:"data" db:"data"`
	PDFURL     *string    `json:"pdf_url,omitempty" db:"pdf_url"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	SentAt     *time.Time `json:"sent_at,omitempty" db:"sent_at"`
}

// AITask represents an AI processing task
type AITask struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
	TaskType    string     `json:"task_type" db:"task_type"`
	Status      string     `json:"status" db:"status"`
	InputData   JSONB      `json:"input_data" db:"input_data"`
	OutputData  JSONB      `json:"output_data,omitempty" db:"output_data"`
	TokensUsed  *int       `json:"tokens_used,omitempty" db:"tokens_used"`
	Cost        *float64   `json:"cost,omitempty" db:"cost"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// User represents a system user (for authentication)
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Name         string     `json:"name" db:"name"`
	Role         string     `json:"role" db:"role"`
	ClientID     *uuid.UUID `json:"client_id,omitempty" db:"client_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

// AuditTypes
const (
	AuditTypeTechnical   = "technical"
	AuditTypeContent     = "content"
	AuditTypeKeywords    = "keywords"
	AuditTypeBacklinks   = "backlinks"
	AuditTypeCompetitors = "competitors"
	AuditTypeFull        = "full"
)

// AuditStatus
const (
	AuditStatusPending    = "pending"
	AuditStatusProcessing = "processing"
	AuditStatusCompleted  = "completed"
	AuditStatusFailed     = "failed"
)

// SubscriptionTiers
const (
	TierStarter    = "starter"
	TierPro        = "pro"
	TierBusiness   = "business"
	TierEnterprise = "enterprise"
)

// ClientStatus
const (
	ClientStatusActive    = "active"
	ClientStatusInactive  = "inactive"
	ClientStatusSuspended = "suspended"
)

// UserRoles
const (
	RoleAdmin  = "admin"
	RoleClient = "client"
	RoleViewer = "viewer"
)
