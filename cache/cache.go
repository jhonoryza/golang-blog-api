package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheManager struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheManager(client *redis.Client, ttl time.Duration) *CacheManager {
	return &CacheManager{
		client: client,
		ttl:    ttl,
	}
}

func (c *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[CACHE] MISS: %s", key)
		return nil
	}
	if err != nil {
		log.Printf("[CACHE] ERROR getting %s: %v", key, err)
		return err
	}
	log.Printf("[CACHE] HIT: %s", key)
	return json.Unmarshal([]byte(val), dest)
}

func (c *CacheManager) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		log.Printf("[CACHE] ERROR setting %s: %v", key, err)
		return err
	}
	log.Printf("[CACHE] SET: %s", key)
	return nil
}

func HasSearchOrSortParams(r *http.Request) bool {
	values := r.URL.Query()
	_, hasSearch := values["search"]
	_, hasSortBy := values["sortBy"]
	_, hasSortDir := values["sortDir"]
	return hasSearch || hasSortBy || hasSortDir
}

func BuildSurahCacheKey(r *http.Request) string {
	return "surah:all"
}

func BuildAyahCacheKey(r *http.Request, surahId string) string {
	return fmt.Sprintf("ayah:%s", surahId)
}