package service

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func BenchmarkLinkUseCase_Check(b *testing.B) {
	linkUseCase := NewLinkUseCase()

	b.ResetTimer()

	linkUseCase.Check("../../links.example.txt", 100, "console")
}
