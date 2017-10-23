Bufio包实现了缓冲IO，它通过为io.Reader或io.Writer包装一层缓冲实现新的Reader或Writer，为文本IO提供缓冲和一些帮助

#### type Reader

##### func NewReader(rd io.Reader) *Reader

创建一个缓冲reader，缓冲大小默认值4096字节

##### func NewReaderSize(rd io.Reader, size int) *Reader

创建一个最小为size字节的缓冲reader，如果rd一个是个bufio.Reader并且它的大小>=size,则直接返回rd

##### func (*Reader) Read

```go
func (b *Reader) Read(p []byte) (n int, err error)
```

读取数据到p,返回成功读取的长度。数据最多在底层reader读取一次，所以n可能小于len(p)。当遇到EOF时

，返回的长度n为0,err=io.EOF

##### func (*Reader) ReadByte

```go
func (b *Reader) ReadByte() (byte, error)
```

读取一个字节，如果没有则返回error

##### func (* Reader) ReadBytes

```go
func (b *Reader) ReadBytes(delim byte) ([]byte, error)
```

从输入中一直读取到delim, 返回的slice以delim为结尾，当返回error时，最后一个字节一定不是delim。如果是简单的用法，Scanner可能更方便





































