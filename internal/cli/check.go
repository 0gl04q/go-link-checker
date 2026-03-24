package cli

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/0gl04q/go-link-checker/internal/deduplicator"
	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/0gl04q/go-link-checker/internal/handler"
	"github.com/0gl04q/go-link-checker/internal/limiter"
	"github.com/0gl04q/go-link-checker/internal/output"
	"github.com/0gl04q/go-link-checker/pkg/consumer"
	"github.com/0gl04q/go-link-checker/pkg/producer"
	"github.com/0gl04q/go-link-checker/pkg/worker"
	"github.com/spf13/cobra"
)

const (
	redisOutputType   = "redis"
	consoleOutputType = "console"
)

// checkCmd - команда для проверки доступности ссылок, указанных в файле, и сохранения результатов в Redis
func (c *CLI) checkCmd() *cobra.Command {
	var (
		filePath       string
		workerPoolSize int
		outputType     string
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Проверить доступность ссылок",
		Long:  "Проверить доступность ссылок, указанных в конфигурации, и сохранить результаты в Redis",
		Run:   c.checkCmdInit(&filePath, &workerPoolSize, &outputType),
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "links.example.txt", "Путь к файлу со ссылками")
	cmd.Flags().IntVarP(&workerPoolSize, "worker-pool-size", "w", 100, "Размер пула воркеров для асинхронной проверки ссылок")
	cmd.Flags().StringVarP(&outputType, "output", "o", consoleOutputType, "Тип вывода результатов (console, redis)")

	return cmd
}

// checkCmdInit - функция для инициализации команды check, которая создает контекст с поддержкой сигналов и вызывает метод Check у linkUseCase
func (c *CLI) checkCmdInit(filePath *string, workerPoolSize *int, outputType *string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var (
			out    consumer.Output[domain.Link]
			dedup  producer.Deduplicator
			rLimit limiter.RateLimiter
		)

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		switch *outputType {
		case redisOutputType:
			out = output.NewRedisOutput(c.redisClient)
			dedup = deduplicator.NewRedisDeduplicator(c.redisClient)
			rLimit = limiter.NewRedisRateLimiter(c.redisClient)
		default:
			out = output.NewConsoleOutput()
			dedup = deduplicator.NewMemoryDeduplicator()
			rLimit = limiter.NewMemoryRateLimiter()
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        *workerPoolSize,
				MaxIdleConnsPerHost: *workerPoolSize / 10,
				IdleConnTimeout:     30 * time.Second,
				DisableKeepAlives:   false,
			},
		}

		c.linkUseCase.Check(
			ctx,
			*filePath,
			handler.NewLinkHandler(httpClient, rLimit),
			worker.NewPool[domain.Link](*workerPoolSize),
			producer.NewProducer(dedup),
			consumer.NewConsumer[domain.Link](out),
		)
	}
}
