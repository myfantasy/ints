package ints

import (
	"testing"
)

func BenchmarkUInt128Rand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandUUID()
	}
}
