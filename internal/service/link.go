package service

import (
	"bufio"
	"fmt"
	"os"

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
	links, err := l.getLinksFromFile(filePath)
	if err != nil {
		fmt.Printf("ошибка при получении ссылок: %v\n", err)
		return
	}

	pool := worker.NewPool[domain.Result](workerPoolSize)
	results := pool.Start(l.handler.Handle)

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
