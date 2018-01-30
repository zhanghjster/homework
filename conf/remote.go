package conf

import (
	"context"

	"errors"

	"github.com/hashicorp/consul/api"
)

type ProviderType int

const (
	ProviderConsul ProviderType = iota
	// TODO: support etcd
)

type QueryMeta struct {
	LastIndex uint64
}

type QueryOption struct {
	WithPrefix bool
	LastIndex  uint64
}

type Provider interface {
	GetRemote(ctx context.Context, key string, opt *QueryOption) (*KVPair, *QueryMeta, error)
	GetRemoteList(ctx context.Context, prefix string, opt *QueryOption) (KVPairs, *QueryMeta, error)
	WatchRemote(ctx context.Context, key string, opt *QueryOption) (<-chan *KVPair, error)
}

type KVPair struct {
	Key   string
	Value []byte
}

type KVPairs []*KVPair

type Consul struct {
	Addr string
	kv   *api.KV
}

func NewRemoteProvider(pt ProviderType, addr string) (p Provider, err error) {
	switch pt {
	case ProviderConsul:
		return newConsul(addr)
	default:
		err = errors.New("provider not supported")
	}

	return p, err
}

func newConsul(addr string) (*Consul, error) {
	var client, err = api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}

	return &Consul{
		Addr: addr, kv: client.KV(),
	}, nil
}

func (c *Consul) GetRemote(ctx context.Context, key string, opt *QueryOption) (*KVPair, *QueryMeta, error) {
	p, mt, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, nil, err
	}

	return c.convertPair(p), &QueryMeta{LastIndex: mt.LastIndex}, nil
}

func (c *Consul) GetRemoteList(ctx context.Context, prefix string, opt *QueryOption) (KVPairs, *QueryMeta, error) {
	l, mt, err := c.kv.List(prefix, nil)
	if err != nil {
		return nil, nil, err
	}

	return c.convertPairs(l), &QueryMeta{LastIndex: mt.LastIndex}, nil
}

type watchFunc func() error

func (c *Consul) WatchRemote(ctx context.Context, key string, opt *QueryOption) (<-chan *KVPair, error) {
	var ch = make(chan *KVPair)

	var qOpt = (&api.QueryOptions{}).WithContext(ctx)
	if opt != nil {
		qOpt.WaitIndex = opt.LastIndex
	}

	var watch watchFunc
	if opt.WithPrefix {
		watch = c.watchItemListFunc(key, ch, qOpt)
	} else {
		watch = c.watchItemFunc(key, ch, qOpt)
	}

	go func() {
		defer close(ch)
		for {
			if err := watch(); err != nil {
				return
			}
		}
	}()

	return ch, nil
}

func (c *Consul) watchItemFunc(key string, ch chan *KVPair, qOpt *api.QueryOptions) watchFunc {
	var err error
	var p *api.KVPair
	var mt *api.QueryMeta

	return func() error {
		if mt != nil {
			qOpt.WaitIndex = mt.LastIndex
		}

		p, mt, err = c.kv.Get(key, qOpt)
		if err != nil {
			return err
		}

		if p != nil {
			ch <- c.convertPair(p)
		}

		return nil
	}
}

func (c *Consul) watchItemListFunc(prefix string, ch chan *KVPair, qOpt *api.QueryOptions) watchFunc {
	var err error
	var l api.KVPairs
	var mt *api.QueryMeta

	return func() error {
		if mt != nil {
			qOpt.WaitIndex = mt.LastIndex
		}

		l, mt, err = c.kv.List(prefix, qOpt)
		if err != nil {
			return err
		}

		for _, item := range c.convertPairs(l) {
			ch <- item
		}

		return nil
	}
}

func (c *Consul) convertPair(p *api.KVPair) *KVPair {
	if p == nil {
		return nil
	}

	return &KVPair{Key: p.Key, Value: p.Value}
}

func (c *Consul) convertPairs(ps api.KVPairs) KVPairs {
	var pairs KVPairs
	for _, p := range ps {
		pairs = append(pairs, c.convertPair(p))
	}
	return pairs
}

/*
type Etcd struct {
	Addr   string
	client *clientv3.Client
}

func newEtcd(addr string) (*Etcd, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{addr},
	})
	if err != nil {
		return nil, err
	}

	return &Etcd{Addr: addr, client: client}, nil
}

func (e *Etcd) GetRemote(ctx context.Context, key string, opt *QueryOption) (*KVPair, *QueryMeta, error) {
	var item *KVPair

	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	for _, kv := range resp.Kvs {
		item = &KVPair{
			Key:   string(kv.Key),
			Value: kv.Value,
		}
		break
	}

	qMeta := &QueryMeta{
		LastIndex: uint64(resp.Header.Revision),
	}

	return item, qMeta, nil
}

func (e *Etcd) GetRemoteList(ctx context.Context, prefix string, opt *QueryOption) (KVPairs, *QueryMeta, error) {
	var items KVPairs

	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return items, nil, err
	}

	for _, kv := range resp.Kvs {
		items = append(items, &KVPair{
			Key:   string(kv.Key),
			Value: kv.Value,
		})
	}

	qMeta := &QueryMeta{
		LastIndex: uint64(resp.Header.Revision),
	}

	return items, qMeta, nil
}

func (e *Etcd) WatchRemote(ctx context.Context, key string, opt *QueryOption) (<-chan *KVPair, error) {
	var rCh = make(chan *KVPair)

	go func() {
		var opts []clientv3.OpOption
		if opt != nil {
			if opt.WithPrefix {
				opts = append(opts, clientv3.WithPrefix())
			}
			if opt.LastIndex > 0 {
				opts = append(opts, clientv3.WithRev(int64(opt.LastIndex)+1))
			}
		}

		wCh := e.client.Watch(ctx, key, opts...)
		for resp := range wCh {
			for _, ev := range resp.Events {
				rCh <- &KVPair{
					Key:   string(ev.Kv.Key),
					Value: ev.Kv.Value,
				}
			}
		}
	}()

	return rCh, nil
}

*/
