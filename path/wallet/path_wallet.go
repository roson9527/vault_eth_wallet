package wallet

import (
	"github.com/roson9527/vault_eth_wallet/modules"
)

type pathWallet struct {
	*storage
}

func walletResponseData(wallet *modules.Wallet, extra bool) map[string]any {
	out := map[string]any{
		fieldAddress:    wallet.Address,
		fieldPublicKey:  wallet.PublicKey,
		fieldUpdateTime: wallet.UpdateTime,
		fieldNetwork:    wallet.Network,
	}
	if extra {
		out[fieldPrivateKey] = wallet.PrivateKey
		out[fieldNameSpaces] = wallet.NameSpaces
	}
	return out
}
