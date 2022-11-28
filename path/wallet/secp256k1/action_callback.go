package secp256k1

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (cb *callback) read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	wallet, err := cb.readRaw(ctx, req, data)
	if err != nil {
		return nil, err
	}
	chain := data.Get(doc.FieldChain).(string)
	return &logical.Response{
		Data: walletResponseData(wallet, chain, false),
	}, nil
}

func (cb *callback) signTx(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	chain := data.Get(doc.FieldChain).(string)
	address := data.Get(doc.FieldAddress).(string)
	binaryStr := data.Get(doc.FieldTxBinary).(string)

	var unsignTx types.Transaction
	binary, err := hexutil.Decode(binaryStr)
	if err != nil {
		return nil, err
	}
	err = unsignTx.UnmarshalBinary(binary)
	if err != nil {
		return nil, err
	}

	policy, err := cb.Storage.Policy.Read(ctx, req, namespace, doc.CryptoSECP256K1, chain)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("policy not found")
	}

	// 做chainId约束, 不支持旧交易
	if unsignTx.ChainId() == nil || unsignTx.ChainId().Uint64() == 0 {
		return nil, types.ErrInvalidChainId
	}

	// 交易验证
	ok, err := policy.IsTxAllowed(&unsignTx, unsignTx.ChainId().Uint64())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("tx not allowed")
	}

	// 获取目标钱包
	wallet, err := cb.Storage.read(ctx, req, namespace, chain, address)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.SignEthTx(&unsignTx)
	if err != nil {
		return nil, err
	}

	responseData, err := ethTxResponseData(tx)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: responseData,
	}, nil
}

func (cb *callback) listAlias(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	chain := data.Get(doc.FieldChain).(string)
	aliasList, err := cb.Storage.Alias.List(ctx, req, namespace, aliasType(chain))
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldKeys: aliasList,
		},
	}, nil
}
