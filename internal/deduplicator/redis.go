package deduplicator

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// urlsSetKey Ключ для Redis Set, который будет использоваться для хранения уникальных URL
const urlsSetKey = "search_urls"

// RedisDeduplicator - реализация интерфейса Debouncer, использующая Redis для хранения обработанных URL
type RedisDeduplicator struct {
	r *redis.Client
}

// NewRedisDeduplicator - конструктор для RedisDeduplicator, который принимает клиент Redis
func NewRedisDeduplicator(r *redis.Client) *RedisDeduplicator {
	return &RedisDeduplicator{r: r}
}

// IsSeen - проверяет, был ли URL уже обработан, используя Redis Set для хранения уникальных URL. Если URL уже существует в Set, возвращает true, иначе добавляет его и возвращает false.
func (d *RedisDeduplicator) IsSeen(ctx context.Context, url string) (bool, error) {
	added, err := d.r.SAdd(ctx, urlsSetKey, url).Result()
	if err != nil {
		return false, err
	}
	return added == 1, nil
}

// Clear - очищает все обработанные URL из Redis, удаляя Set
func (d *RedisDeduplicator) Clear(ctx context.Context) error {
	return d.r.Del(ctx, urlsSetKey).Err()
}
