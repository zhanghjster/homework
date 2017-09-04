package math

import "strconv"

func RunLengthEncoding(s string) (res []byte, err error) {
	if len(s) == 0 {
		return
	}

	// 当前字符的个数
	var n = 1
	for _, c := range s {
		if c != int32(cur) {
			res = append(res, []byte(strconv.Itoa(n))...)
			cur, n = c, n+1
		} else {
			n++
		}
	}

	return res
}
