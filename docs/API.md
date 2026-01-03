# API Documentation - Portfolio Dashboard

## Genel Bakış

Portfolio Dashboard REST API, Go backend servisi tarafından sağlanır. API, makale yönetimi, proje yönetimi ve portfolio içerik yönetimi için endpoint'ler içerir.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://api.portfolio.example.com`

## Authentication

API, JWT (JSON Web Token) tabanlı authentication kullanır. Auth Service'den alınan access token, request header'ında gönderilmelidir.

### Authentication Header

```http
Authorization: Bearer <access_token>
```

### Token Alımı

Auth Service'den login endpoint'i ile token alınır:

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "password"
}
```

**Response**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

Detaylar için: [Auth Service Documentation](./AUTH.md)

## API Endpoints

### Health Check

#### GET /healthz

**Açıklama**: Servis sağlık kontrolü

**Request**:
```http
GET /healthz HTTP/1.1
```

**Response** (200 OK):
```
ok
```

---

## Public API Endpoints

### Articles

#### GET /api/v1/articles

**Açıklama**: Makale listesi (pagination ile)

**Query Parameters**:
- `page` (int, default: 1): Sayfa numarası
- `limit` (int, default: 10): Sayfa başına kayıt
- `sort` (string, default: "created_at"): Sıralama alanı
- `order` (string, default: "desc"): Sıralama yönü (asc/desc)

**Request**:
```http
GET /api/v1/articles?page=1&limit=10&sort=created_at&order=desc
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Article Title",
      "slug": "article-slug",
      "excerpt": "Article excerpt...",
      "content": "Full article content...",
      "author": {
        "id": "uuid",
        "name": "Author Name"
      },
      "published_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

#### GET /api/v1/articles/:id

**Açıklama**: Makale detayı

**Request**:
```http
GET /api/v1/articles/uuid
```

**Response** (200 OK):
```json
{
  "id": "uuid",
  "title": "Article Title",
  "slug": "article-slug",
  "excerpt": "Article excerpt...",
  "content": "Full article content...",
  "author": {
    "id": "uuid",
    "name": "Author Name",
    "email": "author@example.com"
  },
  "published_at": "2024-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### GET /api/v1/articles/slug/:slug

**Açıklama**: Slug ile makale getir

**Request**:
```http
GET /api/v1/articles/slug/article-slug
```

**Response**: GET /api/v1/articles/:id ile aynı

### Projects

#### GET /api/v1/projects

**Açıklama**: GitHub projeleri listesi

**Query Parameters**:
- `page` (int, default: 1)
- `limit` (int, default: 10)
- `featured` (bool, optional): Sadece featured projeler

**Request**:
```http
GET /api/v1/projects?page=1&limit=10&featured=true
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Project Name",
      "description": "Project description...",
      "github_url": "https://github.com/user/repo",
      "live_url": "https://project.example.com",
      "technologies": ["Go", "React", "Kubernetes"],
      "featured": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

#### GET /api/v1/projects/:id

**Açıklama**: Proje detayı

**Request**:
```http
GET /api/v1/projects/uuid
```

**Response** (200 OK):
```json
{
  "id": "uuid",
  "name": "Project Name",
  "description": "Project description...",
  "github_url": "https://github.com/user/repo",
  "live_url": "https://project.example.com",
  "technologies": ["Go", "React", "Kubernetes"],
  "featured": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Portfolio

#### GET /api/v1/portfolio

**Açıklama**: Portfolio bilgileri

**Request**:
```http
GET /api/v1/portfolio
```

**Response** (200 OK):
```json
{
  "id": "uuid",
  "name": "John Doe",
  "title": "Full Stack Developer",
  "bio": "Passionate developer...",
  "email": "john@example.com",
  "social_links": {
    "github": "https://github.com/johndoe",
    "linkedin": "https://linkedin.com/in/johndoe",
    "twitter": "https://twitter.com/johndoe"
  },
  "settings": {
    "theme": "dark",
    "language": "en"
  },
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## Admin API Endpoints

**Not**: Tüm admin endpoint'leri authentication gerektirir.

### Articles Management

#### POST /api/v1/admin/articles

**Açıklama**: Yeni makale oluştur

**Authentication**: Required (Admin)

**Request**:
```http
POST /api/v1/admin/articles
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Article Title",
  "slug": "article-slug",
  "excerpt": "Article excerpt...",
  "content": "Full article content...",
  "published": true
}
```

**Response** (201 Created):
```json
{
  "id": "uuid",
  "title": "Article Title",
  "slug": "article-slug",
  "excerpt": "Article excerpt...",
  "content": "Full article content...",
  "author_id": "uuid",
  "published": true,
  "published_at": "2024-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /api/v1/admin/articles/:id

**Açıklama**: Makale güncelle

**Authentication**: Required (Admin)

**Request**:
```http
PUT /api/v1/admin/articles/uuid
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Updated Article Title",
  "content": "Updated content...",
  "published": true
}
```

**Response** (200 OK):
```json
{
  "id": "uuid",
  "title": "Updated Article Title",
  "slug": "article-slug",
  "content": "Updated content...",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### DELETE /api/v1/admin/articles/:id

**Açıklama**: Makale sil

**Authentication**: Required (Admin)

**Request**:
```http
DELETE /api/v1/admin/articles/uuid
Authorization: Bearer <access_token>
```

**Response** (200 OK):
```json
{
  "message": "Article deleted successfully"
}
```

### Projects Management

#### POST /api/v1/admin/projects

**Açıklama**: Yeni proje ekle

**Authentication**: Required (Admin)

**Request**:
```http
POST /api/v1/admin/projects
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Project Name",
  "description": "Project description...",
  "github_url": "https://github.com/user/repo",
  "live_url": "https://project.example.com",
  "technologies": ["Go", "React", "Kubernetes"],
  "featured": true
}
```

**Response** (201 Created):
```json
{
  "id": "uuid",
  "name": "Project Name",
  "description": "Project description...",
  "github_url": "https://github.com/user/repo",
  "live_url": "https://project.example.com",
  "technologies": ["Go", "React", "Kubernetes"],
  "featured": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /api/v1/admin/projects/:id

**Açıklama**: Proje güncelle

**Authentication**: Required (Admin)

**Request**:
```http
PUT /api/v1/admin/projects/uuid
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Updated Project Name",
  "featured": false
}
```

**Response** (200 OK): GET /api/v1/projects/:id ile aynı format

#### DELETE /api/v1/admin/projects/:id

**Açıklama**: Proje sil

**Authentication**: Required (Admin)

**Request**:
```http
DELETE /api/v1/admin/projects/uuid
Authorization: Bearer <access_token>
```

**Response** (200 OK):
```json
{
  "message": "Project deleted successfully"
}
```

### Portfolio Management

#### PUT /api/v1/admin/portfolio

**Açıklama**: Portfolio bilgilerini güncelle

**Authentication**: Required (Admin)

**Request**:
```http
PUT /api/v1/admin/portfolio
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "John Doe",
  "title": "Full Stack Developer",
  "bio": "Updated bio...",
  "email": "john@example.com",
  "social_links": {
    "github": "https://github.com/johndoe",
    "linkedin": "https://linkedin.com/in/johndoe"
  }
}
```

**Response** (200 OK): GET /api/v1/portfolio ile aynı format

### Cache Management

#### POST /api/v1/admin/cache/invalidate

**Açıklama**: Cache'i invalidate et

**Authentication**: Required (Admin)

**Query Parameters**:
- `type` (string, optional): Cache tipi (articles, projects, all)

**Request**:
```http
POST /api/v1/admin/cache/invalidate?type=articles
Authorization: Bearer <access_token>
```

**Response** (200 OK):
```json
{
  "message": "Cache invalidated successfully",
  "type": "articles"
}
```

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "Bad Request",
  "message": "Invalid request parameters",
  "code": "INVALID_PARAMS",
  "details": {
    "field": "title",
    "reason": "Title is required"
  }
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing authentication token",
  "code": "UNAUTHORIZED"
}
```

### 403 Forbidden

```json
{
  "error": "Forbidden",
  "message": "Insufficient permissions",
  "code": "FORBIDDEN"
}
```

### 404 Not Found

```json
{
  "error": "Not Found",
  "message": "Resource not found",
  "code": "NOT_FOUND"
}
```

### 429 Too Many Requests

```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "retry_after": 60
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred",
  "code": "INTERNAL_ERROR"
}
```

## Rate Limiting

API endpoint'leri rate limiting ile korunur:

- **Public Endpoints**: 100 requests per minute per IP
- **Admin Endpoints**: 200 requests per minute per user

**Headers**:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1704067200
```

## Pagination

Liste endpoint'leri pagination kullanır:

**Query Parameters**:
- `page`: Sayfa numarası (default: 1)
- `limit`: Sayfa başına kayıt (default: 10, max: 100)

**Response Format**:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Filtering & Sorting

### Articles

**Filtering**:
- `published` (bool): Sadece yayınlanmış makaleler

**Sorting**:
- `sort`: created_at, updated_at, title
- `order`: asc, desc

### Projects

**Filtering**:
- `featured` (bool): Sadece featured projeler

**Sorting**:
- `sort`: created_at, updated_at, name
- `order`: asc, desc

## CORS

**Allowed Origins**:
- `http://localhost:5173` (development)
- `https://portfolio.example.com` (production)
- `https://admin.portfolio.example.com` (production)

**Headers**:
```http
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## API Versioning

- **Current Version**: v1
- **Version Header**: `API-Version: v1`
- **URL Versioning**: `/api/v1/...`

## Testing

### cURL Examples

```bash
# Get articles
curl http://localhost:8080/api/v1/articles

# Get article by ID
curl http://localhost:8080/api/v1/articles/uuid

# Create article (admin)
curl -X POST http://localhost:8080/api/v1/admin/articles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Test content"}'

# Update article (admin)
curl -X PUT http://localhost:8080/api/v1/admin/articles/uuid \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title"}'

# Delete article (admin)
curl -X DELETE http://localhost:8080/api/v1/admin/articles/uuid \
  -H "Authorization: Bearer <token>"
```

## SDK/Client Libraries

- **JavaScript/TypeScript**: `@portfolio/api-client`
- **Go**: `github.com/portfolio/api-client-go`
- **Python**: `portfolio-api-client`

## OpenAPI/Swagger

API dokümantasyonu Swagger UI'da mevcuttur:

- **Development**: `http://localhost:8080/swagger`
- **Production**: `https://api.portfolio.example.com/swagger`

## Changelog

### v1.0.0
- ✅ Public article endpoints
- ✅ Public project endpoints
- ✅ Portfolio info endpoint
- ✅ Admin article management
- ✅ Admin project management
- ✅ Admin portfolio management
- ✅ Cache invalidation endpoint
