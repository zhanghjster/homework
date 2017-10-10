package alg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	var n = Pow(2,3)
	assert.Equal(t, n, 8)

}

func BenchmarkPow(b *testing.B) {
	for i :=0; i < b.N; i++ {
		Pow(2,10)
	}
}

func BenchmarkPowEasy(b *testing.B) {
	for i :=0; i < b.N; i++ {
	}
}
