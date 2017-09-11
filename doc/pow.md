#### 

在看Golang的strings库的strings.Index(s, sep string)函数时时发现了一个计算n的m次幂的方法。代码如下：

````
func hashStr(sep string) (uint32, uint32) {
	....
	var pow, sq uint32 = 1, primeRK
	for i := len(sep); i > 0; i >>= 1 {
		if i&1 != 0 {
			pow *= sq
		}
		sq *= sq
	}
	....
}
````
发现它并不是简单的m次n的乘积得到，其中必有一些技巧的存在。

#### left-to-right Exponentiation

这是一个很古老的办法，具体逻辑是这样，比如：要求解x^25:

   1. 首先将25用二进制表示 11001.
   2. 去掉最高位的1， 变为 1001.
   3. 用D(double)代表0， D1（double and +1）代表1。得到计算的指令 D1DDD1
   4. 从左到右，按照上面得到的指令进行计算。顺序得到x^2,x^3,x^6,x^12,x^24,x25. 共执行了6次。

这个办法存在一个问题就是需要先计算出从左到右的指令集，并需要一定的空间保存。程序上实现有一定的复杂度。 这个办法也叫 Addition Chain Exponetiation


#### right-to-left Exponetiation

这个算法也比较古老，实现方法如下,

````
	Power(x,n) {
	  r := 1
	  y := x
	  while (n > 1) { 
	    if odd(n) then r := r*y
	    n := floor(n/2)
	    y := y*y
	  }
	  r := r*y
	  return(r)
	}
````

开头的golang代码就是这个算法的实现。与left-to-right比较实现起来很容易，不需要额外的空间存储指令。


#### 参考：

https://en.wikipedia.org/wiki/Addition-chain_exponentiation

http://primes.utm.edu/glossary/xpage/BinaryExponentiation.html

https://www.quora.com/What-are-some-fast-algorithms-for-computing-the-nth-power-of-a-number