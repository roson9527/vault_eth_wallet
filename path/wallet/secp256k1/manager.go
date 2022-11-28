package secp256k1

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type handler struct {
	callback *callback
}

type callback struct {
	Storage *storageEx
}

type manager struct {
	handler *handler
}

func newManager() *manager {
	return &manager{
		handler: &handler{
			callback: &callback{
				Storage: &storageEx{
					storage.Standard(),
				},
			},
		},
	}
}

func (m *manager) Path() []*framework.Path {
	out := make([]*framework.Path, 0)
	out = append(out, m.handler.admin()...)
	out = append(out, m.handler.action()...)
	return out
}
