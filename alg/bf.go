package alg

// 检查t中是否包含子串p
// 返回出现子串的位置，如不存在则返回-1
func BruteForce(t, p string) int {
	if len(t) < len(p) {
		return -1
	}

	for i := 0; i <= len(t) - len(p); i++ {
		// 检查t从i开始len(p)的是否与p相等
		if t[i:i+len(p)] == p {
			return i
		}

		// 更为原始的做法
		/*
		j := 0
		for ;j<len(p);j++ {
			if t[i+j] != p[j] {
				break
			}
		}
		if j == len(p) {
			return i
		}
		 */
	}

	return -1
}

func BruteForceOld(t, p string) int {
	if len(t) < len(p) {
		return -1
	}

	for i:=0; i <= len(t) - len(p); i++ {

	}

	return -1
}

