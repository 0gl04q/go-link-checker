package output

import (
	"context"
	"fmt"

	"github.com/0gl04q/go-link-checker/internal/domain"
)

// ConsoleOutput - реализация интерфейса Output для вывода результатов в консоль
type ConsoleOutput struct{}

// NewConsoleOutput - конструктор для ConsoleOutput
func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

// Write - выводит результат в консоль
func (o *ConsoleOutput) Write(_ context.Context, l *domain.Link) error {
	fmt.Println(l.URL, l.Status, l.Message, l.Err, l.Timestamp)
	return nil
}
