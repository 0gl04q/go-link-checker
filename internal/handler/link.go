package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	urlpkg "net/url"
	"time"

	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/0gl04q/go-link-checker/internal/limiter"
)

// ErrEmptyResult - ошибка для случаев, когда результат проверки ссылки пустой
var ErrEmptyResult = fmt.Errorf("пустой результат")

// LinkHandler - базовый воркер для проверки ссылок
type LinkHandler struct {
	RateLimiter limiter.RateLimiter
	Client      *http.Client
}

// NewLinkHandler - конструктор для LinkHandler
func NewLinkHandler(c *http.Client, l limiter.RateLimiter) *LinkHandler {
	return &LinkHandler{Client: c, RateLimiter: l}
}

// Handle - воркер для проверки ссылок, получает ссылки из канала jobs и отправляет результат в канал results
func (l *LinkHandler) Handle(ctx context.Context, jobs <-chan string, results chan<- *domain.Link) {
	for link := range jobs {
		results <- l.processLink(ctx, link)
	}
}

// processLink - обрабатывает ссылку и возвращает результат проверки
func (l *LinkHandler) processLink(ctx context.Context, link string) *domain.Link {
	timestamp := time.Now().Unix()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := l.sendGetRequest(ctx, link)
	if err != nil {
		slog.Error("Ошибка при отправке GET запроса", "url", link, "err", err)
		return domain.NewLink(link, 0, "", timestamp, err)
	}
	if res == nil {
		slog.Error("Пустой результат при отправке GET запроса", "url", link)
		return domain.NewLink(link, 0, "", timestamp, ErrEmptyResult)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 400 {
		slog.Debug("Ссылка доступна", "url", link, "status_code", res.StatusCode)
		return domain.NewLink(link, res.StatusCode, "Ссылка доступна", timestamp, nil)
	}

	slog.Debug("Ссылка не доступна", "url", link, "status_code", res.StatusCode)
	return domain.NewLink(link, res.StatusCode, "Ссылка не доступна", timestamp, nil)
}

// sendGetRequest - отправляет GET запрос по ссылке и возвращает результат
func (l *LinkHandler) sendGetRequest(ctx context.Context, url string) (*http.Response, error) {
	u, parsErr := urlpkg.ParseRequestURI(url)
	if parsErr != nil {
		return nil, fmt.Errorf("%s ссылка не валидная\n", url)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		allowed, err := l.RateLimiter.Allow(ctx, u.Hostname())
		if err != nil {
			return nil, fmt.Errorf("ошибка при проверке rate limiter для хоста %s: %v\n", u.Hostname(), err)
		}

		if allowed {
			break
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("контекст истек при ожидании разрешения от rate limiter для хоста %s\n", u.Hostname())
		case <-ticker.C:
		}
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkChecker/1.0)")

	for i := 0; i < 3; i++ {
		response, err := l.Client.Do(r)
		if err == nil {
			return response, nil
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("контекст истек при отправке запроса к %s: %v\n", url, err)
		}

		if i == 2 {
			break
		}

		delay := 500 * time.Millisecond * time.Duration(1<<i)
		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("контекст истек при ожидании следующей попытки для %s: %v\n", url, ctx.Err())
		case <-time.After(delay + jitter):
		}
	}

	return nil, fmt.Errorf("не удалось получить ответ от %s после 3 попыток", url)
}
