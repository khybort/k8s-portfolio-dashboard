# Development Guide - Portfolio Dashboard

## Genel Bakış

Bu rehber, Portfolio Dashboard projesinin geliştirme ortamını kurmak ve geliştirme yapmak için gerekli tüm bilgileri içerir.

## Proje Yapısı

```
k8s-portfolio-dashboard/
├── frontend/              # React frontend
│   ├── public/
│   ├── src/
│   │   ├── admin/         # Admin panel
│   │   ├── public/        # Public site
│   │   ├── shared/        # Shared components
│   │   ├── services/      # API clients
│   │   ├── store/          # State management
│   │   └── utils/
│   ├── package.json
│   └── vite.config.ts
├── backend/               # Go backend API
│   ├── cmd/
│   │   ├── server/        # Main application
│   │   └── migrate/       # Migration tool
│   ├── internal/
│   │   ├── api/           # HTTP handlers
│   │   ├── service/       # Business logic
│   │   ├── repository/    # Data access
│   │   ├── model/          # Domain models
│   │   ├── config/         # Configuration
│   │   └── kafka/          # Kafka integration
│   ├── migrations/         # Database migrations
│   ├── go.mod
│   └── go.sum
├── auth-service/          # Go auth service
│   ├── cmd/
│   │   ├── server/
│   │   └── migrate/
│   ├── internal/
│   │   ├── api/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   └── jwt/
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
├── k8s/                   # Kubernetes manifests
│   ├── frontend/
│   ├── backend/
│   ├── auth-service/
│   └── infrastructure/
├── docker-compose.yml      # Local development
├── docs/                   # Documentation
└── README.md
```

## Backend Development (Go)

### Gereksinimler

- Go 1.23+
- Air (hot reload) - `go install github.com/cosmtrek/air@latest`
- PostgreSQL client tools
- Make (opsiyonel)

### Proje Kurulumu

```bash
cd backend

# Dependencies yükle
go mod download

# Environment variables
cp .env.example .env
# .env dosyasını düzenle

# Database migration'ları çalıştır
go run cmd/migrate/main.go

# Database seed'leri çalıştır (Portfolio, Articles, Projects)
go run cmd/seed/main.go
```

### Development Server

**Normal çalıştırma**:
```bash
go run cmd/server/main.go
```

**Hot reload ile** (önerilen):
```bash
air
```

Air yapılandırması (`.air.toml`):
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### Kod Yapısı

#### API Handlers

```go
// internal/api/handlers/article.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ArticleHandler struct {
    articleService service.ArticleService
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
    // Handler logic
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    // Handler logic
}
```

#### Service Layer

```go
// internal/service/article.go
package service

type ArticleService interface {
    GetArticles(ctx context.Context, page, limit int) ([]model.Article, error)
    GetArticleByID(ctx context.Context, id string) (*model.Article, error)
    CreateArticle(ctx context.Context, article *model.Article) error
    UpdateArticle(ctx context.Context, id string, article *model.Article) error
    DeleteArticle(ctx context.Context, id string) error
}

type articleService struct {
    repo repository.ArticleRepository
    kafka kafka.Producer
}

func (s *articleService) CreateArticle(ctx context.Context, article *model.Article) error {
    // Business logic
    if err := s.repo.Create(ctx, article); err != nil {
        return err
    }
    
    // Publish Kafka event
    event := kafka.Event{
        Type: "article.created",
        Data: article,
    }
    s.kafka.Publish(ctx, "portfolio.articles", event)
    
    return nil
}
```

#### Repository Layer (GORM Best Practices)

Repository'ler GORM best practices'e göre refactor edilmiştir:

```go
// internal/repository/article.go
package repository

import (
    "context"
    "errors"
    "github.com/portfolio/backend/internal/model"
    "gorm.io/gorm"
)

var (
    ErrArticleNotFound = errors.New("article not found")
)

type ArticleRepository interface {
    Create(ctx context.Context, article *model.Article) error
    GetByID(ctx context.Context, id string) (*model.Article, error)
    List(ctx context.Context, page, limit int) ([]model.Article, int64, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
    Delete(ctx context.Context, id string) error
}

type articleRepository struct {
    db *gorm.DB
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
    if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
        return err
    }
    return nil
}

func (r *articleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
    var article model.Article
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&article).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrArticleNotFound
        }
        return nil, err
    }
    return &article, nil
}

func (r *articleRepository) List(ctx context.Context, page, limit int) ([]model.Article, int64, error) {
    var articles []model.Article
    var total int64
    
    offset := (page - 1) * limit
    
    // Count total
    if err := r.db.WithContext(ctx).Model(&model.Article{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated results
    if err := r.db.WithContext(ctx).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&articles).Error; err != nil {
        return nil, 0, err
    }
    
    return articles, total, nil
}

func (r *articleRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
    result := r.db.WithContext(ctx).
        Model(&model.Article{}).
        Where("id = ?", id).
        Updates(updates)
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return ErrArticleNotFound
    }
    
    return nil
}

func (r *articleRepository) Delete(ctx context.Context, id string) error {
    result := r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&model.Article{})
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return ErrArticleNotFound
    }
    
    return nil
}
```

**GORM Best Practices Uygulanan Özellikler**:
- ✅ Context kullanımı (`WithContext`)
- ✅ Custom error handling (`ErrArticleNotFound`)
- ✅ Pagination desteği
- ✅ Partial updates (`Updates` ile map)
- ✅ `RowsAffected` kontrolü
- ✅ Transaction desteği (gerekli yerlerde)
- ✅ Model hooks (`BeforeCreate` için UUID)
- ✅ Index'ler (model tanımlarında)

### Testing

```bash
# Tüm testleri çalıştır
go test ./...

# Coverage ile
go test -cover ./...

# Belirli bir paket
go test ./internal/service/...

# Verbose
go test -v ./...
```

### Database Migrations

```bash
# Migration'ları çalıştır
go run cmd/migrate/main.go

# Migration dosyaları: backend/migrations/
# - 001_create_articles.up.sql
# - 002_create_projects.up.sql
# - 003_create_portfolio.up.sql
```

**Not**: GORM AutoMigrate server başlangıcında otomatik çalışır. SQL migration'lar opsiyoneldir.

### Database Seeding

Seeder'lar resume bilgilerine göre initial data oluşturur:

```bash
# Backend seeder (Portfolio, Articles, Projects)
go run cmd/seed/main.go

# Auth service seeder (Admin user)
cd ../auth-service
go run cmd/seed/main.go
```

**Seed Data**:
- **Portfolio**: Muhsin Kılıç'ın bilgileri (email, bio, social links)
- **Articles**: 5 teknik makale (Go, Kubernetes, React, Kafka, Redis)
- **Projects**: 7 GitHub projesi (RTB Platform, Campaign Dashboard, vb.)
- **Admin User**: admin@portfolio.com / Admin123!

## Frontend Development (React)

### Gereksinimler

- Node.js 18+
- npm 9+ veya yarn/pnpm

### Proje Kurulumu

```bash
cd frontend

# Dependencies yükle
npm install

# Environment variables
cp .env.example .env.local
# .env.local dosyasını düzenle

# Development server başlat
npm run dev
```

### Kod Yapısı

#### Component Yapısı

```typescript
// src/admin/components/ArticleList.tsx
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { articleService } from '@/services/api';

export const ArticleList: React.FC = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['articles'],
    queryFn: () => articleService.getArticles(),
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data?.map(article => (
        <ArticleCard key={article.id} article={article} />
      ))}
    </div>
  );
};
```

#### API Service

```typescript
// src/services/api/articleService.ts
import { apiClient } from './client';
import { Article } from '@/types';

export const articleService = {
  getArticles: async (): Promise<Article[]> => {
    const response = await apiClient.get('/api/v1/articles');
    return response.data;
  },

  getArticle: async (id: string): Promise<Article> => {
    const response = await apiClient.get(`/api/v1/articles/${id}`);
    return response.data;
  },

  createArticle: async (article: CreateArticleDto): Promise<Article> => {
    const response = await apiClient.post('/api/v1/admin/articles', article);
    return response.data;
  },
};
```

#### State Management

```typescript
// src/store/authSlice.ts
import { createSlice } from '@reduxjs/toolkit';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null,
    token: null,
    isAuthenticated: false,
  } as AuthState,
  reducers: {
    setCredentials: (state, action) => {
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.isAuthenticated = true;
    },
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
    },
  },
});
```

### Development Scripts

```bash
# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview

# Linting
npm run lint

# Type checking
npm run type-check

# Testing
npm run test
```

### Routing

```typescript
// src/admin/App.tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ProtectedRoute } from '@/components/ProtectedRoute';

export const AdminApp = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/articles"
          element={
            <ProtectedRoute>
              <ArticlesPage />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
};
```

## Auth Service Development

### Development Server

```bash
cd auth-service

# Dependencies
go mod download

# Environment
cp .env.example .env

# Migrations
go run cmd/migrate/main.go up

# Run server
go run cmd/server/main.go
# veya
air
```

### JWT Implementation

```go
// internal/jwt/token.go
package jwt

func GenerateAccessToken(userID string, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(15 * time.Minute).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

## Kafka Development

### Producer Example

```go
// internal/kafka/producer.go
package kafka

func (p *Producer) PublishArticleCreated(ctx context.Context, article *model.Article) error {
    event := Event{
        Type:      "article.created",
        Timestamp: time.Now(),
        Data:      article,
    }
    
    data, _ := json.Marshal(event)
    
    msg := &kafka.Message{
        Topic: "portfolio.articles",
        Value: data,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte("article.created")},
        },
    }
    
    return p.producer.WriteMessages(ctx, msg)
}
```

### Consumer Example

```go
// internal/kafka/consumer.go
package kafka

func (c *Consumer) ConsumeArticles(ctx context.Context) error {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  []string{"localhost:9092"},
        Topic:     "portfolio.articles",
        GroupID:   "email-service",
        MinBytes:  10e3,
        MaxBytes:  10e6,
    })
    defer reader.Close()
    
    for {
        msg, err := reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        
        var event Event
        json.Unmarshal(msg.Value, &event)
        
        switch event.Type {
        case "article.created":
            c.handleArticleCreated(ctx, event)
        }
    }
}
```

## Redis Development

### Cache Example

```go
// internal/cache/redis.go
package cache

func (c *RedisCache) GetArticle(ctx context.Context, id string) (*model.Article, error) {
    key := fmt.Sprintf("article:%s", id)
    
    data, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, ErrCacheMiss
    }
    if err != nil {
        return nil, err
    }
    
    var article model.Article
    json.Unmarshal([]byte(data), &article)
    return &article, nil
}

func (c *RedisCache) SetArticle(ctx context.Context, article *model.Article, ttl time.Duration) error {
    key := fmt.Sprintf("article:%s", article.ID)
    data, _ := json.Marshal(article)
    
    return c.client.Set(ctx, key, data, ttl).Err()
}
```

## Debugging

### Backend Debugging

**VS Code Launch Configuration**:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/server",
      "env": {
        "ENV": "development"
      }
    }
  ]
}
```

### Frontend Debugging

**React DevTools**: Browser extension
**Redux DevTools**: Browser extension
**Network Tab**: API request'leri izleme

## Best Practices

### Go

1. **Error Handling**: Her zaman error'ları kontrol et
2. **Context**: Timeout ve cancellation için context kullan
3. **Interfaces**: Dependency injection için interface'ler kullan
4. **Testing**: Unit test ve integration test yaz
5. **Logging**: Structured logging (zap) kullan

### React

1. **TypeScript**: Type safety için TypeScript kullan
2. **Component Composition**: Küçük, reusable component'ler
3. **Custom Hooks**: Logic'i hook'lara ayır
4. **Error Boundaries**: Error handling için Error Boundary
5. **Code Splitting**: Lazy loading için React.lazy

## Git Workflow

```bash
# Feature branch oluştur
git checkout -b feature/article-management

# Commit
git add .
git commit -m "feat: add article management"

# Push
git push origin feature/article-management

# Pull request oluştur
```

## Sonraki Adımlar

- [API Documentation](./API.md) - API endpoint'leri
- [Deployment Guide](./DEPLOYMENT.md) - Production deployment
- [Kafka Integration](./KAFKA.md) - Kafka detayları
- [Redis Integration](./REDIS.md) - Redis caching
