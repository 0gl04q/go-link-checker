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
	"github.com/0gl04q/go-link-checker/pkg/consumer"
	"github.com/0gl04q/go-link-checker/pkg/producer"
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
func (l *LinkUseCase) Check(
	filePath string,
	httpClient *http.Client,
	wp *worker.Pool[domain.Link],
	prod *producer.Producer,
	con *consumer.Consumer[domain.Link],
) {
	l.handler.Client = httpClient

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	urls, err := l.getLinksFromFile(filePath)
	if err != nil {
		fmt.Printf("ошибка при получении ссылок: %v\n", err)
		return
	}

	results := wp.Start(ctx, l.handler.Handle)

	prod.Produce(ctx, urls, wp.Jobs)

	if errs := con.Consume(ctx, results); errs != nil {
		fmt.Printf("ошибок при потреблении: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
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
