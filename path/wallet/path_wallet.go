package wallet

import (
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type handler struct {
	callback callback
}

type callback struct {
	Storage *storage.Core
}

func walletResponseData(wallet *modules.Wallet, extra bool) map[string]any {
	out := map[string]any{
		doc.FieldAddress:    wallet.Address,
		doc.FieldPublicKey:  wallet.PublicKey,
		doc.FieldUpdateTime: wallet.UpdateTime,
		doc.FieldNetwork:    wallet.Network,
	}
	if extra {
		out[doc.FieldPrivateKey] = wallet.PrivateKey
		out[doc.FieldNameSpaces] = wallet.NameSpaces
		out[doc.FieldExtra] = wallet.Extra
	}
	return out
}
