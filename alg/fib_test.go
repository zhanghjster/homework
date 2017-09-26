package alg

import (
	"testing"
)

func TestFibDesc(t *testing.T) {
	if FibDesc(7) != 13 {
		t.Error("fib of 7 should be 13")
	}
}

func TestFibAsc(t *testing.T) {
	if FibAsc(7) != 13 {
		t.Error("fib of 7 should be 13")
	}
}

func BenchmarkFibAsc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibDesc(10)
	}
}

func BenchmarkFibDesc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibAsc(10)
	}
}
