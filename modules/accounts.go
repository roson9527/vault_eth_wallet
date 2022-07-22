package modules

// Account is an Ethereum account
type Account struct {
	Address      string `json:"address"` // Ethereum account address derived from the private key
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`    // Ethereum public key derived from the private key
	CreationTime int64  `json:"creation_time"` // key pair creation time
}
