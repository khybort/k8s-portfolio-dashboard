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

func (c *RedisCache) GetArticle(ctx context.Context, id string) ([]byte, error) {
	key := fmt.Sprintf("article:detail:%s", id)
	return c.Get(ctx, key)
}

func (c *RedisCache) SetArticle(ctx context.Context, id string, article interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("article:detail:%s", id)
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, data, ttl)
}

func (c *RedisCache) InvalidateArticles(ctx context.Context) error {
	return c.DeletePattern(ctx, "articles:*")
}

var ErrCacheMiss = fmt.Errorf("cache miss")

