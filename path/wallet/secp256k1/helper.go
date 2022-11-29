package secp256k1

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

var (
	defAliasType = aliasType(doc.ChainETH)
)

func aliasType(chain string) string {
	return fmt.Sprintf(doc.AliasWallet, doc.CryptoSECP256K1, chain)
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
	}
	return out
}

func ethTxResponseData(tx *types.Transaction) (map[string]any, error) {
	binary, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		doc.FieldTxBinary: hexutil.Encode(binary),
		doc.FieldTxHash:   tx.Hash().Hex(),
	}, nil
}
