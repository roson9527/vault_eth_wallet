package addresses

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/accounts"
	"github.com/roson9527/vault_eth_wallet/path/base"
)

func pathVerify(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldAddress: {Type: framework.TypeString},
			fieldData: {
				Type:        framework.TypeString,
				Description: fieldDataDesc,
			},
			fieldSignature: {
				Type:        framework.TypeString,
				Description: fieldSignatureDesc,
			},
		},
		// 执行的位置，有read，list，create，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: verify,
			},
		},
		// TODO:ExistenceCheck
		ExistenceCheck:  nil,
		HelpSynopsis:    pathVerifySyn,
		HelpDescription: pathVerifyDesc,
	}
}

func verify(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	addr := data.Get(fieldAddress).(string)
	address, err := readByAddr(ctx, req, addr)
	if err != nil {
		return nil, err
	}
	if address.Name == "" {
		return nil, nil
	}

	return verifySignature(ctx, req, data, address.Name)
}

func verifySignature(
	ctx context.Context, req *logical.Request, data *framework.FieldData, name string) (*logical.Response, error) {
	acct, err := accounts.ReadByName(ctx, req, name)
	if err != nil || acct == nil {
		return nil, err
	}

	signature := data.Get(fieldSignature).(string)
	dataToSign := data.Get(fieldData).(string)
	encoding := data.Get(fieldEncoding).(string)
	isHash := data.Get(fieldIsHash).(bool)

	dataBytes, err := base.FormatData(dataToSign, encoding)
	if err != nil {
		return nil, err
	}
	verified, err := base.Verify(acct, dataBytes, signature, isHash)

	return &logical.Response{
		Data: map[string]interface{}{
			fieldVerified: verified,
			fieldAddress:  acct.Address,
		},
	}, nil
}
