package config

const (
	pathConfigDesc = `Configure the trustee plugin.`
	pathConfigSyn  = `Configure the trustee plugin.`

	chainIdDesc = `Ethereum network ID`
	rpcUrlDesc  = `The RPC address of the Ethereuem network.`

	storageKey = "config"

	fieldChainId = "chain_id"
	fieldRpcUrl  = "rpc_url"

	errNotConfiguredEthBackend = `the ethereum backend is not configured properly`
)
