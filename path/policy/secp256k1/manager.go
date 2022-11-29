package eth

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type handler struct {
	callback *callback
}

type callback struct {
	Storage *storage.Core
}

type manager struct {
	handler *handler
}

func newManager() *manager {
	return &manager{
		handler: &handler{
			callback: &callback{
				Storage: storage.Standard(),
			},
		},
	}
}

func (m *manager) Path() []*framework.Path {
	out := make([]*framework.Path, 0)
	out = append(out, m.handler.admin()...)
	return out
}
