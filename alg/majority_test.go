package alg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindMajority(t *testing.T) {
	var items = []int{1, 5, 2, 5, 3, 4, 5, 5, 6, 5, 5, 5}
	m := FindMajority(items)
	assert.Equal(t, m, 5)

}
