package utils

import (
	"math/rand"
	"testing"
)

func BenchmarkPow(b *testing.B) {
	x := rand.Intn(100)
	y := rand.Intn(6)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		PowInt(x, y)
	}
}
