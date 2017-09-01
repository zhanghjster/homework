package math

import (
	"testing"
)

func TestFib(t *testing.T) {
	if FibDesc(7) != 13 {
		t.Error("fib of 7 should be 13")
	}
}
