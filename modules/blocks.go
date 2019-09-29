package modules

// Transaction is an Ethereum transaction
type Transaction struct {
	Value     string `json:"value"`
	Gas       uint64 `json:"gas"`
	GasPrice  uint64 `json:"gas_price"`
	Nonce     uint64 `json:"nonce"`
	AddressTo string `json:"address_to"`
}
