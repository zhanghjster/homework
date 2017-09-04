package math

import (
	"testing"
)

func TestIsPrimitive(t *testing.T) {
	var p, np uint64 = 13, 30

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
	FindAllPrimitive(uint64(b.N))
}


