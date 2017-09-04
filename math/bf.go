package math

func BruteForce(d, s string) int {
	if len(d) < len(s) {
		return -1
	}

	for i := 0; i <= len(d) - len(s); i++ {
		var j = 0
		for ;j<len(s); j++ {
			if d[i+j] != s[j] {
				break
			}
		}

		if j == len(s) {
			return i
		}
	}

	return -1
}
