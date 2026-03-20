package output

import (
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
func (o *ConsoleOutput) Write(r domain.Result) error {
	if r.Err != nil {
		fmt.Printf("%v\n", r.Err)
	} else {
		fmt.Print(r.Message)
	}
	return nil
}
