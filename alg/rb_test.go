package alg

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRabinKarp(t *testing.T) {
	var txt = "hello from mars"
	var pat = "mars"
	i := RabinKarp(txt, pat, 256, 103)
	assert.Equal(t, i, 11)
}
