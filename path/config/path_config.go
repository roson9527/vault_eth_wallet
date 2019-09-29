package config

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		Fields: map[string]*framework.FieldSchema{
			fieldChainId: {
				Type:        framework.TypeString,
				Description: chainIdDesc,
			},
			fieldRpcUrl: {
				Type:        framework.TypeString,
				Description: rpcUrlDesc,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: write,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: write,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: read,
			},
		},
		HelpSynopsis:    pathConfigSyn,
		HelpDescription: pathConfigDesc,
	}
}
