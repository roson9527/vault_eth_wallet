package socialid

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type handler struct {
	callback *callback
}

type callback struct {
	Storage *storageEx
}

type storageEx struct {
	*storage.Core
}

type manager struct {
	handler *handler
}

func newManager() *manager {
	return &manager{
		handler: &handler{
			callback: &callback{
				&storageEx{
					storage.Standard(),
				},
			},
		},
	}
}

func (m *manager) Path() []*framework.Path {
	out := make([]*framework.Path, 0)
	for _, app := range []string{doc.AppDiscord, doc.AppTwitter} {
		out = append(out, m.handler.admin(app)...)
	}
	out = append(out, m.handler.action()...)
	return out
}
