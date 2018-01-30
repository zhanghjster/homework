package conf

import (
	"context"

	"github.com/spf13/viper"
)

type Configure struct {
	viper.Viper
	remote Provider
}

func New() *Configure { return &Configure{} }

func (c *Configure) SetRemoteProvider(provider ProviderType, addr string) (err error) {
	c.remote, err = NewRemoteProvider(provider, addr)
	return
}

func (c *Configure) GetRemote(ctx context.Context, key string, opt *QueryOption) (*KVPair, *QueryMeta, error) {
	return c.remote.GetRemote(ctx, key, opt)
}

func (c *Configure) GetRemoteList(ctx context.Context, prefix string, opt *QueryOption) (KVPairs, *QueryMeta, error) {
	return c.remote.GetRemoteList(ctx, prefix, opt)
}

func (c *Configure) WatchRemote(ctx context.Context, key string, opt *QueryOption) (<-chan *KVPair, error) {
	return c.remote.WatchRemote(ctx, key, opt)
}
