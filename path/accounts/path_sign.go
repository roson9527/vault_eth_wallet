package accounts

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/utils"
)

func pathSignTx(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldName: {Type: framework.TypeString},
			fieldAddressTo: {
				Type:        framework.TypeString,
				Description: "The address of the account to send ETH to.",
			},
			fieldData: {
				Type:        framework.TypeString,
				Description: "The data to sign.",
			},
			fieldEncoding: {
				Type:        framework.TypeString,
				Default:     valueUTF8,
				Description: "The encoding of the data to sign.",
			},
			fieldAmount: {
				Type:        framework.TypeString,
				Description: "Amount of ETH (in wei).",
			},
			fieldNonce: {
				Type:        framework.TypeString,
				Description: "The transaction nonce.",
			},
			fieldGasLimit: {
				Type:        framework.TypeString,
				Description: "The gas limit for the transaction - defaults to 21000.",
				Default:     "21000",
			},
			fieldGasPrice: {
				Type:        framework.TypeString,
				Description: "The gas price for the transaction in wei.",
				Default:     "0",
			},
		},
		// 执行的位置，有read，list，create，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: signTx,
			},
		},
		ExistenceCheck:  base.PathExistenceCheck,
		HelpSynopsis:    pathSignTxSyn,
		HelpDescription: pathSignTxDesc,
	}
}

func pathSign(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldName: {Type: framework.TypeString},
			fieldData: {
				Type:        framework.TypeString,
				Description: fieldDataDesc,
			},
			fieldEncoding: {
				Type:        framework.TypeString,
				Default:     valueUTF8,
				Description: "The encoding of the data to sign.",
			},
			fieldIsHash: {
				Type:    framework.TypeBool,
				Default: false,
			},
		},
		// 执行的位置，有read，list，create，update
		ExistenceCheck: base.PathExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: sign,
			},
		},
		HelpSynopsis:    pathSignSyn,
		HelpDescription: pathSignDesc,
	}
}

func sign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get(fieldName).(string)
	signData := data.Get(fieldData).(string)
	encoding := data.Get(fieldEncoding).(string)
	isHash := data.Get(fieldIsHash).(bool)

	dataToSign, err := base.FormatData(signData, encoding)

	account, err := ReadByName(ctx, req, name)
	if err != nil {
		return nil, err
	}

	signRet, err := base.Signature(account, &modules.SignParams{
		Data:       dataToSign,
		IsHashData: isHash,
	})
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			fieldSigned:   signRet.Signed,
			fieldSignHash: signRet.TransactionHash,
		},
	}, nil
}

func signTx(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// 基本属性读取
	// * 包括nonce需要包含在dataOrFile里 *
	name := data.Get(fieldName).(string)

	account, err := ReadByName(ctx, req, name)
	if err != nil {
		return nil, err
	}

	signParams, err := readParams(data)
	if err != nil {
		return nil, err
	}

	signRet, err := base.SignatureTx(account, signParams)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			fieldTransactionHash:   signRet.TransactionHash,
			fieldSignedTransaction: signRet.Signed,
			fieldAddress:           account.Address,
		},
	}, nil
}

func readParams(data *framework.FieldData) (*modules.SignParams, error) {
	txData := data.Get(fieldData).(string)
	encoding := data.Get(fieldEncoding).(string)
	chainId := data.Get(fieldChainId).(string)
	gasLimit := data.Get(fieldGasLimit).(string)
	gasPrice := data.Get(fieldGasPrice).(string)
	nonce := data.Get(fieldNonce).(string)
	toAddress := data.Get(fieldAddressTo).(string)
	amount := data.Get(fieldAmount).(string)
	isHash := data.Get(fieldIsHash).(bool)

	addr := common.HexToAddress(toAddress)
	fd, err := base.FormatData(txData, encoding)
	if err != nil {
		return nil, err
	}
	// TODO 处理一些默认值的问题, 比如自动ChainID, gas***, nonce等
	return &modules.SignParams{
		Nonce:      utils.ValidNumber(nonce).Uint64(),
		ToAddress:  &addr,
		Amount:     utils.ValidNumber(amount),
		GasLimit:   utils.ValidNumber(gasLimit).Uint64(),
		GasPrice:   utils.ValidNumber(gasPrice),
		Data:       fd,
		IsHashData: isHash,
		ChainId:    utils.ValidNumber(chainId),
	}, nil
}
