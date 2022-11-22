package wallet

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"math/big"
)

func (pmgr *pathWallet) walletSignTxPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
			doc.FieldAddress:   {Type: framework.TypeString, Required: true},
			doc.FieldChainType: {Type: framework.TypeString, Default: "ETH"},
			doc.FieldTxBinary:  {Type: framework.TypeString},
			doc.FieldChainId:   {Type: framework.TypeInt64, Default: int64(1)},
		},
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: pmgr.signTxCallBack,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: pmgr.signTxCallBack,
			},
		},
		HelpSynopsis:    doc.PathSignSyn,
		HelpDescription: doc.PathSignDesc,
	}
}

func (pmgr *pathWallet) signTxCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
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

	// TODO：ChainId 约束: Config中设置
	chainId := data.Get(doc.FieldChainId).(int64)

	policy, err := pmgr.Storage.Policy.Read(ctx, req, namespace)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("policy not found")
	}

	// 做chainId约束
	cId := big.NewInt(chainId)
	if unsignTx.ChainId() == nil || unsignTx.ChainId().Cmp(cId) != 0 {
		return nil, types.ErrInvalidChainId
	}

	// 交易验证
	ok, err := policy.IsTxAllowed(&unsignTx, uint64(chainId))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("tx not allowed")
	}

	// 获取目标钱包
	wallet, err := pmgr.Storage.Wallet.Read(ctx, req, namespace, address)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.SignEthTx(&unsignTx, chainId)
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

func ethTxResponseData(tx *types.Transaction) (map[string]any, error) {
	binary, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		doc.FieldTxBinary: hexutil.Encode(binary),
		doc.FieldTxHash:   tx.Hash().Hex(),
	}, nil
}
