package handler

import (
	"context"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"time"

	"github.com/0gl04q/go-link-checker/internal/domain"
)

var ErrEmptyResult = fmt.Errorf("пустой результат")

// LinkHandler - базовый воркер для проверки ссылок
type LinkHandler struct {
	Client *http.Client
}

// NewLinkHandler - конструктор для LinkHandler
func NewLinkHandler() *LinkHandler {
	return &LinkHandler{}
}

// Handle - воркер для проверки ссылок, получает ссылки из канала jobs и отправляет результат в канал results
func (l *LinkHandler) Handle(ctx context.Context, jobs <-chan string, results chan<- *domain.Result) {
	for link := range jobs {
		results <- l.processLink(ctx, link)
	}
}

// processLink - обрабатывает ссылку и возвращает результат проверки
func (l *LinkHandler) processLink(ctx context.Context, link string) *domain.Result {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	res, err := l.sendGetRequest(ctx, link)
	if err != nil {
		return domain.NewResult(link, 0, "", err)
	}
	if res == nil {
		return domain.NewResult(link, 0, "", ErrEmptyResult)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 400 {
		return domain.NewResult(link, res.StatusCode, "Ссылка доступна", nil)
	}

	return domain.NewResult(link, res.StatusCode, "Ссылка не доступна", nil)
}

// sendGetRequest - отправляет GET запрос по ссылке и возвращает результат
func (l *LinkHandler) sendGetRequest(ctx context.Context, link string) (*http.Response, error) {
	_, err := urlpkg.ParseRequestURI(link)
	if err != nil {
		return nil, fmt.Errorf("%s ссылка не валидная\n", link)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(r)
}
