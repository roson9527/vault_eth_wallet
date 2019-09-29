package accounts

import (
	"github.com/hashicorp/vault/sdk/framework"
)

func Path() []*framework.Path {
	return []*framework.Path{
		pathList(patternStr + "?"),
		pathAccount(patternStr + framework.GenericNameRegex(fieldName)),
		pathSignTx(patternStr + framework.GenericNameRegex(fieldName) + "/sign-tx"),
		pathSign(patternStr + framework.GenericNameRegex(fieldName) + "/sign"),
	}
}
