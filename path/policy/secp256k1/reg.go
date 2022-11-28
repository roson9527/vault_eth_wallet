package eth

import "github.com/roson9527/vault_eth_wallet/path/factory"

var aliasType = "policy"

func init() {
	// Register the wallet type
	factory.Register(aliasType, newManager())
}
