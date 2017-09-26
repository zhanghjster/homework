package alg

import (
	"testing"
)

func TestBruteForce(t *testing.T) {
	dst := "I have a good "
	src := "have"

	i := BruteForce(dst, src)
	if i != 2 {
		t.Errorf("%q should contain %q", dst, src)

	}

	dst = "i Have a good"
	i = BruteForce(dst, src)
	if i == 2 {
		t.Errorf("%q should not container %q", dst, src)
	}
}
