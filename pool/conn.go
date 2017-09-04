package pool

import (
	"net"
	"sync"
	"github.com/pkg/errors"
)

var (
	ErrNoConnCreated = errors.New("no connection created")
	ErrPut = errors.New("fail put the connection back")
	ErrGetEmpty = errors.New("no connection free")
	ErrPoolClosed = errors.New("pool closed")
)

type Connection interface {
	Get() (net.Conn, error)
	Put(conn net.Conn) error
	Len() int
	Close() error
}

type Dialer func() (net.Conn, error)

type ChannelConnectionPool struct {
	conn     chan net.Conn
	capacity int
	length   int
	mux      sync.Mutex
	d        Dialer
}

func NewConnectionPool(length, capacity int, d Dialer) (Connection, error) {
	p := &ChannelConnectionPool{
		conn:     make(chan net.Conn, capacity),
		capacity: capacity,
		d:        d,
	}

	// initial all connections
	for i:=0; i < length; i++ {
		if c, err := d(); err == nil {
			p.length++
			p.conn <- c
		}
	}

	if len(p.conn) == 0 {
		return nil, ErrNoConnCreated
	}

	return p, nil
}

func (p *ChannelConnectionPool) Get() (net.Conn, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.conn == nil {
		return nil, ErrPoolClosed
	}

	select {
	case c := <- p.conn:
		return c, nil
	default:
		// create new one
		if p.length < p.capacity {
			if c, err := p.d(); err == nil {
				p.length++
				return c, nil
			}
		}

		return nil, ErrGetEmpty
	}
}

func (p *ChannelConnectionPool) Put(conn net.Conn) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.conn == nil {
		return ErrPoolClosed
	}

	select {
	case p.conn <- conn:
	default:
		return ErrPut
	}

	return nil
}

func (c *ChannelConnectionPool) Len() int{
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.conn)
}

func (c *ChannelConnectionPool) Close() error{
	c.mux.Lock()
	defer c.mux.Unlock()

	for conn := range c.conn {
		conn.Close()
	}

	close(c.conn)

	return nil
}
