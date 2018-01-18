package kv

import (
	"context"

	"github.com/pkg/errors"
)

const (
	BackendConsul = iota
	BackendEtcd
)

type Config struct {
	Backend int
	Addr    string
}
type Pair struct {
	Key   string
	Value []byte
}

type Pairs []*Pair

type WriteOption struct {
	Context context.Context
}
type QueryOption struct {
	Context context.Context
}

type IClient interface {
	Get(key string, option *QueryOption) (*Pair, error)
	GetList(prefix string, option *QueryOption) (Pairs, error)

	Put(pair *Pair, option *WriteOption) error
	Delete(key string, options *WriteOption) error

	Watch(key string, option *QueryOption) (chan *Pair, error)
	WatchList(prefix string, option *QueryOption) (chan Pairs, error)
}

func New(conf *Config) (c IClient, err error) {
	switch conf.Backend {
	case BackendConsul:
		c, err = newConsul(conf)
	default:
		err = errors.New("backend not supported")
	}

	return c, nil
}
