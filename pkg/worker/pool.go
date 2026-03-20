package worker

import (
	"context"
	"sync"
)

// Handler - функция-обработчик, которая принимает канал с заданиями и канал для отправки результатов
type Handler[T any] func(ctx context.Context, jobs <-chan string, results chan<- *T)

// Pool - структура, которая управляет количеством воркеров и каналом для заданий
type Pool[T any] struct {
	Size int
	Jobs chan string
}

// NewPool - конструктор для Pool, который инициализирует канал для заданий
func NewPool[T any](size int) *Pool[T] {
	return &Pool[T]{
		Size: size,
		Jobs: make(chan string, size),
	}
}

// Start - запускает воркеры и возвращает канал для получения результатов
func (w *Pool[T]) Start(ctx context.Context, handler Handler[T]) <-chan *T {
	var wg sync.WaitGroup

	results := make(chan *T, w.Size)

	wg.Add(w.Size)

	for i := 0; i < w.Size; i++ {
		go func() {
			defer wg.Done()
			handler(ctx, w.Jobs, results)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
