---
title: Bytes Buffer 
date: 2017-10-19 22:09:32
tags:
  - buffer
  - bytes 
  - 注释
categories: go
---

bytes.Buffer是带有读写方法的变长字节缓冲区，零值的Buffer是一个开箱即用的空缓存，到目前用了多次，所以是时候将部分注释拎出来做个总结，强化一下记忆与理解

```go
var buf bytes.buffer
```

<!-- more -->

#### func NewBuffer

```go
func NewBuffer(buf []byte) *Buffer
```

使用 buf 作为初始内容创建一个Buffer，新的buffer接管 buf，调用者不应再使用buf。通过NewBuffer可以创建一个缓冲区来读取现有内容(buff)。或者初始化写入的大小，但buff应该用所需的capacity和为0的length来初始化

通常情况下使用new(Buffer)或者直接声明Buffer变量就可以高效的初始化一个Buffer

#### func NewBufferString

```go
func NewStringBuffer(s string) *Buffer	
```

使用s作为初始内容生成一个Buffer，可以用来创建一个用于读取一个字符串的Buffer

#### func (*Buffer) Bytes()

```go
func (b *Buffer) Bytes() []byte
```



返回长度为b.Len()的slice，为buffer里还没有被读出的内容。slice会一直有效到再次对buffer进行 Write、Read、Reset、Truncate等更新操作，所以更新slice会影响到将来的read操作

#### func (*Buffer) Next()

```go
func (b *Buffer) Next(n int) []byte				
```

返回 包含n个byte的slice，如果buffer里有不够n个未读数据则返回全部内容。slice只有效到Buffer的下一次读写操作

#### func (*Buffer) Read

```go
func (b *Buffer) Read(p []byte) (n int, err error)
```

读取len(p)个字节数据直到读完，返回值n表示读出的字节数，如果已经没有数据可读(n=0)，返回err为io.EOF

```go
buf := bytes.NewBufferString("abcd")
var data = make([]byte, 3)
fmt.Println(buf.Read(data))
fmt.Println(buf.Read(data))
fmt.Println(buf.Read(data))
```

结果为

```
3 <nil>
1 <nil>
0 EOF
```

#### func (*Buffer) ReadByte() 

```go
func (b *Buffer) ReadByte() (byte, error)
```

读取buffer里的下一个字节， 如果没有则返回 io.EOF

```go
buf := bytes.NewBufferString("abcd")
var data = make([]byte, 3)
fmt.Println(buf.Read(data))
fmt.Println(buf.ReadByte())
fmt.Println(buf.ReadByte())
```

结果为

```
3 <nil>
100 <nil>
0 EOF
```

#### func (*Buffer) ReadBytes

```go
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) 
```

读取buffer里一直到第一次出现delim的所有字节，返回包含delim在内的slice。如果在找打delimiter之前遇到error，则返回已经读取的数据和 error(通常是io.EOF)。当且仅当line不以delim为结尾时才 err != nil

```go
buf := bytes.NewBufferString("abcd")
line, err := buf.ReadBytes('c')
fmt.Println(string(line), err)

line, err = buf.ReadBytes('e')
fmt.Println(string(line), err)
```

结果为

```
abc <nil>
d EOF
```

#### func (*Buffer) ReadRune

```go
func (b *Buffer) ReadRune() (r rune, size int, err error)
```

读取并返回下一个UTF8编码的Code Point。如果没有有效的字节，返回io.EOF。如果是错误的UTF8编码，则只读取一个字节，并且返回 'U+FFFD', 1

```go
var s = "大"
fmt.Printf("%+q %x %b\n", s, s, []byte(s))
buf := bytes.NewBufferString(s)

// 读取一个字符破坏编码
buf.ReadByte()

// 无效的编码，读取一个字节
// 返回 "U+FFFD", 1, nil
r, n, err := buf.ReadRune()
fmt.Printf("%+q %d, %v\n", r, n, err)

// 无效的编码吗，读取一个字节
// 返回 "U+FFFD", 1, nil
r, n, err = buf.ReadRune()
fmt.Printf("%+q %d, %v\n", r, n, err)

// 已无数据
r, n, err = buf.ReadRune()
fmt.Printf("%+q %d, %v\n", r, n, err)
```

返回

```
"\u5927" e5a4a7 [11100101 10100100 10100111]
'\ufffd' 1, <nil>
'\ufffd' 1, <nil>
'\x00' 0, EOF
```

#### func (*Buffer) ReadString

```go
func (b *Buffer) ReadString(delim byte) (line string, err error)
```

读取一直到dlim的字节并返回包含delim在内的字符串，如果在遇到delim之前遇到错误，则返回已经读取的内容和error(通常为io.EOF)。当且仅当err!=nil时，返回的字符串不包含delim

#### func ReadFrom

```go
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error)
```

从ioreader中读取数据并追加到缓存里直到遇到read的io.EOF，buffer的长度会自动按需增加。返回值 n 为成功读取的字节数，在读的过程中遇到的除io.EOF的error外也返回。如果buffer太大，则返回 ErrTwoLarge

#### func (*Buffer) Cap() int

返回底层用于保存数据的slice的容量

#### func (*Buffer) Len() int

返回还未读出的bytes数量 b.Len == len(b.Bytes())

#### func (*Buffer) Grow()

```go
func (b *Buffer) Grow(n int)	
```

增加Buffer的容量n使其能够至少再容纳n个bytes， 如果n为负值 Grow会panic，如果buffer不能grow，则panic ErrTooLarge

#### func (*Buffer) Reset

```go
func (b *Buffer) Reset()
```

重置buffer为空，但保留底层的存储空间，和Trancate(0)相同

#### func (*Buffer) Truncate

```go
func (b *Buffer) Truncate(n int)
```

丢弃前n个未读字节后的所有字节，但继续使用相同的底层存储空间，如果n为负数或者大于buffer的长度则panic

#### func (*Buffer) String()

```go
func (b *Buffer) String() string
```

以字符串形式返回未读内容，但不会将这些内容置为已读，如果buffer是nil，则返回nil

#### func（*Buffer) UnreadByte

```go
func (b *Buffer) UnreadByte() error
```

将最新一次成功读取的字节设置为未读。如果上一次读之后有写的操作或者上一次读返回error或read返回了0字节，则返回error(***似乎有问题，上一次read如果返回io.EOF并没有让unreadbyte返回error***）

```go
buf := bytes.NewBufferString("abcd")
fmt.Println(buf.String())

buf.ReadByte()
fmt.Println(buf.UnreadByte())

fmt.Println(buf.ReadString('e'))

// 此处读取返回了 io.EOF
fmt.Println(buf.ReadString('e'))
// 此处并没有返回error
fmt.Println(buf.UnreadByte())

fmt.Println(buf.String())
```

结果

```
abcd
<nil>
abcd EOF
 EOF
<nil>
d
```

#### func UnreadRune()

```go
func (b *Buffer) UnreadRune() error
```

将上一次读出的rune置位未读，如果上一次读取返回error，则UnreadRune也返回error

#### func (*Buffer) Write

```go
func (b *Buffer) Write(p []byte) (n int, err error)
```

将p里内容写入缓冲并返回成功p的长度，err永远为nil，如果buffer太长则返回 ErrTooLarge

#### func (*Buffer) WriteByte

```go
func (b *Buffer) WriteByte(c byte) error
```

将字节c写入缓冲，永远返回nil的error，如果buffer太大，则返回ErrTooLarge

#### func (*Buffer) WriteRune

```go
func (b *Buffer) WriteRune(r rune) (n int, err error)
```

将UTF8编码的Unicode point写入buffer，返回字符的长度和永远为nil的error，如果buffer太长则返回ErrTooLarge

#### func (*Buffer) WriteString

```go
func (b *Buffer) WriteString(s string) (n int, err error)
```

将字符串写入到buffer，err永远为空

#### func (*Buffer) WriteTo

```go
func (b *Buffer) WriteTo(w io.Wirter) (n int64, err error)
```

将buffe中内容写入到writer，返回写入的字节数和写入过程中遇到的任何错误

#### 总结

1. 读取方法

   ```go
   func (b *Buffer) Read(p []byte) (n int, err error)
   func (b *Buffer) ReadByte() (byte, error)
   func (b *Buffer) ReadBytes(delim byte) ([]byte, error) 
   func (b *Buffer) ReadString(delim byte) (line string, err error)
   func (b *Buffer) ReadRune() (r rune, err error)
   func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error)
   func (b *Buffer) Next(n int) []byte
   func (b *Buffer) Bytes() []byte
   func (b *Buffer) String() string

   ```

2. 写入方法

   ```go
   func (b *Buffer) Write(p []byte) (n int, err error)
   func (b *Buffer) WriteByte(c byte) error
   func (b *Buffer) WriteString(s string) (n int, err error)
   func (b *Buffer) WriteRune(r rune) (n int, err error)
   func (b *Buffer) WriteTo(w io.Wirter) (n int64, err error)
   ```

3. 其他

   ```
   func (b *Buffer) Cap() int
   func (b *Buffer) Len() int
   func (b *Buffer) Trancate(n int)
   func (b *Buffer) Grow(n int)	
   func (b *Buffer) Reset()
   func (b *Buffer) UnReadByte() error
   func (b *Buffer) UnReadRune() error
   ```

   ​

   ​
