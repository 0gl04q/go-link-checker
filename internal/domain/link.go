package domain

// Result - структура для хранения результатов проверки ссылок
type Result struct {
	Link    string
	Status  int
	Message string
	Err     error
}

func NewResult(link string, status int, message string, err error) *Result {
	return &Result{
		Link:    link,
		Status:  status,
		Message: message,
		Err:     err,
	}
}
