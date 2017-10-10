package alg

import (
	"fmt"
	"strings"
)

// d 为基数, q 为模,
// d为文本t的所有字符种类的和, q为满足d*q接近计算机字长的最大的素数
//
// 实际应用中d如果选择合适则不需要进行hash运算，直接利用无符号整数溢出来实现取模
// 这样既可以发挥现代计算机乘法效率高于除法优势，还可以减少除法的不必要运算
// Golang的strings库就采用了这个办法，d的选择是16777619(FNV算法也用到)
func RabinKarp(t, p string, d, q int) int {
	var n, m int = len(t), len(p)
	if n < m || m == 0 {
		return -1
	}

	// 计算 h = (d%q)^m
	// right to left binary exponentiation
	var h, s int = 1, d
	for i:=m; i>0;i>>=1 {
		if i&1 != 0 {
			h = (h*s)%q
		}
		s = (s*s)%q
	}

	// 计算t[0..m]和p的hash
	var ht, hp int
	for i:=0; i<m; i++ {
		ht = (d*ht + int(t[i]))%q
		hp = (d*hp + int(p[i]))%q
	}

	// 滚动检查
	for i:=0; i <= n-m; i++  {
		// 找到潜在匹配后做字符串验证
		fmt.Println(t[i:i+m])
		fmt.Println("hp => ", hp, " ht => ", ht)
		showp(t[i:i+m], d, q)
		if ht == hp && t[i:i+m] == p {
			return i
		}

		// 根据ht[i]计算ht[i+1]
		ht = (d*ht  - int(t[i])*h + int(t[i+m]))%q

		// ht为负数时转化成与之同余的正数以便于与p进行比较
		// 如果ht = (dt*ht + q - (int(t[i])*h)%q + int(t[i+m])%q)
		// 就不需要下面转化，因为ht肯定为正数，但它效率不高
		fmt.Println(ht)
		if ht < 0 {
			ht += q
		}
	}

	return -1
}

func showp(s string, d, q int) {
	var a = make([]string, len(s))
	var b = make([]string, len(s))
	var t int
	for i:=0; i<len(s);i++ {
		a[i] = fmt.Sprintf("%c*%d^%d",s[i],d,len(s)-i-1)
		b[i] = fmt.Sprintf("%d*%d^%d",byte(s[i]),d,len(s)-i-1)
		t = d*t + int(s[i])
	}

	fmt.Println(strings.Join(a, " + "))
	fmt.Println(strings.Join(b, " + "))
	fmt.Println(t%q)
}