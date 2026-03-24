package cli

import (
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// checkCmd - команда для проверки доступности ссылок, указанных в файле, и сохранения результатов в Redis
func (c *CLI) reportCmd() *cobra.Command {
	var (
		outputType string
	)

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Получить отчет о проверенных ссылках из Redis",
		Long:  "Получить отчет о проверенных ссылках из Redis и вывести его в консоль или сохранить в файл",
		Run:   c.reportCmdInit(&outputType),
	}

	cmd.Flags().StringVarP(&outputType, "output", "o", consoleOutputType, "Тип вывода результатов (console, file)")

	return cmd
}

// reportCmdInit - функция для инициализации команды report, которая создает контекст с поддержкой сигналов и вызывает метод Report у linkUseCase
func (c *CLI) reportCmdInit(outputType *string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		c.linkUseCase.Report(ctx, c.redisClient)
	}
}
