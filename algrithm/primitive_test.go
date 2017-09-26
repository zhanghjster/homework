package algrithm

import (
	"testing"
)

func TestIsPrimitive(t *testing.T) {
	var p, np uint64 = 16777619, 30

	if !IsPrimitive(p) {
		t.Fatal(p, "should be primitive")
	}

	if IsPrimitive(np) {
		t.Fatal(np, "should not be primitive")
	}
}

func TestFindAllPrimitive(t *testing.T) {
	//FindAllPrimitive(1000)
}

func BenchmarkFindAllPrimitive(b *testing.B) {
	EratosthenesSieve(uint64(b.N))
}


