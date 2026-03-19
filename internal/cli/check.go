package cli

import "github.com/spf13/cobra"

func (c *CLI) checkCmd() *cobra.Command {
	var (
		filePath   string
		dryRun     bool
		syncMethod bool
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Проверить доступность ссылок",
		Long:  "Проверить доступность ссылок, указанных в конфигурации, и сохранить результаты в Redis",
		Run: func(cmd *cobra.Command, args []string) {
			c.linkUseCase.CheckLinks(filePath, dryRun, syncMethod)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Путь к файлу со ссылками (по умолчанию используется конфигурация)")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Показать результаты проверки без сохранения в Redis")
	cmd.Flags().BoolVarP(&syncMethod, "sync", "s", false, "Проверять ссылки синхронно (без многопоточности)")

	return cmd
}
