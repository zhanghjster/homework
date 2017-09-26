package alg

func BruteForce(d, s string) int {
	if len(d) < len(s) {
		return -1
	}

	for i := 0; i <= len(d) - len(s); i++ {
		if d[i:i+len(s)] == s {
			return i
		}
	}

	return -1
}
