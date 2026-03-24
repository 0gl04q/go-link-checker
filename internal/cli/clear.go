package cli

import (
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// checkCmd - команда для проверки доступности ссылок, указанных в файле, и сохранения результатов в Redis
func (c *CLI) clearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Очистить все данные о проверенных ссылках из Redis",
		Long:  "Очистить все данные о проверенных ссылках из Redis, удалив все ключи, связанные с проверкой ссылок",
		Run:   c.clearCmdInit(),
	}
	return cmd
}

// clearCmdInit - функция для инициализации команды clear, которая создает контекст с поддержкой сигналов и вызывает метод Clear у linkUseCase
func (c *CLI) clearCmdInit() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		c.linkUseCase.Clear(ctx, c.redisClient)
	}
}
