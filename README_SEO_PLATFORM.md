# 🔥 PHOENIX SEO Platform

**Professionelle SEO-Beratungsplattform mit KI-Integration**

Eine moderne, skalierbare SEO-Analyse- und Beratungsplattform, die fortschrittliche KI-Technologien (Claude & GPT-4) nutzt, um Unternehmen zu helfen, ihre Online-Sichtbarkeit zu maximieren.

---

## 🎯 Mission & Vision

**Ziel:** Positionierung als führende SEO-Beratungsinstanz mit technologischem Vorsprung durch KI-Integration.

**Timeline:**
- **Start:** Februar 2026
- **Ziel:** Q2 2029 - Etabliertes, profitables SaaS-Geschäft

---

## ✨ Features

### 🔍 SEO-Analyse Engine
- **Website Crawler:** Intelligente Analyse von Websites mit robots.txt-Respektierung
- **On-Page SEO:** Meta-Tags, Überschriften, Content-Qualität
- **Technical SEO:** Performance, Mobile-Friendly, HTTPS
- **Content-Analyse:** Keyword-Dichte, Lesbarkeit, Struktur
- **Scoring-System:** Detaillierte Bewertung mit Handlungsempfehlungen

### 🤖 KI-Integration
- **Claude (Anthropic):** Deep Content-Analyse, strategische Empfehlungen
- **GPT-4 (OpenAI):** Keyword-Research, Meta-Tag-Generierung, Content-Ideen
- **Automatische Insights:** KI-generierte Optimierungsvorschläge
- **Smart Recommendations:** Priorisierte Handlungsempfehlungen

### 📊 Client Management
- Kundenverwaltung mit Projekten
- Multi-Domain Support
- Subscription-Tiers (Starter, Pro, Business, Enterprise)
- Report-Generierung

### 🚀 Performance
- Asynchrone Crawling-Engine
- Redis-Caching für schnelle Antworten
- PostgreSQL für robuste Datenhaltung
- Skalierbare Microservice-Architektur

---

## 🏗️ Architektur

```
┌─────────────────────────────────────────────────────────┐
│              PHOENIX SEO Platform                        │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  Frontend (React)  ◄──► Backend (Go)  ◄──► AI APIs      │
│       ▼                      ▼                 ▼         │
│  Dashboard            PostgreSQL         Claude/GPT-4    │
│  Reports              Redis Cache                        │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

**Technologie-Stack:**
- **Backend:** Go 1.24+
- **Database:** PostgreSQL 16
- **Cache:** Redis 7
- **AI:** Anthropic Claude, OpenAI GPT-4
- **Deployment:** Docker, Docker Compose

**Siehe:** [Detaillierte Architektur-Dokumentation](docs/ARCHITECTURE.md)

---

## 🚀 Quick Start

### Voraussetzungen

- Docker & Docker Compose
- Go 1.24+ (für lokale Entwicklung)
- PostgreSQL 16+ (optional, wenn nicht Docker)
- API-Keys: Claude, OpenAI

### 1. Repository klonen

```bash
git clone https://github.com/EricFreesoul/phoenix-feuer-os.git
cd phoenix-feuer-os
```

### 2. Umgebungsvariablen konfigurieren

```bash
cp .env.example .env
```

Bearbeite `.env` und füge deine API-Keys ein:

```env
CLAUDE_API_KEY=your_claude_api_key_here
OPENAI_API_KEY=your_openai_api_key_here
DB_PASSWORD=your_secure_password
```

### 3. Mit Docker starten

```bash
# Alle Services starten
docker-compose up -d

# Logs verfolgen
docker-compose logs -f api
```

### 4. Datenbank initialisieren

```bash
make db-migrate
```

### 5. API testen

```bash
curl http://localhost:8080/api/v1/health
```

---

## 📖 API-Dokumentation

### Health Check

```bash
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1704369600,
  "service": "phoenix-seo-api"
}
```

### URL Analysieren

```bash
POST /api/v1/seo/analyze
Content-Type: application/json

{
  "url": "https://example.com",
  "keywords": ["seo", "marketing"],
  "use_ai": true
}
```

**Response:**
```json
{
  "crawl_result": {
    "url": "https://example.com",
    "status_code": 200,
    "title": "Example Domain",
    "meta_description": "Example description",
    "h1_tags": ["Example Domain"],
    "load_time_ms": 345,
    "word_count": 150
  },
  "seo_score": {
    "overall": 75.5,
    "technical": 85,
    "content": 70,
    "on_page": 80,
    "performance": 90,
    "issues": [...],
    "opportunities": [...]
  },
  "ai_insights": "Claude's analysis here...",
  "analyzed_at": "2026-01-04T10:00:00Z"
}
```

### Keywords generieren

```bash
POST /api/v1/seo/keywords/generate
Content-Type: application/json

{
  "topic": "E-Commerce SEO",
  "count": 10
}
```

### Meta-Tags optimieren

```bash
POST /api/v1/seo/meta/optimize
Content-Type: application/json

{
  "content": "Your page content here...",
  "keywords": ["seo", "optimization"]
}
```

---

## 🛠️ Entwicklung

### Lokale Entwicklung (ohne Docker)

```bash
# Dependencies installieren
make install

# API starten
make run

# Mit Auto-Reload (requires air)
make dev
```

### Tests ausführen

```bash
# Alle Tests
make test

# Mit Coverage-Report
make test-coverage
```

### Code-Qualität

```bash
# Formatieren
make fmt

# Linting
make lint

# Security-Check
make security
```

### Datenbank-Management

```bash
# PostgreSQL Shell öffnen
make db-shell

# Redis CLI öffnen
make redis-cli

# Migrations ausführen
make db-migrate
```

---

## 📁 Projektstruktur

```
phoenix-feuer-os/
├── cmd/
│   ├── api/              # API Server
│   ├── worker/           # Background Jobs
│   └── cli/              # Admin CLI
├── internal/
│   ├── api/              # API Layer
│   │   ├── handlers/     # HTTP Handlers
│   │   ├── middleware/   # CORS, Auth, Logging
│   │   └── routes/       # Route Definitions
│   ├── seo/              # SEO Engine
│   │   ├── analyzer/     # SEO Analysis
│   │   ├── crawler/      # Web Crawler
│   │   ├── keywords/     # Keyword Research
│   │   └── technical/    # Technical SEO
│   ├── ai/               # AI Integration
│   │   ├── claude/       # Anthropic Claude
│   │   └── openai/       # OpenAI GPT
│   ├── clients/          # Client Management
│   ├── database/         # Database Layer
│   │   ├── models/       # Data Models
│   │   ├── migrations/   # SQL Migrations
│   │   └── repositories/ # Data Access
│   └── queue/            # Job Queue
├── pkg/
│   ├── config/           # Configuration
│   ├── logger/           # Logging
│   └── utils/            # Utilities
├── docs/                 # Documentation
├── web/                  # Frontend (React)
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## 💰 Geschäftsmodell

### Pricing Tiers

| Tier | Preis/Monat | Features |
|------|-------------|----------|
| **Starter** | 99€ | 1 Domain, 10 Keywords, Basis-Reports |
| **Pro** | 299€ | 5 Domains, 100 Keywords, KI-Analyse |
| **Business** | 699€ | 20 Domains, Unlimited Keywords, API-Access |
| **Enterprise** | Custom | Unlimited, White-Label, Support |

### Skalierungs-Roadmap

**Q1 2026:** MVP mit Core-Features
- 5-10 Beta-Kunden

**Q2-Q3 2026:** Feature-Expansion
- KI-Integration vollständig
- 25-50 zahlende Kunden

**Q4 2026 - Q2 2027:** Growth Phase
- API-Integrationen
- 100-200 Kunden

**Q3 2027 - Q2 2029:** Scale
- White-Label-Lösung
- 500+ Kunden, profitabel

---

## 🔒 Sicherheit & Compliance

### DSGVO-Konformität
- EU-Server-Standort
- Datenexport-Funktionalität
- Löschfunktionen
- Transparente Datenschutzerklärung

### Sicherheits-Features
- HTTPS-Only
- API-Key-Authentifizierung
- Rate-Limiting
- Input-Validierung
- SQL-Injection-Schutz

---

## 📊 Monitoring & Metrics

### Business-KPIs
- MRR (Monthly Recurring Revenue)
- Churn Rate
- Customer Acquisition Cost (CAC)
- Lifetime Value (LTV)

### Technical-KPIs
- API Response Time (<200ms Ziel)
- Uptime (99.9%+ Ziel)
- Crawl Success Rate (>95%)
- AI Task Success Rate (>98%)

---

## 🚢 Deployment

### Docker Deployment

```bash
# Produktions-Build
docker build -t phoenix-seo-api:latest .

# Container starten
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  phoenix-seo-api:latest
```

### Cloud Deployment (Hetzner, DigitalOcean)

1. VPS mit Ubuntu 22.04+ erstellen
2. Docker & Docker Compose installieren
3. Repository klonen
4. `.env` konfigurieren
5. `docker-compose up -d` ausführen
6. Nginx als Reverse Proxy konfigurieren
7. SSL-Zertifikat (Let's Encrypt) einrichten

---

## 🤝 Support & Kontakt

**Projekt:** PHOENIX SEO Platform
**Version:** 1.0.0 (MVP)
**Lizenz:** Proprietär

**Für Support:**
- GitHub Issues: [Issues erstellen](https://github.com/EricFreesoul/phoenix-feuer-os/issues)
- Email: support@phoenix-seo.com (coming soon)

---

## 📝 Roadmap

### ✅ Phase 1 - MVP (Q1 2026)
- [x] Backend-Architektur
- [x] SEO-Crawler & Analyzer
- [x] KI-Integration (Claude + GPT-4)
- [x] API-Endpoints
- [ ] Frontend-Dashboard
- [ ] Client-Management
- [ ] Report-Generierung

### 🔄 Phase 2 - Feature Expansion (Q2-Q3 2026)
- [ ] Automated Reports
- [ ] Email-Automation
- [ ] Keyword-Tracking über Zeit
- [ ] Competitor-Analysis
- [ ] Backlink-Monitoring
- [ ] Advanced Analytics

### 🚀 Phase 3 - Growth (Q4 2026+)
- [ ] API für Integrationen
- [ ] White-Label-Option
- [ ] Multi-Language Support
- [ ] Mobile App
- [ ] Marketplace für SEO-Tools

---

## 📄 Lizenz

© 2026 PHOENIX SEO Platform. Alle Rechte vorbehalten.

Dieses Projekt ist proprietäre Software und nicht für öffentliche Nutzung oder Weiterverbreitung lizenziert.

---

**Gebaut mit ❤️ und 🔥 für die Zukunft der SEO-Beratung**
