package utils

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidIP   = errors.New("invalid ip address")
	ErrInvalidMask = errors.New("invalid network mask")
)

// Convert ip4 string to int
func Ip4ToInt(ip string) (uint32, error) {
	if ip[0] == '.' || ip[len(ip)-1] == '.' {
		return 0, ErrInvalidIP
	}

	var m uint
	var seg, res uint32
	for i, c := range ip {
		if c >= '0' && c <= '9' {
			seg = 10*seg + uint32(c-'0')
		}

		if c == '.' || i == (len(ip) - 1) {
			res = seg<<(8*(3-m)) + res
			seg, m = 0, m+1
		}

		if c == '.' && ip[i-1] == '.' || seg > 0xFF {
			return 0, ErrInvalidIP
		}
	}

	if m != 4 {
		return 0, ErrInvalidIP
	}

	return res, nil
}

// convert mask string to int
func MaskToInt(m string) (uint32, error) {
	if m == "" || len(m) > 2 {
		return 0, ErrInvalidMask
	}

	var s int
	for _, c := range m {
		if c >= '0' && c <= '9' {
			s = 10*s + int(c-'0')
		} else {
			return 0, ErrInvalidMask
		}
	}

	if s < 1 || s > 32 {
		return 0, ErrInvalidMask
	}

	var mask uint32
	for i:=0; i < s; i++ {
		mask = mask | 1 << uint32(31 - i)
	}

	return mask, nil
}