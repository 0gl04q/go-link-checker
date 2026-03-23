package producer

import (
	"context"
	"fmt"
)

// Deduplicator - интерфейс для проверки, был ли URL уже обработан, и для очистки всех обработанных URL
type Deduplicator interface {
	IsSeen(ctx context.Context, url string) (bool, error)
	Clear(ctx context.Context) error
}

// Producer - структура, которая использует Deduplicator для проверки URL перед их обработкой
type Producer struct {
	deduplicator Deduplicator
}

// NewProducer - конструктор для Producer, который принимает Deduplicator
func NewProducer(deduplicator Deduplicator) *Producer {
	return &Producer{deduplicator: deduplicator}
}

// Produce - принимает список URL и канал для отправки URL на обработку. Для каждого URL проверяет, был ли он уже обработан с помощью Deduplicator. Если URL не был обработан, отправляет его в канал. После обработки всех URL очищает все обработанные URL с помощью Deduplicator и закрывает канал.
func (p *Producer) Produce(ctx context.Context, urls []string, jobs chan<- string) {
	go func() {
		for _, url := range urls {
			add, err := p.deduplicator.IsSeen(ctx, url)
			if err != nil {
				fmt.Printf("ошибка при дебаунсе ссылки %s: %v\n", url, err)
				continue
			}
			if !add {
				fmt.Printf("ссылка %s уже была проверена, пропускаем\n", url)
				continue
			}

			jobs <- url
		}

		err := p.deduplicator.Clear(ctx)
		if err != nil {
			return
		}
		
		close(jobs)
	}()
}
