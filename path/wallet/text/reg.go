package text

import "github.com/roson9527/vault_eth_wallet/path/factory"

func init() {
	// Register the wallet type
	factory.Register(aliasType(""), newManager())
}
