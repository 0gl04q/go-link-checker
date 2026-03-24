package limiter

import (
	"context"
	"sync"
	"time"
)

// MemoryRateLimiter - реализация rateLimiter, использующая встроенную мапу для хранения счетчика запросов к каждому хосту
type MemoryRateLimiter struct {
	mu    sync.Mutex
	store map[string]int
}

// NewMemoryRateLimiter - конструктор для RedisRateLimiter, который принимает Redis клиент
func NewMemoryRateLimiter() *MemoryRateLimiter {
	rl := &MemoryRateLimiter{
		store: make(map[string]int),
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			rl.mu.Lock()
			rl.store = make(map[string]int)
			rl.mu.Unlock()
		}
	}()

	return rl
}

// Allow - проверяет, можно ли отправлять запрос к данному хосту, увеличивая счетчик запросов в хранилище и устанавливая время жизни ключа на 1 секунду при первом запросе
func (l *MemoryRateLimiter) Allow(_ context.Context, hostname string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.store[hostname]++

	return l.store[hostname] <= rateLimit, nil
}
