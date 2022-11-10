package wallet

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"math/big"
)

func (pmgr *pathWallet) walletSignTxPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNameSpace: {Type: framework.TypeString, Required: true},
			fieldAddress:   {Type: framework.TypeString, Required: true},
			fieldChainType: {Type: framework.TypeString, Default: "ETH"},
			fieldTxBinary:  {Type: framework.TypeString},
			fieldChainId:   {Type: framework.TypeInt64, Default: int64(1)},
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
		HelpSynopsis:    pathSignSyn,
		HelpDescription: pathSignDesc,
	}
}

func (pmgr *pathWallet) signTxCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(fieldNameSpace).(string)
	address := data.Get(fieldAddress).(string)

	binaryStr := data.Get(fieldTxBinary).(string)
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
	chainId := data.Get(fieldChainId).(int64)

	policy, err := pmgr.policyStorage.readPolicy(ctx, req, namespace)
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
	wallet, err := pmgr.walletStorage.readWallet(ctx, req, namespace, address)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.SignETH(&unsignTx, chainId)
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
		fieldTxBinary: hexutil.Encode(binary),
		fieldTxHash:   tx.Hash().Hex(),
	}, nil
}
