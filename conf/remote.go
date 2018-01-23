package conf

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type ProviderType int

const (
	ProviderEtcd ProviderType = iota
	ProviderConsul
)

type Provider interface {
	GetRemote(ctx context.Context, key string) (*RemoteItem, error)
	GetRemoteList(ctx context.Context, prefix string) (RemoteItems, error)
	WatchRemote(ctx context.Context, key string) (<-chan *RemoteItem, error)
	WatchRemoteList(ctx context.Context, prefix string) (<-chan RemoteItems, error)
}

type RemoteItem struct {
	Key   string
	Value []byte
}

type RemoteItems []*RemoteItem

type Consul struct {
	Addr string
	kv   *api.KV
}

func (c *Consul) GetRemote(ctx context.Context, key string) (*RemoteItem, error) {
	p, _, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return c.convertPair(p), nil
}

func (c *Consul) GetRemoteList(ctx context.Context, prefix string) (RemoteItems, error) {
	l, _, err := c.kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	return c.convertPairs(l), nil
}

func (c *Consul) WatchRemote(ctx context.Context, key string) (<-chan *RemoteItem, error) {
	p, mt, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var ch = make(chan *RemoteItem)
	go func() {
		defer close(ch)

		for {
			ch <- c.convertPair(p)

			var qOpt = api.QueryOptions{
				WaitIndex: mt.LastIndex,
			}.WithContext(ctx)

			p, mt, err = c.kv.Get(key, qOpt)
		}
	}()

	return ch, nil
}

func (c *Consul) WatchRemoteList(ctx context.Context, prefix string) (<-chan RemoteItems, error) {
	l, mt, err := c.kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	var ch = make(chan RemoteItems)
	go func() {
		defer close(ch)
		for {
			ch <- c.convertPairs(l)

			var qOpt = api.QueryOptions{
				WaitIndex: mt.LastIndex,
			}.WithContext(ctx)

			l, mt, _ = c.kv.List(prefix, qOpt)
		}
	}()

	return ch, nil
}

func (c *Consul) convertPair(p *api.KVPair) *RemoteItem {
	if p == nil {
		return nil
	}

	return &RemoteItem{Key: p.Key, Value: p.Value}
}

func (c *Consul) convertPairs(ps api.KVPairs) RemoteItems {
	var pairs RemoteItems
	for _, p := range ps {
		pairs = append(pairs, c.convertPair(p))
	}
	return pairs
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

func NewRemoteProvider(pt ProviderType, addr string) (p Provider, err error) {
	switch pt {
	case ProviderEtcd:
		return newEtcd(addr)
	case ProviderConsul:
		return newConsul(addr)
	default:
		err = errors.New("provider not supported")
	}

	return p, err
}

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

func (e *Etcd) GetRemote(ctx context.Context, key string) (*RemoteItem, error) {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var item *RemoteItem
	for _, kv := range resp.Kvs {
		item = &RemoteItem{
			Key:   string(kv.Key),
			Value: kv.Value,
		}
		break
	}

	return item, nil
}

func (e *Etcd) GetRemoteList(ctx context.Context, prefix string) (RemoteItems, error) {
	var items RemoteItems
	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return items, err
	}

	for _, kv := range resp.Kvs {
		items = append(items, &RemoteItem{
			Key:   string(kv.Key),
			Value: kv.Value,
		})
	}

	return items, nil
}

func (e *Etcd) WatchRemote(ctx context.Context, key string) (<-chan *RemoteItem, error) {
	var rCh = make(chan *RemoteItem)
	go func() {
		wCh := e.client.Watch(ctx, key)
		for resp := range wCh {
			for _, ev := range resp.Events {
				rCh <- &RemoteItem{
					Key: string(ev.Kv.Key), Value: ev.Kv.Value,
				}
			}
		}
	}()
}

func (e *Etcd) WatchRemoteList(ctx context.Context, prefix string) (<-chan RemoteItems, error) {
	panic("implement me")
}
