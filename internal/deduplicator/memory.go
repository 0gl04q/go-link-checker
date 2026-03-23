package deduplicator

import (
	"context"
	"sync"
)

// MemoryDeduplicator - реализация интерфейса Deduplicator, которая хранит обработанные URL в памяти
type MemoryDeduplicator struct {
	mu            sync.Mutex
	processedURLs map[string]struct{}
}

// NewMemoryDeduplicator - конструктор для MemoryDeduplicator, который инициализирует срез для хранения обработанных URL
func NewMemoryDeduplicator() *MemoryDeduplicator {
	return &MemoryDeduplicator{
		processedURLs: make(map[string]struct{}),
	}
}

// IsSeen - проверяет, был ли URL уже обработан, и если нет, добавляет его в срез обработанных URL
func (d *MemoryDeduplicator) IsSeen(_ context.Context, url string) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.processedURLs[url]; !ok {
		d.processedURLs[url] = struct{}{}
		return true, nil
	}

	return false, nil
}

// Clear - очищает все обработанные URL, сбрасывая срез к начальной пустой длине
func (d *MemoryDeduplicator) Clear(_ context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.processedURLs = make(map[string]struct{})
	return nil
}
