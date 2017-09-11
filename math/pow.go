package math

/*
 * right-to-left binary exponentiation
 */
func Pow(n, m uint32) uint32 {
	var res, sq uint32 = 1, n
	for i := m; i>0; i>>=1 {
		if i & 1 != 0 {
			res *= sq
		}

		sq *= sq
	}
	return res
}