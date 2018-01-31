package pool

import (
	"net"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrGetEmpty      = errors.New("no connection free")
	ErrPoolClosed    = errors.New("pool closed")
	ErrNoConnCreated = errors.New("no connection created")
)

type ConnectionPool interface {
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

func NewConnectionPool(length, capacity int, d Dialer) (ConnectionPool, error) {
	p := &ChannelConnectionPool{
		conn:     make(chan net.Conn, capacity),
		capacity: capacity,
		d:        d,
	}

	// initial all connections
	for i := 0; i < length; i++ {
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
	case c := <-p.conn:
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
	}

	return nil
}

func (p *ChannelConnectionPool) Len() int {
	p.mux.Lock()
	defer p.mux.Unlock()
	return len(p.conn)
}

func (p *ChannelConnectionPool) Close() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	for conn := range p.conn {
		conn.Close()
	}

	close(p.conn)

	return nil
}
