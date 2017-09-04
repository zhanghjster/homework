package utils

import (
	"testing"
)

func TestIp4ToInt(t *testing.T) {
	var ip = "127.0.0.1"
	var ex = uint32(0x7F << 24 + 1)

	i, err := Ip4ToInt(ip)
	if err != nil {
		t.Error("nil error expected but got ", err)
	}

	if i != ex {
		t.Errorf("%b expected but got %b", ex, i)
	}

	for _, ip := range []string{
		"127.0..1", "127.0.256.1", ".127.0.0.1", "127.0.0.1.",
	} {
		if _, err := Ip4ToInt(ip); err != ErrInvalidIP {
			t.Error("invalid ip error expected")
		}
	}
}

func TestMaskToInt(t *testing.T) {
	var mask = "32"
	var ex = uint32(0xFFFFFFFF)
	m, err := MaskToInt(mask)
	if err != nil {
		t.Error("nil err expected but got", err)
	}

	if m != ex {
		t.Errorf("%b expected but get %b", ex, m)
	}

	for _, mask := range []string{
		"0", "33", "3m",
	} {
		if _, err := MaskToInt(mask); err != ErrInvalidMask {
			t.Error("invalid mask error expected")
		}
	}
}
