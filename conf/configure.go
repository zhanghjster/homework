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

func (c *Configure) GetRemote(ctx context.Context, key string) (*RemoteItem, error) {
	return c.remote.GetRemote(ctx, key)
}

func (c *Configure) GetRemoteList(ctx context.Context, prefix string) (RemoteItems, error) {
	return c.remote.GetRemoteList(ctx, prefix)
}

func (c *Configure) WatchRemote(ctx context.Context, key string) (<-chan *RemoteItem, error) {
	return c.remote.WatchRemote(ctx, key)
}

func (c *Configure) WatchRemoteList(ctx context.Context, prefix string) (<-chan RemoteItems, error) {
	return c.remote.WatchRemoteList(ctx, prefix)
}
