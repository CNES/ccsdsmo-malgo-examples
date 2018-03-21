package provider

import (
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
)

// Define Provider's structure
type Provider struct {
	ctx  *Context
	cctx *ClientContext
}

// Allow to close the context of a specific provider
func (provider *Provider) Close() {
	provider.ctx.Close()
}

// Create a provider
func CreateProvider(url string) (*Provider, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}

	provider := &Provider{ctx, cctx}

	return provider, nil
}
