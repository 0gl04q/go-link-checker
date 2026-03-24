package service

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/0gl04q/go-link-checker/internal/deduplicator"
	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/0gl04q/go-link-checker/internal/handler"
	"github.com/0gl04q/go-link-checker/internal/limiter"
	"github.com/0gl04q/go-link-checker/internal/output"
	"github.com/0gl04q/go-link-checker/pkg/consumer"
	"github.com/0gl04q/go-link-checker/pkg/producer"
	"github.com/0gl04q/go-link-checker/pkg/worker"
	"github.com/redis/go-redis/v9"
)

func BenchmarkLinkUseCase_Check(b *testing.B) {
	linkUseCase := NewLinkUseCase()
	wps := 1000

	b.ResetTimer()

	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	out := output.NewRedisOutput(rClient)
	dedup := deduplicator.NewRedisDeduplicator(rClient)
	rLimit := limiter.NewRedisRateLimiter(rClient)

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        wps,
			MaxIdleConnsPerHost: wps / 10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	linkUseCase.Check(
		ctx,
		"../../links.example.txt",
		handler.NewLinkHandler(httpClient, rLimit),
		worker.NewPool[domain.Link](wps),
		producer.NewProducer(dedup),
		consumer.NewConsumer[domain.Link](out),
	)
}
