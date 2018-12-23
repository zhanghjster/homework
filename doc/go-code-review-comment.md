---
title: Go Code Review Comment
date: 2018-09-24 12:52:30
tags:
---

#### Gofmt

使用gofmt自动修复大多数代码排版问题，使用goimports增加或移除import

#### Commet Sentences

注释语句应该是完整的语句，即便看上去是冗余的。这会让godoc提取的文档具有良好的格式。注释应该以要描述的事物的名称开头并以点号结束：

~~~go
// Request represents a request to run a command.
type Request struct { ...

// Encode writes the JSON encoding of req to w.
func Encode(w io.Writer, req *Request) { ...
~~~

#### Contexts

context.Context类型的值包含跨API和进程的安全凭证(security credentials)、追踪信息(tracing information)、截止时间(deadline)和取消信号(cancellation signal)。Context在Go程序从RPC或HTTP或传入、传出，在整个函数的调用链显示的传递。

大多数函数使用第一个参数传递context:

~~~go
func F(ctx context.Context, /* other arguments */) {}
~~~

A function that is never request-specific may use context.Background(), but err on the side of passing a Context even if you think you don't need to. The default case is to pass a Context; only use context.Background() directly if you have a good reason why the alternative is a mistake.

对于结构类型，不要包含Context成员，而是将ctx参数添加到需要传递它的该类型的每个方法上，而那些必须与标准库或第三方库中的接口匹配的方法除外。

Don't create custom Context types or use interfaces other than Context in function signatures.

如果有要传递的应用程序数据，请将其放在参数(parameter)、接收器(receiver)、全局变量(global)或者context中。

Context是不可更改的，因此可以将其传递给多个调用来共享相同的安全凭证(security credentials)、追踪信息(tracing information)、截止时间(deadline)和取消信号(cancellation signal)

#### Declaring Empty Slices

当声明一个空的Slice，首选:

~~~go
var t []int
~~~

而不是：

~~~go
var t = []string{}
~~~

前一个方法声明了一个空的slice，而后者则声明非空但长度为0的slice(分配了内存空间)。他们在功能上相同（$len$和$cap$都为0)，但应首选前者。

在有限的一些情况下，后者为首选。例如，编码JSON对象时，nil Slice编码为null而[]string{}则编码为[]。

在设计接口时，避免区分nil slice和non-nil,zero-length Slice。因为这会导致细微的编程错误。

更多关于nil请点[这里](https://www.youtube.com/watch?v=ynoY2xz-F8s)

[这里](https://stackoverflow.com/questions/44305170/nil-slices-vs-non-nil-slices-vs-empty-slices-in-go-language) 也有nil Slice和空slice的区别的介绍

#### Crypto Rand

不要使用math/rand生成秘钥，即便是一次性的。没有种子的话，生成器是完全可预测的。用time.Nanoseconds()生成种子也只有几位的墒。可以只用crypto/rand的Reader，如果需要文本，使用base64或打印成Hex

 ~~~go
import (
    "crypto/rand"
    // "encoding/base64"
    // "encoding/hex"
    "fmt"
)

func Key() string {
    buf := make([]byte, 16)
    _, err := rand.Read(buf)
    if err != nil {
        panic(err)  // out of randomness, should never happen
    }
    return fmt.Sprintf("%x", buf)
    // or hex.EncodeToString(buf)
    // or base64.StdEncoding.EncodeToString(buf)
}
 ~~~

#### Doc Commet

所有顶级的导出名称都需要有注释，不常用的未导出的类型或函数声明也需要有注释。[https://golang.org/doc/effective_go.html#commentary](https://golang.org/doc/effective_go.html#commentary)有更多关于注释的约定。

#### Don't Panic

不要使用panic处理普通的error，使用error或多返回值。[https://golang.org/doc/effective_go.html#errors](https://golang.org/doc/effective_go.html#errors)有更多说明

#### Error Strings

Error串不要首字母大写，不要有句号。因为error很多时候会被拼接打印。比如使用fmt.Errorf("something bad")而不是fmt.Errorf("Something bad")，这样log.Printf("Reading %s: %v", filename, err)就不会在打印结果的字符串中出现大写字母和句号。

#### Examples

添加一个包的时候，对常用的用法增加可运行的实例或测试来演示如何使用。

[这里](https://blog.golang.org/examples)有更多说明

#### Goroutine Lifetimes

当产生goroutine时，要清楚什么时间或什么条件下它会退出

如果被管道的读取或发送阻塞，goroutine可能发生泄漏，即便管道无法访问GC不会终止goroutine。

即便goroutine不泄露，当它不再有用还任由它运行还是会导致其他的难以诊断的问题。向关闭的管道写入会触发panic。在"不需要结果"之后修改仍在使用的输入会导致数据竞争。任由gorouting长时间运行会导致不可预测的内存使用。

尽量保持并发代码足够简单来确保goroutine的生命周期简洁清晰。如果不行就加入注释说明goroutine退出的时机。

#### Handle Errors





















