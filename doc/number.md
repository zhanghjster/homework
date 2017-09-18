### 补码法

补码法是在数学和计算中通过加一个正数来实现减法的技术，在机械计算机中很常用，在现代计算机中仍广泛使用

#### 基数补码

一个以b为基数的N位数y的基数补码用  <math xmlns="http://www.w3.org/1998/Math/MathML"><msup><mi>b</mi><mi>N</mi></msup><mo>-</mo><mi>y</mi><mspace linebreak="newline"/></math> 定义，可以由基数反码加一直接得到。

#### 基数反码

一个以b为基数的N位数y的基数反码是它的基数补码减一 , <math xmlns="http://www.w3.org/1998/Math/MathML"><msup><mi>b</mi><mi>N</mi></msup><mo>-</mo><mn>1</mn><mo>-</mo><mi>y</mi><mspace linebreak="newline"/></math>。因为 <math xmlns="http://www.w3.org/1998/Math/MathML"><msup><mi>b</mi><mi>N</mi></msup><mo>-</mo><mn>1</mn><mo>=</mo><mo>(</mo><mi>b</mi><mo>-</mo><mn>1</mn><mo>)</mo><msup><mi>b</mi><mrow><mi>n</mi><mo>-</mo><mn>1</mn></mrow></msup><mo>+</mo><mo>.</mo><mo>.</mo><mo>.</mo><mo>+</mo><mo>(</mo><mi>b</mi><mo>-</mo><mn>1</mn><mo>)</mo><mspace linebreak="newline"/></math> ，所以基数反码是y的每一位用基数-1减去每一位后的值替换后得到的值。 反过来，一个数的基数补码是其基数反码+1.

 比如， 在十进制下三位数 678 的基数反码是 321 = 999 - 678

x-y可以用 x 加 y 的基数补码得到，即 <math xmlns="http://www.w3.org/1998/Math/MathML"><mi>x</mi><mo>+</mo><msup><mi>b</mi><mi>N</mi></msup><mo>-</mo><mi>y</mi><mspace linebreak="newline"/></math> 

为什么用 x 加 y的补码结果和 x - y 结果相同呢，因为计算器的字长是有限的，比如为N位，当计算结果超出N位时，最高位将被抛弃。 实际的计算结果相当于取 <math xmlns="http://www.w3.org/1998/Math/MathML"><msup><mi>b</mi><mi>N</mi></msup><mspace linebreak="newline"/></math> 的模。

因为：

* <math><mo>(</mo><mi>x</mi><mo>-</mo><mi>y</mi><mo>)</mo><mo>&#x2261;</mo><mo>(</mo><mi>x</mi><mo>-</mo><mi>y</mi><mo>)</mo><mo>&#xA0;</mo><mi>m</mi><mi>o</mi><mi>d</mi><mo>&#xA0;</mo><msup><mi>b</mi><mi>N</mi></msup></math>

* <math xmlns="http://www.w3.org/1998/Math/MathML"><mn>0</mn><mo>&#x2261;</mo><mo>&#xA0;</mo><msup><mi>b</mi><mrow><mi>N</mi><mo>&#xA0;</mo></mrow></msup><mi>m</mi><mi>o</mi><mi>d</mi><mo>&#xA0;</mo><msup><mi>b</mi><mi>N</mi></msup></math>

所以：

<math xmlns="http://www.w3.org/1998/Math/MathML"><mo>(</mo><mi>x</mi><mo>-</mo><mi>y</mi><mo>)</mo><mo>&#x2261;</mo><mo>(</mo><mi>x</mi><mo>+</mo><msup><mi>b</mi><mi>N</mi></msup><mo>-</mo><mi>y</mi><mo>)</mo><mo>&#xA0;</mo><mi>m</mi><mi>o</mi><mi>d</mi><mo>&#xA0;</mo><msup><mi>b</mi><mi>N</mi></msup><mo>&#x2261;</mo><mo>(</mo><mi>x</mi><mo>+</mo><mi>y</mi><mi>&#x7684;&#x57FA;&#x6570;&#x8865;&#x7801;</mi><mo>)</mo><mi>m</mi><mi>o</mi><mi>d</mi><mo>&#xA0;</mo><msup><mi>b</mi><mi>N</mi></msup><mo>=</mo><mo>(</mo><mi>x</mi><mo>+</mo><mi>y</mi><mi>&#x7684;&#x57FA;&#x6570;&#x53CD;&#x7801;</mi><mo>+</mo><mn>1</mn><mo>)</mo><mi>m</mi><mi>o</mi><mi>d</mi><mo>&#xA0;</mo><msup><mi>b</mi><mi>N</mi></msup><mspace linebreak="newline"/><mspace linebreak="newline"/></math>

所以，如果用y的基数补码表示y的负数，则减法就转化成加法，结果如果为负数，也是一个正数的基数补码

扩展到现代计算机中，则基数为2，基数反码对应one's complement， 基数补码对应two's complement

### 计算机中有符号整数的表示方法

在早期的计算机系统中，有很多关于如何实现硬件和数学计算的对立观点。其中一个就是如何表示有符号数，尤其是负数，不同的阵营都有着各自鲜明的观点。其中一个阵营支持‘补码‘(Two's Compelement)，另一个阵营则支持‘反码’(One's Complement)。如今，前者则是普遍接受与使用。更古老的一个方法则是‘源码’(Sign and Magnitude).


#### 原码（Signed and Magnitude)

这种办法使用一bit位来标示符号(通常是最高位), 符号位为1表示负数，0为正数，其余位用来表示数的绝对值。比如用八位的二进制源码来标示
‘+1’ 和 ‘-1’分别为

```
+1 = 0000 0001
-1 = 1000 0001
```

一个八位的原码可以表示的有符号整数范围： -127 - +127.

源码的表示方法存在着一些问题，其中之一是它不能直接从寄存器交给计算单元计算，比如 计算 1 + (-1), 如果直接将 000000001  10000001 交给计算单元进行加法，结果是 10000010 = -2, 显然是不对的，所以还得进行转换的操作(可以转成反码)，计算完后再将结果转换回来，这样就需要设计更多的电路来判断正负和进行转换，增加了系统复杂度也增加了分立式晶体管的成本。

另一个问题是，0 有两种表示方法，比如八位下有 ‘10000000’(-1) 和 ‘00000000’(+0), 这在判断数字是否为0时需要考虑两种情况，也增加了电子设备的设计复杂度。

#### 反码(One's Complement)

一个二进制数取‘反码’指将数的每个bit位取反所代表的数，在一些计算操作中用来将一个数取负。

用反码来表示有符号整数时，正数用其自身表示，负数用它的绝对值取反码表示。

比如 +1和-1用反码表示为：

````
	+1 = 0000 0001
	-1 = 1111 1110
````

Reference：

https://en.wikipedia.org/wiki/
Method_of_complements#Binary_example

https://en.wikipedia.org/wiki/Signed_number_representations








































	
	