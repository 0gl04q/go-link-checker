package output

import (
	"context"

	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/redis/go-redis/v9"
)

// RedisOutput - реализация интерфейса Output для сохранения результатов в Redis
type RedisOutput struct {
	client *redis.Client
}

// NewRedisOutput - конструктор для RedisOutput, который принимает клиент Redis
func NewRedisOutput(client *redis.Client) *RedisOutput {
	return &RedisOutput{client: client}
}

// Write - сохраняет результат в Redis
func (o *RedisOutput) Write(ctx context.Context, l *domain.Link) error {
	return o.client.HSet(ctx, "link:"+l.URL, l).Err()
}
