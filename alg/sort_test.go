package alg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var qs = []int{9,4,1,5,2,8,7,3,6}
var qsSorted = []int{1,2,3,4,5,6,7,8,9}

func TestLomutoQuickSort(t *testing.T) {
	var a = qs
	LomutoQuickSort(a)
	assert.ObjectsAreEqual(qsSorted, a)
}

func TestHoareQuickSort(t *testing.T) {
	var a = qs
	HoareQuickSort(a)
	assert.ObjectsAreEqual(qsSorted, a)
}

func TestInsertionSort(t *testing.T) {
	var a = qs
	InsertionSort(a)
	assert.ObjectsAreEqual(qsSorted, a)
}

func BenchmarkLomutoQuickSort(b *testing.B) {
	var a = qs
	for i:=0; i <b.N; i++ {
		LomutoQuickSort(a)
	}
}

func BenchmarkHoareQuickSort(b *testing.B) {
	var a = qs
	for i:=0; i <b.N; i++ {
		HoareQuickSort(a)
	}
}

func BenchmarkInsertionSort(b *testing.B) {
	var a = []int{9,4,1,5,2,8,7,3,6}
	for i:=0; i <b.N; i++ {
		InsertionSort(a)
	}
}
