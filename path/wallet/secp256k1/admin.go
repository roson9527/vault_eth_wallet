package secp256k1

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) admin() []*framework.Path {
	return []*framework.Path{
		{
			// create
			Pattern: fmt.Sprintf(storage.PatternWalletAdmin, doc.NameSpaceGlobal, doc.CryptoSECP256K1),
			Fields: map[string]*framework.FieldSchema{
				//doc.FieldNameSpace:  {Type: framework.TypeString, Default: doc.NameSpaceGlobal},
				doc.FieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
				doc.FieldPrivateKey: {Type: framework.TypeString},
				doc.FieldExtra:      {Type: framework.TypeMap, Default: map[string]any{}},
			},
			ExistenceCheck: h.callback.addressExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				// 这个路径只能用于创建
				logical.CreateOperation: &framework.PathOperation{
					Callback: h.callback.create,
				},
			},
		},
		{
			// update + delete
			Pattern: fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.CryptoSECP256K1,
				framework.GenericNameRegex(doc.FieldAddress)),
			Fields: map[string]*framework.FieldSchema{
				doc.FieldAddress:      {Type: framework.TypeString, Required: true},
				doc.FieldNameSpaces:   {Type: framework.TypeCommaStringSlice, Default: []string{}},
				doc.FieldAddressAlias: {Type: framework.TypeKVPairs, Default: map[string]string{}},
				doc.FieldExtra:        {Type: framework.TypeMap, Default: map[string]any{}},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.update,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: h.callback.delete,
				},
			},
		},
		{
			// export
			Pattern: fmt.Sprintf(storage.PatternWalletByChain, framework.GenericNameRegex(doc.FieldNameSpace),
				doc.CryptoSECP256K1, framework.GenericNameRegex(doc.FieldChain),
				framework.GenericNameRegex(doc.FieldAddress)) + doc.PathSubExport,
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldChain:     {Type: framework.TypeString, Required: true},
				doc.FieldAddress:   {Type: framework.TypeString, Required: true},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.export,
				},
			},
		},
		{
			// list
			Pattern: fmt.Sprintf(storage.PatternWalletAdmin, doc.NameSpaceGlobal, doc.CryptoSECP256K1+"/?"),
			Fields:  map[string]*framework.FieldSchema{
				// 不需要任何参数
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: h.callback.list,
				},
			},
		},
	}
	//	HelpSynopsis:    doc.PathCreateSyn,
	//	HelpDescription: doc.PathCreateDesc,
}

func (cb *callback) addressExistenceCheck(ctx context.Context, request *logical.Request, data *framework.FieldData) (bool, error) {
	// 如果没有私钥，就是内部生成
	if data.Get(doc.FieldPrivateKey).(string) == "" {
		return false, nil
	}
	// 如果有私钥，那么就是外部导入，需要检查地址是否已经存在
	w, err := base.CryptoETH.PrivateToWallet(data.Get(doc.FieldPrivateKey).(string))
	if err != nil {
		return false, err
	}
	wIns, err := cb.Storage.Wallet.Read(ctx, request, doc.NameSpaceGlobal, doc.CryptoSECP256K1, w.Address)
	if err != storage.ErrWalletNotFound {
		return false, err
	}
	if wIns != nil {
		return true, errors.New(fmt.Sprintf("address [%s] already exists", w.Address))
	}
	return false, nil
}
