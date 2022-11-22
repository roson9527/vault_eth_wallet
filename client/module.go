package client

import "github.com/roson9527/vault_eth_wallet/modules"

type Wallet struct {
	PrivateKey string              `json:"private_key"`          // PrivateKey is the private key of the wallet
	PublicKey  string              `json:"public_key,omitempty"` // PublicKey is the public key of the wallet
	Address    string              `json:"address"`
	UpdateTime int64               `json:"update_time"` // key pair update time
	NameSpaces []string            `json:"namespaces,omitempty"`
	Network    string              `json:"network,omitempty"`
	Meta       modules.WalletExtra `json:"meta,omitempty"`
}
