package service

import (
	"bufio"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"os"
	"sync"
)

type result struct {
	message string
	err     error
}

// LinkUseCase - базовый юзкейс для проверки ссылок
type LinkUseCase struct{}

// NewLinkUseCase - конструктор для LinkUseCase
func NewLinkUseCase() *LinkUseCase {
	return &LinkUseCase{}
}

// CheckLinks - проверяет доступность ссылок и выводит результат
func (l *LinkUseCase) CheckLinks(filePath string, dryRun, syncMethod bool) {
	links, err := l.getLinksFromFile(filePath)
	if err != nil {
		fmt.Printf("ошибка при получении ссылок: %v\n", err)
		return
	}

	switch syncMethod {
	case true:
		l.syncCheck(links)
	case false:
		l.asyncCheck(links)
	}
}

// syncCheck - проверяет доступность ссылок синхронно и выводит результат
func (l *LinkUseCase) syncCheck(links []string) {
	for _, link := range links {
		res, err := l.sendGetRequest(link)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		if res.StatusCode >= 200 && res.StatusCode < 400 {
			fmt.Printf("%s ссылка доступна, статус %s \n", link, res.Status)
		} else {
			fmt.Printf("%s ссылка не доступна, статус %s\n", link, res.Status)
		}
	}
}

// asyncCheck - проверяет доступность ссылок асинхронно и выводит результат
func (l *LinkUseCase) asyncCheck(links []string) {
	var wg sync.WaitGroup
	ch := make(chan result, len(links))

	wg.Add(len(links))

	for _, link := range links {
		go l.checkLink(&wg, ch, link)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		if res.err != nil {
			fmt.Printf("%v\n", res.err)
		} else {
			fmt.Print(res.message)
		}
	}
}

// sendGetRequest - отправляет GET запрос по ссылке и возвращает результат
func (l *LinkUseCase) sendGetRequest(link string) (*http.Response, error) {
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

// checkLink - проверяет доступность ссылки и отправляет результат в канал
func (l *LinkUseCase) checkLink(wg *sync.WaitGroup, ch chan<- result, link string) {
	defer wg.Done()

	r, err := l.sendGetRequest(link)
	if err != nil {
		ch <- result{"", err}
		return
	}

	if r.StatusCode >= 200 && r.StatusCode < 400 {
		ch <- result{fmt.Sprintf("%s ссылка доступна, статус %s \n", link, r.Status), nil}
	} else {
		ch <- result{fmt.Sprintf("%s ссылка не доступна, статус %s\n", link, r.Status), nil}
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
