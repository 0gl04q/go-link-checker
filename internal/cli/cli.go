package cli

import (
	"github.com/0gl04q/go-link-checker/internal/config"
	"github.com/0gl04q/go-link-checker/internal/service"
	"github.com/spf13/cobra"
)

type Services struct {
	linkUseCase *service.LinkUseCase
}

type CLI struct {
	Services

	root *cobra.Command
	cfg  *config.Config
}

func New(cfg *config.Config) *CLI {
	c := &CLI{cfg: cfg}

	c.root = &cobra.Command{
		Use:   "go-link-checker",
		Short: "Проверка доступности ссылок",
		Long:  "Проверка доступности ссылок с помощью многопоточности и сохранение результатов в Redis",
	}

	c.linkUseCase = service.NewLinkUseCase()

	c.root.AddCommand(c.checkCmd())

	return c
}

func (c *CLI) Run() error {
	return c.root.Execute()
}
