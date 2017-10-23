client.Do(req *Request)

  ->client.send(req *Request, deadline time.Time) 

​    ->send(ireq *Request, rt RoundTripper, deadline time.Time)

​      ->rt.RoundTrip(req)

​	->t.getConn(treq, cm)

​	    ->t.getIdleConn(cm)

​            ->t.dialConn(ctx, cm)

​            ->如果从idleconn里取得了，就把新的conn放到idleconn

​        ->pc.roundTrip(req)

bufio.NewWriter(wr *io.Writer）

建立在writer基础上的buffer,

conn ->包装成writer->包装成bufferWriter













func (*Transport) getConn

```go
func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (*persistConn, error)
```

1. pc, idleSince := t.getIdleConn(cm)， 取得IdleConn