package service

import (
	"testing"
)

func BenchmarkLinkUseCase_Check(b *testing.B) {
	linkUseCase := NewLinkUseCase()

	b.ResetTimer()

	linkUseCase.Check("../../links.example.txt", 1000)
}
