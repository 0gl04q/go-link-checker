package report

import (
	"strconv"
	"time"

	"github.com/0gl04q/go-link-checker/internal/domain"
	"github.com/pterm/pterm"
)

// PrintReport - выводит отчёт о проверке ссылок в консоль с помощью pterm
func PrintReport(links []*domain.Link) {
	var ok, redirects, clientErr, serverErr, timeouts int
	tableData := pterm.TableData{
		{"URL", "Статус", "Сообщение", "Время"},
	}

	for _, l := range links {
		status := strconv.Itoa(l.Status)
		ts := time.Unix(l.Timestamp, 0).Format("15:04:05")

		switch {
		case l.Err != "":
			timeouts++
			tableData = append(tableData, []string{l.URL, "ERR", l.Err, ts})
		case l.Status >= 500:
			serverErr++
			tableData = append(tableData, []string{l.URL, status, l.Message, ts})
		case l.Status >= 400:
			clientErr++
			tableData = append(tableData, []string{l.URL, status, l.Message, ts})
		case l.Status >= 300:
			redirects++
			tableData = append(tableData, []string{l.URL, status, l.Message, ts})
		case l.Status >= 200:
			ok++
			tableData = append(tableData, []string{l.URL, status, l.Message, ts})
		}
	}

	total := len(links)

	pterm.DefaultHeader.Println("Отчёт проверки ссылок")

	pterm.DefaultSection.Println("Статистика")
	pterm.Success.Printf("Доступно (2xx):    %d (%.0f%%)\n", ok, percent(ok, total))
	pterm.Info.Printf("Редиректы (3xx):   %d (%.0f%%)\n", redirects, percent(redirects, total))
	pterm.Warning.Printf("Клиентские (4xx):  %d (%.0f%%)\n", clientErr, percent(clientErr, total))
	pterm.Error.Printf("Серверные (5xx):   %d (%.0f%%)\n", serverErr, percent(serverErr, total))
	pterm.Error.Printf("Таймауты:          %d (%.0f%%)\n", timeouts, percent(timeouts, total))

	pterm.DefaultSection.Println("Все ссылки")
	pterm.DefaultTable.
		WithHasHeader().
		WithData(tableData).
		Render()
}

func percent(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}
