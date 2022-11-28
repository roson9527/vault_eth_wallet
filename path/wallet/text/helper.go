package text

import (
	"fmt"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func aliasType(chain string) string {
	return fmt.Sprintf(doc.AliasWallet, doc.CryptoTEXT, chain)
}

func walletResponseData(wallet *modules.WalletExtra, chain string, extra bool) map[string]any {
	out := map[string]any{
		doc.FieldAddress:      wallet.GetAddress(chain),
		doc.FieldPublicKey:    wallet.PublicKey,
		doc.FieldUpdateTime:   wallet.UpdateTime,
		doc.FieldAddressAlias: wallet.AddressAlias,
	}
	if extra {
		out[doc.FieldCryptoType] = wallet.CryptoType
		out[doc.FieldPrivateKey] = wallet.PrivateKey
		out[doc.FieldNameSpaces] = wallet.NameSpaces
		out[doc.FieldExtra] = wallet.Extra
		out[doc.FieldMnemonic] = wallet.Mnemonic
	}
	return out
}
