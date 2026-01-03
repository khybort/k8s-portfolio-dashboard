# Redis Integration Documentation

## Genel Bakış

Portfolio Dashboard, caching ve session management için Redis kullanır. Redis, yüksek performanslı in-memory data store olarak, API response'larını cache'lemek ve kullanıcı session'larını yönetmek için kullanılır.

## Kullanım Senaryoları

### 1. API Response Caching

Frequently accessed data cache'lenir:

- **Article List**: Makale listesi (TTL: 5 dakika)
- **Article Detail**: Makale detayı (TTL: 10 dakika)
- **Project List**: Proje listesi (TTL: 5 dakika)
- **Project Detail**: Proje detayı (TTL: 10 dakika)
- **Portfolio Info**: Portfolio bilgileri (TTL: 1 saat)

### 2. Session Storage

- **JWT Refresh Tokens**: Refresh token'lar Redis'te saklanır
- **User Sessions**: Active user session'ları
- **Rate Limiting**: API rate limiting counters

### 3. Cache Invalidation

Kafka event'leri ile otomatik cache invalidation:

- Article oluşturuldu/güncellendi → Article list cache'i invalidate
- Project eklendi/güncellendi → Project list cache'i invalidate

## Cache Strategy

### Cache Key Patterns

```
articles:list:page:{page}:limit:{limit}
article:detail:{id}
article:slug:{slug}
projects:list:page:{page}:limit:{limit}
project:detail:{id}
portfolio:info
user:session:{user_id}
rate:limit:{ip}:{endpoint}
```

### TTL (Time To Live)

| Cache Type | TTL |
|------------|-----|
| Article List | 5 minutes |
| Article Detail | 10 minutes |
| Project List | 5 minutes |
| Project Detail | 10 minutes |
| Portfolio Info | 1 hour |
| User Session | 7 days |
| Rate Limit | 1 minute |

## Implementation

### Go Redis Client

```go
// internal/cache/redis.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
    return &RedisCache{
        client: redis.NewClient(&redis.Options{
            Addr:     addr,
            Password: password,
            DB:       db,
        }),
    }
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
    data, err := c.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, ErrCacheMiss
    }
    return data, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
    keys, err := c.client.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }
    if len(keys) > 0 {
        return c.client.Del(ctx, keys...).Err()
    }
    return nil
}
```

### Article Cache Example

```go
// internal/service/article.go
func (s *articleService) GetArticles(ctx context.Context, page, limit int) ([]model.Article, error) {
    // Cache key
    cacheKey := fmt.Sprintf("articles:list:page:%d:limit:%d", page, limit)
    
    // Try cache first
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil {
        var articles []model.Article
        json.Unmarshal(cached, &articles)
        return articles, nil
    }
    
    // Cache miss - fetch from database
    articles, err := s.repo.List(ctx, page, limit)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    data, _ := json.Marshal(articles)
    s.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    
    return articles, nil
}

func (s *articleService) CreateArticle(ctx context.Context, article *model.Article) error {
    // Save to database
    if err := s.repo.Create(ctx, article); err != nil {
        return err
    }
    
    // Invalidate cache
    s.cache.DeletePattern(ctx, "articles:list:*")
    
    return nil
}
```

### Session Storage

```go
// internal/cache/session.go
func (c *RedisCache) StoreSession(ctx context.Context, userID string, token string, ttl time.Duration) error {
    key := fmt.Sprintf("user:session:%s", userID)
    return c.client.Set(ctx, key, token, ttl).Err()
}

func (c *RedisCache) GetSession(ctx context.Context, userID string) (string, error) {
    key := fmt.Sprintf("user:session:%s", userID)
    return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) DeleteSession(ctx context.Context, userID string) error {
    key := fmt.Sprintf("user:session:%s", userID)
    return c.client.Del(ctx, key).Err()
}
```

### Rate Limiting

```go
// internal/middleware/ratelimit.go
func RateLimitMiddleware(cache *cache.RedisCache, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        endpoint := c.FullPath()
        key := fmt.Sprintf("rate:limit:%s:%s", ip, endpoint)
        
        count, err := cache.client.Incr(c.Request.Context(), key).Result()
        if err != nil {
            c.JSON(500, gin.H{"error": "Internal server error"})
            c.Abort()
            return
        }
        
        if count == 1 {
            cache.client.Expire(c.Request.Context(), key, window)
        }
        
        if count > int64(limit) {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## Cache Invalidation

### Manual Invalidation

```go
// Admin API endpoint
func (h *AdminHandler) InvalidateCache(c *gin.Context) {
    cacheType := c.Query("type")
    
    switch cacheType {
    case "articles":
        h.cache.DeletePattern(c.Request.Context(), "articles:*")
    case "projects":
        h.cache.DeletePattern(c.Request.Context(), "projects:*")
    case "all":
        h.cache.DeletePattern(c.Request.Context(), "*")
    }
    
    c.JSON(200, gin.H{"message": "Cache invalidated"})
}
```

### Automatic Invalidation (Kafka)

```go
// internal/kafka/cache_invalidator.go
func (c *CacheInvalidator) ConsumeArticles(ctx context.Context) error {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        
        var event Event
        json.Unmarshal(msg.Value, &event)
        
        switch event.EventType {
        case "article.created", "article.updated", "article.deleted":
            // Invalidate article caches
            c.cache.DeletePattern(ctx, "articles:*")
        }
    }
}
```

## Data Structures

### Strings

**Use Cases**: Simple key-value caching
```go
cache.Set(ctx, "article:1", articleJSON, 10*time.Minute)
```

### Hashes

**Use Cases**: Object caching with fields
```go
cache.client.HSet(ctx, "article:1", "title", "Article Title")
cache.client.HSet(ctx, "article:1", "content", "Article Content")
cache.client.Expire(ctx, "article:1", 10*time.Minute)
```

### Sets

**Use Cases**: Unique collections
```go
cache.client.SAdd(ctx, "article:tags:1", "go", "kubernetes", "docker")
```

### Sorted Sets

**Use Cases**: Ranked lists
```go
cache.client.ZAdd(ctx, "articles:popular", redis.Z{
    Score:  viewCount,
    Member: articleID,
})
```

## Local Development

### Docker Compose

```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data

volumes:
  redis-data:
```

### Redis CLI

```bash
# Connect to Redis
docker exec -it portfolio-redis redis-cli

# Test
PING
# Output: PONG

# Set value
SET article:1 "{\"id\":\"1\",\"title\":\"Test\"}"

# Get value
GET article:1

# Set with TTL
SETEX article:1 600 "{\"id\":\"1\",\"title\":\"Test\"}"

# Delete
DEL article:1

# Pattern delete
KEYS articles:*
DEL articles:*

# Exit
exit
```

## Production Considerations

### High Availability

**Redis Sentinel**:
- Master-Slave replication
- Automatic failover
- Multiple sentinel instances

**Redis Cluster**:
- Sharding across multiple nodes
- High availability
- Automatic failover

### Persistence

**RDB (Snapshot)**:
```conf
save 900 1
save 300 10
save 60 10000
```

**AOF (Append Only File)**:
```conf
appendonly yes
appendfsync everysec
```

### Memory Management

**Max Memory Policy**:
- `allkeys-lru`: Evict least recently used keys
- `volatile-lru`: Evict LRU keys with expiration
- `allkeys-random`: Evict random keys
- `noeviction`: Don't evict (return errors)

### Monitoring

**Metrics**:
- Memory usage
- Hit rate
- Commands per second
- Connected clients
- Keyspace size

**Redis Insight**:
```bash
docker run -d \
  -p 8001:8001 \
  redis/redisinsight:latest
```

## Best Practices

1. **Key Naming**: Consistent, descriptive key names
2. **TTL**: Appropriate TTL values
3. **Serialization**: JSON for complex objects
4. **Error Handling**: Handle cache misses gracefully
5. **Monitoring**: Monitor hit rate and memory usage
6. **Eviction Policy**: Configure appropriate eviction policy

## Troubleshooting

### High Memory Usage

```bash
# Check memory usage
redis-cli INFO memory

# Check largest keys
redis-cli --bigkeys

# Check key count
redis-cli DBSIZE
```

### Slow Operations

```bash
# Monitor slow commands
redis-cli SLOWLOG GET 10

# Monitor commands in real-time
redis-cli MONITOR
```

### Connection Issues

```bash
# Check connections
redis-cli INFO clients

# Check if Redis is responsive
redis-cli PING
```

