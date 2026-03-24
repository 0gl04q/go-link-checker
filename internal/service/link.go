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
	"github.com/0gl04q/go-link-checker/internal/report"
	"github.com/0gl04q/go-link-checker/internal/timer"
	"github.com/0gl04q/go-link-checker/pkg/consumer"
	"github.com/0gl04q/go-link-checker/pkg/producer"
	"github.com/0gl04q/go-link-checker/pkg/worker"
	"github.com/redis/go-redis/v9"
)

// LinkUseCase - базовый юзкейс для проверки ссылок
type LinkUseCase struct{}

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
	defer timer.Duration(timer.Track("Проверка ссылок завершена"))

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	urls, err := l.getLinksFromFile(filePath)
	if err != nil {
		slog.Error("Ошибка при получении ссылок", "err", err)
		return
	}

	slog.Info("Получено ссылок для проверки", "count", len(urls))

	results := wp.Start(ctx, handler.Handle)

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

// Report - выводит отчёт по результатам проверки ссылок, получая данные из Redis
func (l *LinkUseCase) Report(ctx context.Context, r *redis.Client) {
	links, err := l.getAllLinks(ctx, r)
	if err != nil {
		slog.Error("ошибка при получении ссылок из Redis", "err", err)
	}
	report.PrintReport(links)
}

// Clear - удаляет все данные о проверенных ссылках из Redis, используя SCAN для итерации по ключам и удаления их пакетами
func (l *LinkUseCase) Clear(ctx context.Context, r *redis.Client) {
	var cursor uint64
	for {
		keys, nextCursor, err := r.Scan(ctx, cursor, "link:*", 100).Result()
		if err != nil {
			slog.Error("ошибка при сканировании ключей в Redis", "err", err)
			return
		}
		if len(keys) > 0 {
			r.Del(ctx, keys...)
		}
		cursor = nextCursor
		if cursor == 0 {
			break
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

// getAllLinks - получает все ссылки из Redis и возвращает их в виде слайса структур Link
func (l *LinkUseCase) getAllLinks(ctx context.Context, r *redis.Client) ([]*domain.Link, error) {
	var keys []string
	var cursor uint64

	for {
		batch, nextCursor, err := r.Scan(ctx, cursor, "link:*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании ключей в Redis: %v", err)
		}

		keys = append(keys, batch...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	pipe := r.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.HGetAll(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении pipeline в Redis: %v", err)
	}

	var links []*domain.Link

	for _, cmd := range cmds {
		var link domain.Link
		if err := cmd.Scan(&link); err != nil {
			slog.Error("Ошибка при сканировании данных из Redis в структуру Link", "err", err)
			continue
		}
		links = append(links, &link)
	}

	return links, nil
}
