package cli

import (
	"github.com/0gl04q/go-link-checker/internal/config"
	"github.com/0gl04q/go-link-checker/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

// Services - структура для хранения всех сервисов, которые будут использоваться в CLI
type Services struct {
	linkUseCase *service.LinkUseCase
}

// CLI - структура для хранения всех данных, необходимых для работы CLI
type CLI struct {
	Services

	redisClient *redis.Client

	root *cobra.Command
	cfg  *config.Config
}

// New - конструктор для CLI, который инициализирует все сервисы и команды
func New(cfg *config.Config) *CLI {
	c := &CLI{cfg: cfg}

	c.root = &cobra.Command{
		Use:   "go-link-checker",
		Short: "Проверка доступности ссылок",
		Long:  "Проверка доступности ссылок с помощью многопоточности и сохранение результатов в Redis",
	}

	c.linkUseCase = service.NewLinkUseCase()

	c.redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	c.root.AddCommand(c.checkCmd())
	c.root.AddCommand(c.reportCmd())
	c.root.AddCommand(c.clearCmd())

	return c
}

// Run - метод для запуска CLI, который выполняет корневую команду и обрабатывает ошибки
func (c *CLI) Run() error {
	return c.root.Execute()
}
