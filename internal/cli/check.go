package cli

import (
	"net/http"
	"time"

	"github.com/0gl04q/go-link-checker/internal/deduplicator"
	"github.com/0gl04q/go-link-checker/internal/domain"
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

func (c *CLI) checkCmdInit(filePath *string, workerPoolSize *int, outputType *string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var out consumer.Output[domain.Link]
		var dedup producer.Deduplicator

		switch *outputType {
		case redisOutputType:
			out = output.NewRedisOutput(c.redisClient)
			dedup = deduplicator.NewRedisDeduplicator(c.redisClient)
		default:
			out = output.NewConsoleOutput()
			dedup = deduplicator.NewMemoryDeduplicator()
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
			*filePath,
			httpClient,
			worker.NewPool[domain.Link](*workerPoolSize),
			producer.NewProducer(dedup),
			consumer.NewConsumer[domain.Link](out),
		)
	}
}
