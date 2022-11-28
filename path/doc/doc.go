package doc

const (
	NameSpaceGlobal = "global"
)

const (
	PathListDesc   = `All the addresses of accounts will be listed.`
	PathListSyn    = `List all the account addresses.`
	PathCreateSyn  = `Create a new account.`
	PathCreateDesc = `Create a new account.`
	PathReadSyn    = `Read an account.`
	PathReadDesc   = `Read an account.`
	PathSignSyn    = `SignEthTx data.`
	PathSignDesc   = `SignEthTx data.`
)

const (
	PathSubSignTx = "/sign_tx"
	PathSubExport = "/export"
	PathSubNew    = "new"
)

const (
	FieldMnemonic     = "mnemonic"
	FieldGAuth        = "g_auth"
	FieldSocialID     = "social_id"
	FieldApp          = "app"
	FieldUser         = "user"
	FieldAddress      = "address"
	FieldAddressAlias = "address_alias"
	FieldCryptoType   = "crypto_type"
	FieldChain        = "chain"
	FieldExtra        = "extra"
	FieldUpdateTime   = "update_time"
	FieldNameSpace    = "namespace"
	FieldNameSpaces   = "namespaces"
	FieldPrivateKey   = "private_key"
	FieldPublicKey    = "public_key"
	FieldTxBinary     = "tx_binary"
	FieldTxHash       = "tx_hash"
	FieldPolicyHCL    = "policy_hcl"
	FieldKeys         = "keys"
	FieldPolicy       = "policy"
)

const (
	AliasWallet = "wallet/%s/%s"
	AliasSocial = "social/%s"
)

const (
	ChainETH     = "eth"
	ChainDefault = "default"
)

const (
	AppDiscord = "discord"
	AppTwitter = "twitter"
)

const (
	CryptoSECP256K1 = "secp256k1"
	CryptoTEXT      = "text"
)
