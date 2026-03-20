package cli

import "github.com/spf13/cobra"

// checkCmd - команда для проверки доступности ссылок, указанных в файле, и сохранения результатов в Redis
func (c *CLI) checkCmd() *cobra.Command {
	var (
		filePath       string
		workerPoolSize int
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Проверить доступность ссылок",
		Long:  "Проверить доступность ссылок, указанных в конфигурации, и сохранить результаты в Redis",
		Run: func(cmd *cobra.Command, args []string) {
			c.linkUseCase.Check(filePath, workerPoolSize)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "links.example.txt", "Путь к файлу со ссылками")
	cmd.Flags().IntVarP(&workerPoolSize, "worker-pool-size", "w", 100, "Размер пула воркеров для асинхронной проверки ссылок")

	return cmd
}
