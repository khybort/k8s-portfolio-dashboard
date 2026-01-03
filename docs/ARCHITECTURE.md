# Portfolio Dashboard - Architecture Documentation

## Genel Bakış

Bu proje, **mikroservis mimarisi** kullanan, **Go backend** ve **React frontend** ile geliştirilmiş, **Kubernetes** üzerinde çalışan profesyonel bir portfolio yönetim sistemidir. Sistem, admin paneli, içerik yönetimi, makale yönetimi ve GitHub proje entegrasyonu içerir.

## Mimari Genel Bakış

### Mikroservis Mimarisi

Sistem aşağıdaki servislerden oluşur:

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Frontend   │  │    Backend   │  │  Auth Service│     │
│  │   (React)    │──│     (Go)     │──│     (Go)     │     │
│  │              │  │              │  │              │     │
│  └──────────────┘  └──────┬───────┘  └──────┬───────┘     │
│                           │                  │              │
│                    ┌──────▼───────┐  ┌──────▼───────┐     │
│                    │   Kafka      │  │    Redis     │     │
│                    │  (Events)    │  │   (Cache)    │     │
│                    └──────────────┘  └──────────────┘     │
│                                                              │
│                    ┌──────────────┐                        │
│                    │  PostgreSQL  │                        │
│                    │  (Database)  │                        │
│                    └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘
```

## Servisler ve Bileşenler

### 1. Frontend Servisi (React)

#### Public Site
- **Teknoloji**: React 18+, TypeScript, Vite
- **Özellikler**:
  - Responsive portfolio showcase
  - Makale listesi ve detay sayfaları
  - GitHub projeleri gösterimi
  - Dark/Light mode
  - SEO optimizasyonu
  - PWA desteği

#### Admin Panel
- **Teknoloji**: React 18+, TypeScript, React Router, React Query
- **Özellikler**:
  - Kullanıcı authentication (login/logout)
  - Dashboard (istatistikler, grafikler)
  - Makale yönetimi (CRUD)
  - GitHub proje yönetimi (CRUD)
  - Portfolio içerik yönetimi
  - Media yönetimi (resim upload)
  - Kullanıcı ayarları

**Dizin Yapısı**:
```
frontend/
├── public/              # Static files
├── src/
│   ├── admin/           # Admin panel
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   └── utils/
│   ├── public/          # Public site
│   │   ├── components/
│   │   ├── pages/
│   │   └── hooks/
│   ├── shared/          # Shared components
│   ├── services/        # API clients
│   ├── store/           # State management
│   └── utils/
├── package.json
└── vite.config.ts
```

### 2. Backend Servisi (Go)

#### Teknoloji Stack
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **ORM**: GORM (database)
- **Validation**: validator
- **Logging**: zap
- **Config**: viper

#### API Endpoints

**Public API**:
- `GET /api/v1/articles` - Makale listesi
- `GET /api/v1/articles/:id` - Makale detayı
- `GET /api/v1/projects` - GitHub projeleri
- `GET /api/v1/projects/:id` - Proje detayı
- `GET /api/v1/portfolio` - Portfolio bilgileri

**Admin API** (Auth gerekli):
- `POST /api/v1/admin/articles` - Makale oluştur
- `PUT /api/v1/admin/articles/:id` - Makale güncelle
- `DELETE /api/v1/admin/articles/:id` - Makale sil
- `POST /api/v1/admin/projects` - Proje ekle
- `PUT /api/v1/admin/projects/:id` - Proje güncelle
- `DELETE /api/v1/admin/projects/:id` - Proje sil
- `PUT /api/v1/admin/portfolio` - Portfolio güncelle

**Dizin Yapısı**:
```
backend/
├── cmd/
│   └── server/
│       └── main.go       # Application entry point
├── internal/
│   ├── api/              # HTTP handlers
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes.go
│   ├── service/          # Business logic
│   │   ├── article.go
│   │   ├── project.go
│   │   └── portfolio.go
│   ├── repository/       # Data access layer
│   │   ├── article.go
│   │   ├── project.go
│   │   └── portfolio.go
│   ├── model/            # Domain models
│   │   ├── article.go
│   │   ├── project.go
│   │   └── portfolio.go
│   ├── config/           # Configuration
│   └── kafka/            # Kafka producer/consumer
├── pkg/                  # Shared packages
│   ├── logger/
│   ├── validator/
│   └── errors/
├── migrations/           # Database migrations
├── go.mod
└── go.sum
```

### 3. Auth Servisi (Go)

#### Teknoloji Stack
- **Language**: Go 1.21+
- **Framework**: Gin
- **JWT**: golang-jwt/jwt
- **Password**: bcrypt
- **Database**: PostgreSQL (kullanıcı bilgileri)

#### Özellikler
- Kullanıcı kaydı (register)
- Kullanıcı girişi (login)
- JWT token üretimi ve doğrulama
- Refresh token desteği
- Password reset
- Role-based access control (RBAC)

#### API Endpoints
- `POST /api/v1/auth/register` - Kullanıcı kaydı
- `POST /api/v1/auth/login` - Kullanıcı girişi
- `POST /api/v1/auth/refresh` - Token yenileme
- `POST /api/v1/auth/logout` - Çıkış
- `POST /api/v1/auth/verify` - Token doğrulama
- `POST /api/v1/auth/password-reset` - Şifre sıfırlama

**Dizin Yapısı**:
```
auth-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   ├── service/
│   ├── repository/
│   ├── model/
│   └── jwt/
├── go.mod
└── go.sum
```

### 4. Kafka (Event Streaming)

#### Kullanım Senaryoları

**1. Event-Driven Architecture**
- Makale oluşturulduğunda → `article.created` eventi
- Makale güncellendiğinde → `article.updated` eventi
- Proje eklendiğinde → `project.created` eventi

**2. Async Processing**
- Email gönderimi (makale yayınlandığında)
- Analytics event'leri
- Cache invalidation
- Search index güncelleme

**3. Topics**
- `portfolio.articles` - Makale event'leri
- `portfolio.projects` - Proje event'leri
- `portfolio.analytics` - Analytics event'leri
- `portfolio.notifications` - Bildirim event'leri

**4. Consumer Groups**
- `email-service` - Email gönderimi
- `analytics-service` - Analytics toplama
- `cache-service` - Cache yönetimi
- `search-service` - Search index güncelleme

### 5. Redis (Caching & Session)

#### Kullanım Senaryoları

**1. Caching**
- Makale listesi cache (TTL: 5 dakika)
- Makale detay cache (TTL: 10 dakika)
- Proje listesi cache (TTL: 5 dakika)
- Portfolio bilgileri cache (TTL: 1 saat)

**2. Session Storage**
- JWT refresh token'ları
- User session'ları
- Rate limiting counters

**3. Cache Strategy**
```go
// Cache key patterns
articles:list:page:{page}:limit:{limit}
article:detail:{id}
projects:list:page:{page}:limit:{limit}
project:detail:{id}
portfolio:info
```

**4. Cache Invalidation**
- Kafka event'leri ile otomatik invalidation
- Manual invalidation API endpoint'leri

### 6. PostgreSQL (Database)

#### Schema

**articles**:
- id (UUID, PK)
- title (VARCHAR)
- slug (VARCHAR, UNIQUE)
- content (TEXT)
- excerpt (TEXT)
- author_id (UUID, FK)
- published_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

**projects**:
- id (UUID, PK)
- name (VARCHAR)
- description (TEXT)
- github_url (VARCHAR)
- live_url (VARCHAR)
- technologies (TEXT[])
- featured (BOOLEAN)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

**portfolio**:
- id (UUID, PK)
- name (VARCHAR)
- title (VARCHAR)
- bio (TEXT)
- email (VARCHAR)
- social_links (JSONB)
- settings (JSONB)
- updated_at (TIMESTAMP)

**users** (Auth Service):
- id (UUID, PK)
- email (VARCHAR, UNIQUE)
- password_hash (VARCHAR)
- role (VARCHAR)
- created_at (TIMESTAMP)

## Veri Akışı

### 1. Kullanıcı Girişi
```
Frontend (Login Form)
    ↓
Auth Service (/api/v1/auth/login)
    ↓
PostgreSQL (User verification)
    ↓
JWT Token Generation
    ↓
Redis (Store refresh token)
    ↓
Frontend (Store JWT, redirect to admin)
```

### 2. Makale Oluşturma
```
Admin Panel (Create Article Form)
    ↓
Backend API (/api/v1/admin/articles)
    ↓
Auth Middleware (JWT verification)
    ↓
PostgreSQL (Save article)
    ↓
Kafka Producer (Publish article.created event)
    ↓
Redis (Invalidate cache)
    ↓
Kafka Consumers:
  - Email Service (Send notification)
  - Analytics Service (Track event)
  - Search Service (Update index)
```

### 3. Makale Listesi (Public)
```
Public Site (Article List Page)
    ↓
Backend API (/api/v1/articles)
    ↓
Redis (Check cache)
    ↓
[Cache Hit] → Return cached data
[Cache Miss] → PostgreSQL Query
    ↓
Redis (Store in cache)
    ↓
Return response
```

## Güvenlik

### Authentication & Authorization
- **JWT Tokens**: Access token (15 dk) + Refresh token (7 gün)
- **HTTPS**: Tüm iletişim şifrelenmiş
- **CORS**: Frontend domain'leri için CORS yapılandırması
- **Rate Limiting**: API endpoint'leri için rate limiting
- **Input Validation**: Tüm input'lar validate edilir

### Database Security
- **Prepared Statements**: SQL injection koruması
- **Connection Pooling**: Güvenli bağlantı yönetimi
- **Encryption**: Hassas veriler şifrelenir

### Container Security
- **Non-root User**: Container'lar non-root kullanıcı ile çalışır
- **Read-only Filesystem**: Mümkün olduğunca read-only
- **Security Scanning**: Image'lar Trivy ile taranır

## Ölçeklenebilirlik

### Horizontal Scaling
- **Frontend**: Stateless, kolayca scale edilebilir
- **Backend**: Stateless API, multiple replicas
- **Auth Service**: Stateless, JWT-based
- **Kafka**: Partition-based scaling
- **Redis**: Cluster mode
- **PostgreSQL**: Read replicas

### Load Balancing
- **Kubernetes Service**: Automatic load balancing
- **Ingress**: Nginx Ingress Controller
- **Database**: Connection pooling

### Caching Strategy
- **Redis**: Hot data caching
- **CDN**: Static assets için CDN
- **Browser Cache**: HTTP cache headers

## Monitoring ve Observability

### Logging
- **Structured Logging**: JSON format, zap logger
- **Log Aggregation**: Loki veya ELK stack
- **Log Levels**: DEBUG, INFO, WARN, ERROR

### Metrics
- **Prometheus**: Application metrics
- **Grafana**: Visualization
- **Custom Metrics**: Request count, latency, error rate

### Tracing
- **OpenTelemetry**: Distributed tracing
- **Jaeger**: Trace visualization

### Health Checks
- **Liveness Probe**: `/healthz`
- **Readiness Probe**: `/ready`
- **Dependencies**: Database, Redis, Kafka connectivity

## Deployment

### Kubernetes Manifests

**Namespace**: `portfolio`

**Services**:
- `frontend` (Deployment + Service)
- `backend` (Deployment + Service)
- `auth-service` (Deployment + Service)
- `kafka` (StatefulSet + Service)
- `redis` (StatefulSet + Service)
- `postgresql` (StatefulSet + Service)

**ConfigMaps**:
- Application configurations
- Environment variables

**Secrets**:
- Database credentials
- JWT secrets
- API keys

**Ingress**:
- Public site: `portfolio.example.com`
- Admin panel: `admin.portfolio.example.com`
- API: `api.portfolio.example.com`

## Geliştirme Ortamı

### Local Development
- **Docker Compose**: Tüm servisleri lokal çalıştırma
- **Hot Reload**: Frontend ve backend hot reload
- **Local Kafka**: Kafka container
- **Local Redis**: Redis container
- **Local PostgreSQL**: PostgreSQL container

### CI/CD
- **GitHub Actions**: Automated build and test
- **Docker Registry**: GHCR veya Docker Hub
- **Kubernetes Deployment**: Automated deployment
- **Rolling Updates**: Zero-downtime deployments

## Production Considerations

- **High Availability**: Multiple replicas, multi-zone
- **Disaster Recovery**: Database backups, disaster recovery plan
- **Security**: WAF, DDoS protection, security scanning
- **Performance**: CDN, caching, database optimization
- **Compliance**: GDPR, data privacy regulations
