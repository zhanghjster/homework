package alg

func RotateArray(list []int, d int) {
	var n = len(list)

	// 假设gcd=(n,d)=g，则n=a*g,d=b*g
	// LCM(n,d) = a*b*g
	// 时间复杂度
	//  = gcd(n,d)*(LCM(n,d)/d)
	//  = g*(g*a*b/d)
	//  = (a*g)*(b*g)/d
	//  = a*g
	//  = n

	for i:=0; i < gcd(n, d); i++ {
		var tmp = list[i]
		j := i
		// LCM(n,d)/d
		for {
			k := j+d
			if k >= n {
				k -= n
			}
			if k == i {
				break
			}
			list[j] = list[k]
			j = k
		}
		list[j] = tmp
	}
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}