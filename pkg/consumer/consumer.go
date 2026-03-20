package consumer

// Output - интерфейс для вывода результатов
type Output[T any] interface {
	Write(*T) error
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
func (c *Consumer[T]) Consume(results <-chan *T) error {
	for result := range results {
		if err := c.output.Write(result); err != nil {
			return err
		}
	}

	return nil
}
