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

	fieldName         = "name"
	fieldAddress      = "address"
	fieldData         = "data"
	fieldEncoding     = "encoding"
	fieldIsHash       = "is_hash"
	fieldSigned       = "signed"
	fieldChainId      = "chain_id"
	fieldGasLimit     = "gas_limit"
	fieldGasPrice     = "gas_price"
	fieldNonce        = "nonce"
	fieldAddressTo    = "address_to"
	fieldAmount       = "amount"
	fieldCreationTime = "creation_time"

	valueUTF8 = "utf8"

	fieldTransactionHash   = "transaction_hash"
	fieldSignedTransaction = "signed_transaction"
	fieldSignHash          = "sign_hash"

	patternStr        = "accounts/"
	patternAddressStr = "addresses/"
)

// 一些默认的签名用参数固定
// 用于测试
