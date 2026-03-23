package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const rateLimitKeyPrefix = "rate_limit:"
const rateLimit = 5

// RateLimiter - интерфейс для ограничения количества запросов к одному хосту
type RateLimiter interface {
	Allow(ctx context.Context, hostname string) (bool, error)
}

// RedisRateLimiter - реализация rateLimiter, использующая Redis для хранения счетчика запросов к каждому хосту
type RedisRateLimiter struct {
	r *redis.Client
}

// NewRedisRateLimiter - конструктор для RedisRateLimiter, который принимает Redis клиент
func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{r: client}
}

// Allow - проверяет, можно ли отправлять запрос к данному хосту, увеличивая счетчик запросов в Redis и устанавливая время жизни ключа на 1 секунду при первом запросе
func (l *RedisRateLimiter) Allow(ctx context.Context, hostname string) (bool, error) {
	key := fmt.Sprintf("%s%s", rateLimitKeyPrefix, hostname)

	count, err := l.r.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = l.r.Expire(ctx, key, 1*time.Second).Err()
	}

	return count <= int64(rateLimit), err
}
