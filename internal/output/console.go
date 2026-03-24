package output

import (
	"context"

	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/pterm/pterm"
)

// ConsoleOutput - реализация интерфейса Output для вывода результатов в консоль
type ConsoleOutput struct{}

// NewConsoleOutput - конструктор для ConsoleOutput
func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

// Write - выводит результат в консоль
func (o *ConsoleOutput) Write(ctx context.Context, l *domain.Link) error {
	switch {
	case l.Err != "":
		pterm.Error.Printf("%s — %s\n", l.URL, l.Err)
	case l.Status >= 500:
		pterm.Fatal.Printf("[%d] %s\n", l.Status, l.URL)
	case l.Status >= 400:
		pterm.Warning.Printf("[%d] %s\n", l.Status, l.URL)
	case l.Status >= 300:
		pterm.Info.Printf("[%d] %s\n", l.Status, l.URL)
	case l.Status >= 200:
		pterm.Success.Printf("[%d] %s\n", l.Status, l.URL)
	}
	return nil
}
