package service

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/0gl04q/go-link-checker/internal/handler"
	"github.com/0gl04q/go-link-checker/internal/output"
	"github.com/0gl04q/go-link-checker/pkg/consumer"
	"github.com/0gl04q/go-link-checker/pkg/worker"
)

// LinkUseCase - базовый юзкейс для проверки ссылок
type LinkUseCase struct {
	handler *handler.LinkHandler
}

// NewLinkUseCase - конструктор для LinkUseCase
func NewLinkUseCase() *LinkUseCase {
	return &LinkUseCase{
		handler: handler.NewLinkHandler(),
	}
}

// Check - проверяет доступность ссылок и выводит результат
func (l *LinkUseCase) Check(filePath string, workerPoolSize int) {
	l.handler.Client = l.buildClient(workerPoolSize)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	links, err := l.getLinksFromFile(filePath)
	if err != nil {
		fmt.Printf("ошибка при получении ссылок: %v\n", err)
		return
	}

	pool := worker.NewPool[domain.Result](workerPoolSize)
	results := pool.Start(ctx, l.handler.Handle)

	go func() {
		for _, link := range links {
			pool.Jobs <- link
		}
		close(pool.Jobs)
	}()

	c := consumer.NewConsumer[domain.Result](output.NewConsoleOutput())
	if err := c.Consume(results); err != nil {
		fmt.Printf("ошибка при потреблении результатов: %v\n", err)
		return
	}
}

// buildClient - создает HTTP клиент с оптимизированными параметрами для работы с большим количеством ссылок
func (l *LinkUseCase) buildClient(workerPoolSize int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        workerPoolSize,
			MaxIdleConnsPerHost: workerPoolSize / 10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
		},
	}
}

// getLinksFromFile - читает ссылки из файла и возвращает их в виде слайса строк
func (l *LinkUseCase) getLinksFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	var links []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла: %v", err)
	}

	return links, nil
}
