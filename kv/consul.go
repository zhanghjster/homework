package kv

import (
	"context"

	"github.com/hashicorp/consul/api"
)

type Consul struct {
	conf *Config
	kv   *api.KV
}

func newConsul(conf *Config) (*Consul, error) {
	client, err := api.NewClient(&api.Config{
		Address: conf.Addr,
	})
	if err != nil {
		return nil, err
	}

	return &Consul{
		conf: conf,
		kv:   client.KV(),
	}, nil
}

func (c *Consul) Get(key string, opt *QueryOption) (*Pair, error) {
	p, _, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return c.convertPair(p), nil
}

func (c *Consul) GetList(prefix string, opt *QueryOption) (Pairs, error) {
	l, _, err := c.kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	return c.convertPairs(l), nil
}

func (c *Consul) Watch(key string, opt *QueryOption) (chan *Pair, error) {
	p, mt, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var ctx = context.Background()
	if opt != nil {
		ctx = opt.Context
	}

	var pairCh = make(chan *Pair)
	go func() {
		defer close(pairCh)

		for {
			pairCh <- c.convertPair(p)

			var qOpt = api.QueryOptions{
				WaitIndex: mt.LastIndex,
			}.WithContext(ctx)

			p, mt, err = c.kv.Get(key, qOpt)

			if c.done(ctx) {
				return
			}
		}
	}()

	return pairCh, nil
}

func (c *Consul) WatchList(prefix string, opt *QueryOption) (chan Pairs, error) {
	l, mt, err := c.kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	var ctx = context.Background()
	if opt != nil {
		ctx = opt.Context
	}

	var pairsCh = make(chan Pairs)
	go func() {
		defer close(pairsCh)
		for {
			pairsCh <- c.convertPairs(l)

			var qOpt = api.QueryOptions{
				WaitIndex: mt.LastIndex,
			}.WithContext(ctx)

			l, mt, _ = c.kv.List(prefix, qOpt)

			if c.done(ctx) {
				return
			}
		}
	}()

	return pairsCh, nil
}

func (c *Consul) Put(pair *Pair, opt *WriteOption) error {
	_, err := c.kv.Put(&api.KVPair{Key: pair.Key, Value: pair.Value}, nil)
	return err
}

func (c *Consul) Delete(key string, opt *WriteOption) error {
	_, err := c.kv.Delete(key, nil)
	return err
}

func (c *Consul) convertPairs(ps api.KVPairs) Pairs {
	var pairs Pairs
	for _, p := range ps {
		pairs = append(pairs, c.convertPair(p))
	}
	return pairs
}

func (c *Consul) convertPair(p *api.KVPair) *Pair {
	if p == nil {
		return nil
	}

	return &Pair{Key: p.Key, Value: p.Value}
}

func (c *Consul) done(ctx context.Context) bool {
	if ctx == nil {
		return true
	}

	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
