-- Migration: 001_initial_schema
-- Created: 2026-01-04
-- Description: Initial database schema for Phoenix SEO Platform

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search optimization

-- ============================================================================
-- USERS & AUTHENTICATION
-- ============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'client' NOT NULL,
    client_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login_at TIMESTAMP,
    CONSTRAINT valid_role CHECK (role IN ('admin', 'client', 'viewer'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- ============================================================================
-- CLIENTS/CUSTOMERS
-- ============================================================================

CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    company VARCHAR(255),
    subscription_tier VARCHAR(50) DEFAULT 'starter',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT valid_subscription CHECK (subscription_tier IN ('starter', 'pro', 'business', 'enterprise')),
    CONSTRAINT valid_status CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_status ON clients(status);
CREATE INDEX idx_clients_subscription ON clients(subscription_tier);

-- Add foreign key to users
ALTER TABLE users ADD CONSTRAINT fk_users_client
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE;

-- ============================================================================
-- PROJECTS/DOMAINS
-- ============================================================================

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    last_crawl_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT valid_project_status CHECK (status IN ('active', 'paused', 'archived'))
);

CREATE INDEX idx_projects_client_id ON projects(client_id);
CREATE INDEX idx_projects_domain ON projects(domain);
CREATE INDEX idx_projects_status ON projects(status);

-- ============================================================================
-- SEO AUDITS
-- ============================================================================

CREATE TABLE seo_audits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    audit_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    score DECIMAL(5,2),
    data JSONB,
    ai_insights JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    CONSTRAINT valid_audit_type CHECK (audit_type IN ('technical', 'content', 'keywords', 'backlinks', 'competitors', 'full')),
    CONSTRAINT valid_audit_status CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

CREATE INDEX idx_audits_project_id ON seo_audits(project_id);
CREATE INDEX idx_audits_type ON seo_audits(audit_type);
CREATE INDEX idx_audits_status ON seo_audits(status);
CREATE INDEX idx_audits_created_at ON seo_audits(created_at DESC);
CREATE INDEX idx_audits_data ON seo_audits USING GIN (data);

-- ============================================================================
-- KEYWORDS
-- ============================================================================

CREATE TABLE keywords (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL,
    search_volume INTEGER,
    difficulty DECIMAL(5,2),
    current_ranking INTEGER,
    target_ranking INTEGER,
    tracked_since TIMESTAMP DEFAULT NOW(),
    last_checked_at TIMESTAMP,
    CONSTRAINT valid_ranking CHECK (current_ranking >= 0 AND current_ranking <= 100),
    CONSTRAINT valid_target CHECK (target_ranking >= 0 AND target_ranking <= 100),
    CONSTRAINT valid_difficulty CHECK (difficulty >= 0 AND difficulty <= 100)
);

CREATE INDEX idx_keywords_project_id ON keywords(project_id);
CREATE INDEX idx_keywords_keyword ON keywords USING GIN (keyword gin_trgm_ops);
CREATE INDEX idx_keywords_ranking ON keywords(current_ranking);

-- ============================================================================
-- KEYWORD HISTORY (for tracking ranking changes)
-- ============================================================================

CREATE TABLE keyword_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    keyword_id UUID NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    ranking INTEGER NOT NULL,
    search_volume INTEGER,
    checked_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_keyword_history_keyword_id ON keyword_history(keyword_id);
CREATE INDEX idx_keyword_history_checked_at ON keyword_history(checked_at DESC);

-- ============================================================================
-- REPORTS
-- ============================================================================

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    report_type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    data JSONB,
    pdf_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    sent_at TIMESTAMP,
    CONSTRAINT valid_report_type CHECK (report_type IN ('weekly', 'monthly', 'quarterly', 'custom'))
);

CREATE INDEX idx_reports_project_id ON reports(project_id);
CREATE INDEX idx_reports_type ON reports(report_type);
CREATE INDEX idx_reports_created_at ON reports(created_at DESC);

-- ============================================================================
-- AI TASKS
-- ============================================================================

CREATE TABLE ai_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    task_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    input_data JSONB,
    output_data JSONB,
    tokens_used INTEGER,
    cost DECIMAL(10,4),
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    CONSTRAINT valid_ai_status CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

CREATE INDEX idx_ai_tasks_project_id ON ai_tasks(project_id);
CREATE INDEX idx_ai_tasks_status ON ai_tasks(status);
CREATE INDEX idx_ai_tasks_created_at ON ai_tasks(created_at DESC);

-- ============================================================================
-- PAGE AUDITS (detailed page-level analysis)
-- ============================================================================

CREATE TABLE page_audits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    seo_audit_id UUID NOT NULL REFERENCES seo_audits(id) ON DELETE CASCADE,
    url VARCHAR(2048) NOT NULL,
    title VARCHAR(500),
    meta_description VARCHAR(1000),
    h1 TEXT[],
    word_count INTEGER,
    load_time_ms INTEGER,
    mobile_friendly BOOLEAN,
    issues JSONB,
    recommendations JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_page_audits_audit_id ON page_audits(seo_audit_id);
CREATE INDEX idx_page_audits_url ON page_audits(url);

-- ============================================================================
-- BACKLINKS
-- ============================================================================

CREATE TABLE backlinks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    source_url VARCHAR(2048) NOT NULL,
    target_url VARCHAR(2048) NOT NULL,
    anchor_text VARCHAR(500),
    domain_authority INTEGER,
    is_dofollow BOOLEAN DEFAULT true,
    first_seen TIMESTAMP DEFAULT NOW(),
    last_seen TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'active',
    CONSTRAINT valid_backlink_status CHECK (status IN ('active', 'lost', 'broken'))
);

CREATE INDEX idx_backlinks_project_id ON backlinks(project_id);
CREATE INDEX idx_backlinks_status ON backlinks(status);
CREATE INDEX idx_backlinks_domain_authority ON backlinks(domain_authority DESC);

-- ============================================================================
-- API KEYS (for client API access)
-- ============================================================================

CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    permissions JSONB,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

CREATE INDEX idx_api_keys_client_id ON api_keys(client_id);
CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);

-- ============================================================================
-- AUDIT LOG (for compliance & debugging)
-- ============================================================================

CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,
    changes JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_entity ON audit_log(entity_type, entity_id);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at DESC);

-- ============================================================================
-- TRIGGERS (auto-update updated_at timestamps)
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON clients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- INITIAL DATA (optional)
-- ============================================================================

-- Insert default admin user (password: change_me_immediately)
-- Password hash for "change_me_immediately" using bcrypt
INSERT INTO users (email, password_hash, name, role)
VALUES ('admin@phoenix-seo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin User', 'admin');
