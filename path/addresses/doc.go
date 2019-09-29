package addresses

const (
	pathListDesc = `All the addresses of accounts will be listed.`
	pathListSyn  = `List all the account addresses.`

	pathAddressDesc = `Lookup a account's name by address.`
	pathAddressSyn  = pathAddressDesc

	pathVerifyDesc = `Verify that data was signed by a particular address`
	pathVerifySyn  = pathVerifyDesc

	fieldDataDesc      = `The data to verify the signature of.`
	fieldSignatureDesc = `The signature to verify.`

	fieldAddress     = "address"
	fieldData        = "data"
	fieldSignature   = "signature"
	fieldVerified    = "verified"
	fieldAccountName = "name"
	fieldEncoding    = "encoding"
	fieldIsHash      = "isHash"

	patternStr = "addresses/"

	errCastingPubToECDSA  = `error casting public key to ECDSA`
	errSignCheckFailed    = `signature not verified`
	errDeserializeAccount = `failed to deserialize account at %s`
)
