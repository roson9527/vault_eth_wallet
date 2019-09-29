package addresses

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAddress(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldAddress: {Type: framework.TypeString},
		},
		// 执行的位置，有read，list，create，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: read,
			},
		},
		// TODO:ExistenceCheck
		ExistenceCheck:  nil,
		HelpSynopsis:    pathAddressSyn,
		HelpDescription: pathAddressDesc,
	}
}
