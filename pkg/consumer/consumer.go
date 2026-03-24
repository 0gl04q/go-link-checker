package consumer

import (
	"context"
)

// Output - интерфейс для вывода результатов
type Output[T any] interface {
	Write(ctx context.Context, result *T) error
}

// Consumer - отвечает за потребление результатов из канала и передачу их в Output
type Consumer[T any] struct {
	output Output[T]
}

// NewConsumer - конструктор для Consumer
func NewConsumer[T any](output Output[T]) *Consumer[T] {
	return &Consumer[T]{
		output: output,
	}
}

// Consume - потребляет результаты из канала и передает их в Output
func (c *Consumer[T]) Consume(ctx context.Context, results <-chan *T) []error {
	var errs []error

	for {
		select {
		case <-ctx.Done():
			return errs
		case result, ok := <-results:
			if !ok {
				return errs
			}
			if err := c.output.Write(ctx, result); err != nil {
				errs = append(errs, err)
			}
		}
	}
}
