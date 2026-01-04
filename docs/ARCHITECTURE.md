# PHOENIX SEO Platform - Architektur-Dokumentation

## Strategisches Über-Ziel

**Mission:** Positionierung als führende SEO-Beratungsinstanz mit KI-Vorsprung

**Zeitrahmen:**
- **Start:** Februar 2026
- **Ziel:** Q2 2029 - Vollständig etabliertes, profitables Geschäft

**Geschäftsziele:**
1. Aufbau einer skalierbaren SEO-Beratungsplattform
2. Integration modernster KI-Technologien für Wettbewerbsvorsprung
3. Automatisierung repetitiver SEO-Analysen
4. Aufbau eines nachhaltigen, wachsenden Kundenstamms

---

## Systemarchitektur

### 1. High-Level Übersicht

```
┌─────────────────────────────────────────────────────────────┐
│                    PHOENIX SEO Platform                      │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Frontend   │  │   Backend    │  │  AI Engine   │      │
│  │   (React)    │◄─┤   (Go)       │◄─┤  (APIs)      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         │                  │                  │              │
│         ▼                  ▼                  ▼              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Dashboard  │  │  PostgreSQL  │  │  Redis Cache │      │
│  │   Reports    │  │  Database    │  │  Job Queue   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 2. Backend-Architektur (Go)

#### 2.1 Projektstruktur

```
phoenix-seo-platform/
├── cmd/
│   ├── api/                    # API Server
│   ├── worker/                 # Background Jobs
│   └── cli/                    # Admin CLI Tools
├── internal/
│   ├── api/
│   │   ├── handlers/          # HTTP Handlers
│   │   ├── middleware/        # Auth, Logging, CORS
│   │   └── routes/            # Route Definitions
│   ├── seo/
│   │   ├── analyzer/          # SEO Analysis Engine
│   │   ├── crawler/           # Web Crawler
│   │   ├── keywords/          # Keyword Research
│   │   ├── backlinks/         # Backlink Analysis
│   │   └── technical/         # Technical SEO Audit
│   ├── ai/
│   │   ├── claude/            # Anthropic Claude Integration
│   │   ├── openai/            # OpenAI Integration
│   │   └── prompts/           # AI Prompt Templates
│   ├── clients/
│   │   ├── management/        # Client Management
│   │   ├── subscriptions/     # Billing & Subscriptions
│   │   └── reports/           # Report Generation
│   ├── database/
│   │   ├── models/            # Data Models
│   │   ├── migrations/        # DB Migrations
│   │   └── repositories/      # Data Access Layer
│   └── queue/
│       ├── jobs/              # Job Definitions
│       └── workers/           # Job Processors
├── pkg/
│   ├── config/                # Configuration
│   ├── logger/                # Logging
│   └── utils/                 # Utilities
└── web/
    ├── src/
    │   ├── components/        # React Components
    │   ├── pages/             # Page Components
    │   ├── services/          # API Services
    │   └── store/             # State Management
    └── public/                # Static Assets
```

#### 2.2 Core-Module

##### SEO-Analyse Engine
- **Crawler:** Intelligente Website-Analyse mit Respektierung von robots.txt
- **On-Page SEO:** Meta-Tags, Headings, Content-Qualität
- **Technical SEO:** Performance, Mobile-Friendly, Core Web Vitals
- **Content-Analyse:** Keyword-Dichte, Lesbarkeit, Struktur
- **Competitor-Analyse:** Vergleich mit Konkurrenten

##### KI-Integration Layer
- **Claude API:** Fortgeschrittene Content-Analyse und Optimierungsvorschläge
- **GPT-4:** Keyword-Generierung, Content-Ideen
- **Automatische Reports:** KI-generierte Zusammenfassungen
- **Trend-Analyse:** Predictive SEO-Trends

##### Client-Management System
- **Kundenprofile:** Vollständige Kundenverwaltung
- **Projekt-Tracking:** Multi-Domain Support pro Kunde
- **Subscription-Management:** Verschiedene Tarife (Basic, Pro, Enterprise)
- **Team-Zugang:** Multi-User Support mit Rollen

---

## 3. Datenbank-Schema

### 3.1 Kern-Entities

```sql
-- Clients/Kunden
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    company VARCHAR(255),
    subscription_tier VARCHAR(50) DEFAULT 'basic',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Projekte/Domains
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    last_crawl_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- SEO Audits
CREATE TABLE seo_audits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    audit_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    score DECIMAL(5,2),
    data JSONB,
    ai_insights JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Keywords
CREATE TABLE keywords (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL,
    search_volume INTEGER,
    difficulty DECIMAL(5,2),
    current_ranking INTEGER,
    target_ranking INTEGER,
    tracked_since TIMESTAMP DEFAULT NOW(),
    last_checked_at TIMESTAMP
);

-- Reports
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    report_type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    data JSONB,
    pdf_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    sent_at TIMESTAMP
);

-- AI Tasks/Jobs
CREATE TABLE ai_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    task_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    input_data JSONB,
    output_data JSONB,
    tokens_used INTEGER,
    cost DECIMAL(10,4),
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);
```

---

## 4. API-Endpoints

### 4.1 Authentication
```
POST   /api/v1/auth/login
POST   /api/v1/auth/register
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
```

### 4.2 Clients
```
GET    /api/v1/clients
POST   /api/v1/clients
GET    /api/v1/clients/:id
PUT    /api/v1/clients/:id
DELETE /api/v1/clients/:id
```

### 4.3 Projects
```
GET    /api/v1/projects
POST   /api/v1/projects
GET    /api/v1/projects/:id
PUT    /api/v1/projects/:id
DELETE /api/v1/projects/:id
POST   /api/v1/projects/:id/crawl
```

### 4.4 SEO Analysis
```
POST   /api/v1/seo/audit
GET    /api/v1/seo/audits/:id
POST   /api/v1/seo/analyze/technical
POST   /api/v1/seo/analyze/content
POST   /api/v1/seo/analyze/keywords
POST   /api/v1/seo/analyze/backlinks
POST   /api/v1/seo/analyze/competitors
```

### 4.5 AI Services
```
POST   /api/v1/ai/analyze-content
POST   /api/v1/ai/suggest-keywords
POST   /api/v1/ai/generate-meta
POST   /api/v1/ai/optimize-content
POST   /api/v1/ai/trend-analysis
```

### 4.6 Reports
```
GET    /api/v1/reports
POST   /api/v1/reports/generate
GET    /api/v1/reports/:id
GET    /api/v1/reports/:id/pdf
POST   /api/v1/reports/:id/send
```

---

## 5. Frontend-Architektur

### 5.1 Technologie-Stack
- **Framework:** React 18+ mit TypeScript
- **State Management:** Zustand oder Redux Toolkit
- **UI Components:** Tailwind CSS + Shadcn/ui
- **Charts:** Recharts oder Chart.js
- **API Client:** Axios mit React Query
- **Routing:** React Router v6

### 5.2 Haupt-Bereiche

#### Dashboard
- Übersicht aller Projekte
- Quick-Stats (Rankings, Audits, Alerts)
- Recent Activity Feed
- Performance-Metriken

#### SEO Analyzer
- Domain-Input & Crawl-Trigger
- Live-Analyse-Status
- Ergebnis-Visualisierung
- AI-Optimierungsvorschläge

#### Client Management
- Client-Liste & Details
- Projekt-Zuordnung
- Subscription-Management
- Kommunikations-Log

#### Reports
- Report-Generierung
- Customizable Templates
- PDF-Export
- Automated Scheduling

#### Settings
- User-Profile
- API-Keys-Management
- Notification-Präferenzen
- Billing & Invoices

---

## 6. KI-Integration Strategie

### 6.1 Claude API (Anthropic)
**Einsatzbereiche:**
- Deep Content-Analyse
- Strategische SEO-Empfehlungen
- Lange, strukturierte Reports
- Competitor-Intelligence

**Vorteile:**
- Sehr gute Textqualität
- Strukturiertes Denken
- Längere Kontexte (200K Tokens)

### 6.2 OpenAI GPT-4
**Einsatzbereiche:**
- Keyword-Research
- Meta-Descriptions generieren
- Content-Briefings
- Quick-Optimierungen

### 6.3 Hybrid-Ansatz
```
SEO-Audit
    ↓
Technical Analysis (Eigenentwicklung)
    ↓
Claude: Strategic Insights
    ↓
GPT-4: Quick Wins & Keywords
    ↓
Kombinierter Report
```

---

## 7. Deployment & Infrastruktur

### 7.1 Produktions-Setup
```
┌─────────────────────────────────────────┐
│         Load Balancer (nginx)           │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────────┐    ┌─────────────┐   │
│  │  API Server │    │  API Server │   │
│  │  (Go)       │    │  (Go)       │   │
│  └─────────────┘    └─────────────┘   │
│                                         │
│  ┌─────────────┐    ┌─────────────┐   │
│  │  Worker     │    │  Worker     │   │
│  │  (Go)       │    │  (Go)       │   │
│  └─────────────┘    └─────────────┘   │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │     PostgreSQL (Primary)         │  │
│  │     + Read Replica               │  │
│  └──────────────────────────────────┘  │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │     Redis (Cache + Queue)        │  │
│  └──────────────────────────────────┘  │
│                                         │
└─────────────────────────────────────────┘
```

### 7.2 Hosting-Optionen
1. **VPS/Dedicated:** Hetzner, DigitalOcean, Linode
2. **Containerized:** Docker + Docker Compose
3. **Cloud:** AWS, GCP, Azure (später bei Skalierung)
4. **CDN:** Cloudflare für Frontend

---

## 8. Geschäftsmodell & Skalierung

### 8.1 Preismodelle

| Tier       | Preis/Monat | Features                           |
|------------|-------------|------------------------------------|
| **Starter**| 99€         | 1 Domain, 10 Keywords, Basis-Reports |
| **Pro**    | 299€        | 5 Domains, 100 Keywords, KI-Analyse |
| **Business**| 699€       | 20 Domains, Unlimited Keywords, API |
| **Enterprise**| Custom   | Unlimited, White-Label, Dedicated  |

### 8.2 Skalierungs-Roadmap

**Phase 1 (Q1 2026):** MVP
- Core SEO Analyzer
- Basic Client Management
- Manual Reports
- 5-10 Beta-Kunden

**Phase 2 (Q2-Q3 2026):** Feature Expansion
- KI-Integration (Claude + GPT-4)
- Automated Reports
- Email-Automation
- 25-50 zahlende Kunden

**Phase 3 (Q4 2026 - Q2 2027):** Growth
- Advanced Analytics
- Competitor-Tracking
- API für Integrationen
- 100-200 Kunden

**Phase 4 (Q3 2027 - Q2 2029):** Scale
- White-Label-Lösung
- Multi-Language Support
- Marketplace für SEO-Tools
- 500+ Kunden, profitabel

---

## 9. Compliance & Rechtliches

### 9.1 DSGVO-Konformität
- Datenschutzerklärung
- Cookie-Consent
- Datenexport-Funktionalität
- Löschfunktionen
- EU-Server-Standort

### 9.2 Geschäftliche Compliance
- Klare AGB
- Transparent pricing
- SLA-Definitionen
- Support-Garantien

### 9.3 PI-Wohlverhalten (bis Q2 2029)
- Transparente Einkommensdokumentation
- Keine Asset-Akkumulation ohne Meldung
- Kommunikation mit Insolvenzverwalter
- Schrittweise Einkommenssteigerung dokumentieren

---

## 10. Monitoring & Success Metrics

### 10.1 Business-KPIs
- **MRR** (Monthly Recurring Revenue)
- **Churn Rate**
- **Customer Acquisition Cost (CAC)**
- **Lifetime Value (LTV)**
- **Net Promoter Score (NPS)**

### 10.2 Technical-KPIs
- API Response Time (<200ms)
- Uptime (99.9%+)
- Crawl Success Rate (>95%)
- AI Task Success Rate (>98%)

### 10.3 SEO-Platform-KPIs
- Audits durchgeführt
- Keywords getrackt
- Reports generiert
- Client-Zufriedenheit

---

## 11. Risiken & Mitigation

| Risiko | Wahrscheinlichkeit | Impact | Mitigation |
|--------|-------------------|--------|------------|
| KI-API-Kosten explodieren | Mittel | Hoch | Rate-Limiting, Caching, Budgets |
| Datenschutz-Verletzung | Niedrig | Sehr Hoch | Security Audits, Encryption, Backups |
| Konkurrenz-Druck | Hoch | Mittel | USP: KI-Vorsprung, Persönlicher Service |
| Technische Schulden | Mittel | Mittel | Code Reviews, Refactoring-Sprints |
| Gesundheitliche Einschränkungen | Mittel | Hoch | Automation, Delegation, Puffer |

---

## 12. Nächste Schritte (Sofort)

1. **Backend-Grundgerüst:** Go-Projekt initialisieren
2. **Datenbank:** PostgreSQL-Schema implementieren
3. **Core-Crawler:** Basis-SEO-Analyzer entwickeln
4. **API:** Erste Endpoints (Auth, Projects, Audits)
5. **Frontend:** React-Dashboard mit Mock-Daten
6. **MVP-Deploy:** Erste Version auf VPS deployen
7. **Beta-Testing:** 3-5 Freunde/Bekannte als Tester

---

**Erstellt:** Januar 2026
**Status:** Living Document - wird kontinuierlich aktualisiert
**Owner:** PHOENIX Team
