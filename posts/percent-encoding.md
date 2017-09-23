---
title: Percent Encoding
date: 2017-09-20 10:22:44
tags:
	- encode
	- decode
---

‘百分号编码'是一种将URI编码的机制, 又叫URL编码。

URI里允许的字符分为‘保留’和‘非保留’两种。‘保留’字符具有特殊的意义，比如'/'用于分割URL. ‘非保留‘字符就没有这种意义。用‘百分号编码‘,‘保留'字符会被转化为特殊的字符串
<!--more-->
保留字符：

```
! * ' ( ) ; : @ & - + $ , / ? # [ ] =  
```
当‘保留’字符的使用目的并不是原意，比如 ‘&’ 不是用来分割参数，需要进行 ‘百分号编码’。方法是 字符的ASCII编码的十六进制字符前加一个 ‘%’。 ‘%’字符必须进行百分比编码。列表如下:

```
!	#	$	&	'	(	)	*	+	,	/	:	;	=	?	@	[	]
%21	%23	%24	%26	%27	%28	%29	%2A	%2B	%2C	%2F	%3A	%3B	%3D	%3F	%40	%5B	%5D

```

非保留字符, 不需要进行编码

```
A-Z a-z 0-9 - _  . ~
```


二进制数或其他字符数据需要先转换成utf8字节序列，然后最字节进行百分比编码

```go
package main

import (
	"net/url"
	"log"
	"fmt"
)

func main() {
	var query = "a=1&b=2"
	es := url.QueryEscape(query)
	fmt.Printf("%s encoded to %s\n", query, es)

	us, err := url.QueryUnescape(es)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s decoded to %s\n", es, us)
}

localhost:run ben$ go run main.go 
a=1&b=2 encoded to a%3D1%26b%3D2
a%3D1%26b%3D2 decoded to a=1&b=2

```

#### 参考

https://zh.wikipedia.org/wiki/百分号编码