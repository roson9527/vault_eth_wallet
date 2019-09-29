package accounts

const (
	pathListDesc = `All the Ethereum accounts will be listed.`
	pathListSyn  = `List all the Ethereum accounts at a path.`

	pathSignDesc = `Sign data using a given Ethereum account.`
	pathSignSyn  = `Sign data`

	pathSignTxDesc = `Sign transaction using a given Ethereum account.`
	pathSignTxSyn  = `Sign a provided transaction.`

	errCastingPubToECDSA = `error casting public key to ECDSA`
	errReadAccountByName = `Error reading account.`

	fieldDataDesc = `The data to hash (keccak) and sign.`

	fieldName      = "name"
	fieldAddress   = "address"
	fieldData      = "data"
	fieldEncoding  = "encoding"
	fieldIsHash    = "isHash"
	fieldSigned    = "signed"
	fieldChainId   = "chainId"
	fieldGasLimit  = "gas_limit"
	fieldGasPrice  = "gas_price"
	fieldNonce     = "nonce"
	fieldToAddress = "to_address"
	fieldAmount    = "amount"

	patternStr        = "accounts/"
	patternAddressStr = "addresses/"
)

// 一些默认的签名用参数固定
// 用于测试
