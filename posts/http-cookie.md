---
title: HTTP Cookie
date: 2017-10-11 21:29:58
tags:
    - cookie
    - 无状态协议
categories: 基础
---

HTTP是一种无状态的协议，服务器与客户端之间的每次交互式是独立的事务，它们不保存前后事务之间的关系。比如，服务器多次处理客户端页面的请求，它不记录也不关心这些请求是否来自同一客户端或请求的是否是相同页面。与无状态协议相反的有状态协议则是事务之间是有联系的，交互双方需保存信息来记录这种联系，比如TCP协议，它要求通信双方在每次传输后都要知道从对方接受了多少数据、是否收到了最后一个数据包，也就是它要求保存前后事务产生的状态变化

HTTP的无状态特性在交互式WEB应用出现后遇到了问题，因为这些应用往往要求保存请求之间的状态，比如，用户在一个页面上登录后，之后所有的请求都要求是登录状态。于是用于HTTP协议保持状态的技术产生了，一个是用于服务器端的Session，另一个则是用户客户端的Cookie

<!-- more -->

Cookie是客户端保存服务器返回信息的小型纯文本，客户端按一定规则保存这些信息后，在接下来的请求中会将其返回给服务器，服务器利用这些信息来识别客户端的状态。比如，将用户登录的状态保存到cookie里，服务器就可以根据它来判定客户端的请求是否是登录状态

#### Cookie创建

服务器通过在返回给客户端的HTTP header里的**Set-Cookie**字段来告诉客户端保存cookie信息，它的格式如下(中括号里为可选)

```
Set-Cookie: NAME=VALUE[;expires=DATE][; path=PATH][; domain=DOMAIN_NAME][; secure]
```

##### ***NAME=VALUE***

它由分号、逗号、空格之外的字符组成，表示cookie的名称和值。它为唯一的必选项，在之后的所有客户端发送给服务器的请求中都要包含到请求的Header中，而可选项则不会再发送给服务器

cookie的值一般情况下推荐要进行URL编码，但不是必须的，不过几乎所有实现都对cookie进行了编码。通常对***NAME*** 和 ***VALUE***分别进行编码，而 ***=*** 不编码

由于cookie存在着数量的限制，开发者门使用subcookies的办法来增加cookie的存储量，方法就是在***VALUE***中存储一些 ***name-value***对，如下格式

```
NAME1=foo=bar&foo1=bar1&foo2=bar2
```

不过这种格式需要开发者自定义解析方式

##### ***expires=DATE***

通过这个选项设置cookie的过期时间，一旦到达了这个时间，cookie就不在有效，应该删除，格式如下

```
Wdy, DD-Mon-YYYY HH:MM:SS GMT
```

格式基于 [RFC 822](https://curl.haxx.se/rfc/rfc0822.txt), [RFC 850](https://curl.haxx.se/rfc/rfc0850.txt), [RFC 1036](https://curl.haxx.se/rfc/rfc1036.txt), 和 [RFC 1123](https://curl.haxx.se/rfc/rfc1123.txt) ,但做了一定修改，日期之间的空格换成了横线，时区必须为GMT

它为可选项，当不设置时，用户的session结束(关闭浏览器)时候就失效

##### ***domain=DOMAIN_NAME***

当客户端需要选择发送给服务器的cookie时，通过将这个字段与要访问的URL的域名做尾部匹配来做判断。尾部匹配是指 ***DOMAIN_NAME***与URL的全域名尾部相同，比如 ***domain=benx.io***的***benx.io***与 "***www.benx.io***" "***wap.benx.io***"都匹配，那么在访问这些域名的URL时也有可能将cookie发送给服务器，说‘可能’是因为还需要进行***PATH***的匹配，两个都通过才会发送

如果不设置，默认被设置成返回这个cookie的服务器的域名

##### ***path=PATH***

用于指定一个于cookie有效的path。如果cookie已经通过了***domain***的匹配，则进行path部分的匹配，采用的办法是前缀匹配，如果访问url里的path开头包含***PATH***, 则通过。比如 '/blog' 与 ‘/blog/archive'匹配

如果不设置，默认被设置成header里对应URL的path部分

##### ***secure***

它只是一个标记，没有值，当请求通过SSL或HTTPS创建时，cookie才会被发送到服务器，默认情况下，HTTPS链接上传输的cookie都自动有这个选项

##### ***HttpOnly***

这个选项是微软在IE6 SP1里引入的，而原始的Cookie说明文档里并没有，后来被更多的浏览器所接受。这个选项的目的是控制 ***document.cookie***的访问，防止javasciprt发起跨站攻击窃取cookie。同***secure***一样，它也只是一个标记，被设置后，***document.cookie***将不能防问该cookie

##### ***Max-Age***

同expires一样用于设置cookie的过期时间，不同的是的它表示的是从cookie被设接收到过期那一刻的时间间隔。开始到过期

为可选项

#### Cookie发送

当客户端想服务器请求一个URL，需要对cookie进行***domain*** ***path***的匹配，符合条件的cookie的 ***name=value***对会被包含进HTTP请求的中，格式如下：

```
Cookie: NAME1=OPAQUE_STRING1; NAME2=OPAQUE_STRING2 ...
```

发送的内容为set-cookie时的原始内容，不需要做任何处理

#### Cookie更新

一个cookie的唯一性由 ***NAME-DOMAIN-PATH-secure***共同决定，要想修改一个cookie，必须发送另一个这四个值相同的cookie，比如

如果原cookie是通过如下设置的

```
Set-Cookie: NAME1=benx; domain=benx.io; path=/
```

则需要用如下方法来更新

```
Set-Cookie: NAME1=lee; domain=benx.io; path=/
```

由于过期时间是可选项，所以如果更新cookie时不带过期时间，则不会更新过期时间

#### 总结

cookie作为一种简单有效的手段弥补了HTTP协议的一些不足，虽然历经很多年的发展，依然在广泛应用。在日常IT工作中会经常与之打交道，还是需要对他进行深入理解

参考：

https://curl.haxx.se/rfc/cookie_spec.html

https://www.nczonline.net/blog/2009/05/05/http-cookies-explained/

https://stackoverflow.com/questions/19899236/is-tcp-protocol-stateless-or-not

https://en.wikipedia.org/wiki/Stateless_protocol













