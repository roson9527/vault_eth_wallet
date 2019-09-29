package modules

// Config contains the configuration for each mount
type Config struct {
	//BoundCIDRList       []string `json:"bound_cidr_list_list" structs:"bound_cidr_list" mapstructure:"bound_cidr_list"`
	RPC string `json:"rpc_url"` // 用于读取Nonce
	//CoinMarketCapAPIKey string   `json:"api_key"`
	ChainID string `json:"chain_id"` // 配置ChainID

}
