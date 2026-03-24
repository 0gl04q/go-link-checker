package service

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
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
	return &LinkUseCase{}
}

// Check - проверяет доступность ссылок и выводит результат
func (l *LinkUseCase) Check(
	ctx context.Context,
	filePath string,
	handler *handler.LinkHandler,
	wp *worker.Pool[domain.Link],
	prod *producer.Producer,
	con *consumer.Consumer[domain.Link],
) {
	l.handler = handler

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	urls, err := l.getLinksFromFile(filePath)
	if err != nil {
		slog.Error("Ошибка при получении ссылок", "err", err)
		return
	}

	slog.Info("Получено ссылок для проверки", "count", len(urls))

	results := wp.Start(ctx, l.handler.Handle)

	slog.Info("Запустили воркеры")

	prod.Produce(ctx, urls, wp.Jobs)

	slog.Info("Запустили продюсера")

	if errs := con.Consume(ctx, results); errs != nil {
		slog.Error("Ошибок при потреблении ресурсов")
		for _, err := range errs {
			slog.Error("err", err)
		}
	}

	slog.Info("Завершили проверку ссылок")
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
