package alg

/*
 * right-to-left binary exponentiation
 */
func Pow(m, n int) int {
	var res = 1
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			res *= m
		}

		m *= m
	}

	return res
}
