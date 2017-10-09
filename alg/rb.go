package alg

func RabinKarp(t, p string) int {
	var n, m int = len(t), len(p)
	if n < m || m == 0 {
		return -1
	}

	var d  = 256 		// 基数
	var q  = 7 	// 模

	// h 为d^(m-1)
	var h  = 1
	for i:=0; i<m; i++ {
		h = (h*d)%q
	}

	// 计算t[0..m]和p的hash
	var ht, hp int
	for i:=0; i<m; i++ {
		ht = (d*ht + int(t[i]))%q
		hp = (d*hp + int(p[i]))%q
	}
	// 做第一次检查
	if ht == hp && t[:m] == p {
		return 0
	}

	// 滚动检查
	for i:=1;i<=n-m;i++ {
		ht = (d*ht - int(t[i-1])*h + int(t[i+m-1]))%q
		if ht < 0 {
			ht += q
		}

		println("ht =>", ht, " d => ", d, " hp => ", hp)

		// 找到潜在匹配后做字符串验证
		if ht == hp && t[i:i+m] == p {
			return i
		}

	}

	// 对于串S[0..m], h = s[0]*d^(m-1)+s[1]*d^(m-2)+..+s[m-1]在d选择比较合适的素数的情况下，
	// 哈希结果分布会比较均匀, 不需要再进行模q的操作， 并且现代计算机乘法的效率要远高于除法
	// 所以在Go的strings包中的hash部分，没有模的计算，使用的d是16777619在FNV算法中也有用到

	return -1
}