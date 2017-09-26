package algrithm

import "testing"

func TestPow(t *testing.T) {

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
