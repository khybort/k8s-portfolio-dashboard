# Authentication Service Documentation

## Genel Bakış

Auth Service, Portfolio Dashboard sisteminin kimlik doğrulama ve yetkilendirme işlemlerini yöneten ayrı bir mikroservistir. JWT (JSON Web Token) tabanlı authentication kullanır.

## Mimari

### Servis Yapısı

```
auth-service/
├── cmd/
│   ├── server/           # Main application
│   └── migrate/          # Database migrations
├── internal/
│   ├── api/              # HTTP handlers
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   └── user.go
│   │   └── middleware/
│   │       └── auth.go
│   ├── service/          # Business logic
│   │   ├── auth.go
│   │   └── user.go
│   ├── repository/        # Data access
│   │   └── user.go
│   ├── model/             # Domain models
│   │   └── user.go
│   └── jwt/               # JWT utilities
│       ├── token.go
│       └── validator.go
└── migrations/
```

## API Endpoints

### Public Endpoints

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "name": "John Doe"
}
```

**Response** (201 Created):
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "User registered successfully"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "admin"
  }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

#### Verify Token
```http
POST /api/v1/auth/verify
Authorization: Bearer <access_token>
```

**Response** (200 OK):
```json
{
  "valid": true,
  "user_id": "uuid",
  "role": "admin",
  "expires_at": "2024-01-01T00:15:00Z"
}
```

#### Password Reset Request
```http
POST /api/v1/auth/password-reset/request
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Response** (200 OK):
```json
{
  "message": "Password reset email sent"
}
```

#### Password Reset
```http
POST /api/v1/auth/password-reset
Content-Type: application/json

{
  "token": "reset-token",
  "new_password": "NewSecurePassword123!"
}
```

**Response** (200 OK):
```json
{
  "message": "Password reset successfully"
}
```

### Protected Endpoints

#### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

**Response** (200 OK):
```json
{
  "message": "Logged out successfully"
}
```

## JWT Token Yapısı

### Access Token

**Payload**:
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "admin",
  "exp": 1704067200,
  "iat": 1704066300,
  "type": "access"
}
```

**Expiration**: 15 dakika
**Algorithm**: HS256

### Refresh Token

**Payload**:
```json
{
  "user_id": "uuid",
  "token_id": "uuid",
  "exp": 1704672000,
  "iat": 1704066300,
  "type": "refresh"
}
```

**Expiration**: 7 gün
**Storage**: Redis (TTL: 7 gün)

## Database Schema

### users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### refresh_tokens Table

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

## Security

### Password Hashing

- **Algorithm**: bcrypt
- **Cost**: 12 rounds
- **Salt**: Automatic (bcrypt handles)

### Token Security

- **Secret Key**: Environment variable'dan alınır
- **Token Rotation**: Refresh token rotation
- **Token Revocation**: Redis'te blacklist
- **HTTPS Only**: Production'da HTTPS zorunlu

### Rate Limiting

- **Login**: 5 attempts per 15 minutes per IP
- **Register**: 3 attempts per hour per IP
- **Password Reset**: 3 attempts per hour per email

## Integration

### Backend Service Integration

Backend servis, auth service'i token doğrulama için kullanır:

```go
// Backend middleware
func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        
        // Verify token with auth service
        resp, err := http.Post(
            authServiceURL + "/api/v1/auth/verify",
            "application/json",
            strings.NewReader(`{"token": "` + token + `"}`),
        )
        
        if err != nil || resp.StatusCode != 200 {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Extract user info from token
        // Set in context
        c.Next()
    }
}
```

### Frontend Integration

```typescript
// Frontend API client
const apiClient = axios.create({
  baseURL: 'http://localhost:8080',
});

// Request interceptor - Add token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - Handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Try to refresh token
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await axios.post(
            'http://localhost:8081/api/v1/auth/refresh',
            { refresh_token: refreshToken }
          );
          localStorage.setItem('access_token', response.data.access_token);
          // Retry original request
          return apiClient.request(error.config);
        } catch (refreshError) {
          // Redirect to login
          window.location.href = '/admin/login';
        }
      }
    }
    return Promise.reject(error);
  }
);
```

## Roles & Permissions

### Roles

- **admin**: Full access
- **user**: Limited access (read-only)

### Permission Matrix

| Action | Admin | User |
|--------|-------|------|
| View Articles | ✅ | ✅ |
| Create Article | ✅ | ❌ |
| Update Article | ✅ | ❌ |
| Delete Article | ✅ | ❌ |
| View Projects | ✅ | ✅ |
| Manage Projects | ✅ | ❌ |
| Manage Portfolio | ✅ | ❌ |

## Environment Variables

```bash
# Server
AUTH_SERVICE_PORT=8081
AUTH_SERVICE_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=portfolio
DB_PASSWORD=password
DB_NAME=auth_db

# JWT
JWT_SECRET=your-secret-key-here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Email (Password Reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-password
```

## Testing

### Unit Tests

```bash
cd auth-service
go test ./internal/service/...
go test ./internal/jwt/...
```

### Integration Tests

```bash
go test ./internal/api/... -tags=integration
```

### Manual Testing

```bash
# Register
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!","name":"Test User"}'

# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!"}'

# Verify Token
curl -X POST http://localhost:8081/api/v1/auth/verify \
  -H "Authorization: Bearer <token>"
```

## Monitoring

### Metrics

- Login attempts (success/failure)
- Token generation count
- Token verification count
- Password reset requests
- Active sessions

### Logging

- Authentication events
- Security events (failed login, token validation failures)
- Error logs

## Best Practices

1. **Never log passwords**: Password'ları asla loglama
2. **HTTPS in production**: Production'da mutlaka HTTPS kullan
3. **Token expiration**: Kısa expiration süreleri kullan
4. **Refresh token rotation**: Her refresh'te yeni token üret
5. **Rate limiting**: Brute force saldırılarına karşı rate limiting
6. **Password policy**: Güçlü password policy uygula

