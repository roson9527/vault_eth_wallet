package modules

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// (nonce, toAddress, amount, gasLimit, gasPrice, txDataToSign)
type SignParams struct {
	Nonce      uint64
	ToAddress  *common.Address
	Amount     *big.Int
	GasLimit   uint64
	GasPrice   *big.Int
	Data       []byte
	IsHashData bool

	ChainId *big.Int
}

type SignResult struct {
	Signed          string `json:"signed"`
	TransactionHash string `json:"transaction_hash"`
}
