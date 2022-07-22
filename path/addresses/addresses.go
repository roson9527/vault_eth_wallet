package addresses

import (
	"github.com/hashicorp/vault/sdk/framework"
)

func Path() []*framework.Path {
	return []*framework.Path{
		pathList(patternStr + "?"),
		pathAddress(patternStr + framework.GenericNameRegex(fieldAddress)),
		pathVerify(patternStr + framework.GenericNameRegex(fieldAddress) + "/verify"),
	}
}
