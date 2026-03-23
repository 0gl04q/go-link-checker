package domain

// Link - структура для хранения результатов проверки ссылок
type Link struct {
	URL       string `redis:"url"`
	Status    int    `redis:"status"`
	Message   string `redis:"message"`
	Timestamp int64  `redis:"timestamp"`
	Err       string `redis:"err"`
}

// NewLink - конструктор для Link, который инициализирует все поля
func NewLink(url string, status int, message string, timestamp int64, err error) *Link {
	errStr := ""
	if err != nil {
		errStr = err.Error() // ← один раз при создании
	}
	return &Link{
		URL:       url,
		Status:    status,
		Message:   message,
		Timestamp: timestamp,
		Err:       errStr,
	}
}
