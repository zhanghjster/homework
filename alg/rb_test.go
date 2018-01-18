package alg

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestRabinKarp(t *testing.T) {
	var txt = "hello from mars"
	var pat = "mars"
	i := RabinKarp(txt, pat, 256, 103)
	assert.Equal(t, i, 11)
}
