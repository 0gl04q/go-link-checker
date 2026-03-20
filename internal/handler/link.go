package handler

import (
	"fmt"
	"net/http"
	urlpkg "net/url"

	"github.com/0gl04q/go-link-checker/internal/domain"
)

// LinkHandler - базовый воркер для проверки ссылок
type LinkHandler struct{}

// NewLinkHandler - конструктор для LinkHandler
func NewLinkHandler() *LinkHandler {
	return &LinkHandler{}
}

// Handle - воркер для проверки ссылок, получает ссылки из канала jobs и отправляет результат в канал results
func (l *LinkHandler) Handle(jobs <-chan string, results chan<- domain.Result) {
	for link := range jobs {
		res, err := l.sendGetRequest(link)
		if err != nil {
			results <- domain.Result{Err: err}
			continue
		}

		if res.StatusCode >= 200 && res.StatusCode < 400 {
			results <- domain.Result{Message: fmt.Sprintf("%s ссылка доступна, статус %s \n", link, res.Status)}
		} else {
			results <- domain.Result{Message: fmt.Sprintf("%s ссылка не доступна, статус %s\n", link, res.Status)}
		}
	}
}

// sendGetRequest - отправляет GET запрос по ссылке и возвращает результат
func (l *LinkHandler) sendGetRequest(link string) (*http.Response, error) {
	_, err := urlpkg.ParseRequestURI(link)
	if err != nil {
		return nil, fmt.Errorf("%s ссылка не валидная\n", link)
	}

	r, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return r, nil
}
