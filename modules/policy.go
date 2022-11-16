package modules

import (
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Policy struct {
	ChainIds             []uint64                  `json:"chain_ids" hcl:"chain_ids" mapstructure:"chain_ids"`                                        // 配置ChainID允许
	EnableCreateContract bool                      `json:"enable_create_contract" hcl:"enable_create_contract" mapstructure:"enable_create_contract"` // 是否允许创建合约
	Contract             map[string]ContractConfig `json:"contract" hcl:"contract" mapstructure:"contract"`                                           // 操作合约白名单
}

func (p *Policy) IsTxAllowed(tx *types.Transaction, chainId uint64) (bool, error) {
	if !p.IsChainIdAllowed(chainId) {
		return false, errors.New("chainId not allowed")
	}

	if tx.To() == nil {
		if !p.EnableCreateContract {
			return false, errors.New("create contract not allowed")
		}
		if tx.Value().Cmp(big.NewInt(0)) != 0 {
			return false, errors.New("create contract can not have value")
		}
		return true, nil
	}

	contract := tx.To().String()
	contractConfig := p.IsContractAllowed(contract, chainId)
	if contractConfig == nil {
		return false, errors.New("contract to address not allowed")
	}
	if tx.Data() == nil {
		return false, errors.New("contract.data is nil")
	}
	if len(tx.Data()) < 4 {
		return false, errors.New("contract.data is too short")
	}
	funcSign := hexutil.Encode(tx.Data()[0:4])
	funcSignConfig := contractConfig.IsFuncSignAllowed(funcSign)
	if funcSignConfig == nil {
		return false, errors.New("contract.funcSign not allowed")
	}
	if funcSignConfig.MaxValue != "" {
		if !funcSignConfig.IsMaxValueAllowed(tx.Value()) {
			return false, errors.New("contract.value is too big")
		}
	}
	return true, nil
}

func (p *Policy) IsChainIdAllowed(chainId uint64) bool {
	for _, id := range p.ChainIds {
		if id == chainId {
			return true
		}
	}
	return false
}

func (p *Policy) IsContractAllowed(contract string, chainId uint64) *ContractConfig {
	for k, c := range p.Contract {
		if k == contract || c.IsContractAllowed(contract, chainId) {
			return &c
		}
	}
	return nil
}

type ContractConfig struct {
	Address   string              `json:"address,omitempty" hcl:"address,omitempty" mapstructure:"address,omitempty"`       // 合约地址
	FuncSigns map[string]FuncSign `json:"func_sign" hcl:"func_sign" mapstructure:"func_sign"`                               // 函数签名
	ChainIds  []uint64            `json:"chain_ids,omitempty" hcl:"chain_ids,omitempty" mapstructure:"chain_ids,omitempty"` // 配置ChainID允许
}

func (c *ContractConfig) IsContractAllowed(contract string, chainId uint64) bool {
	return c.Address == contract && c.IsChainIdAllowed(chainId)
}

func (c *ContractConfig) IsFuncSignAllowed(funcSign string) *FuncSign {
	for k, f := range c.FuncSigns {
		if k == funcSign || f.IsFuncSignAllowed(funcSign) {
			return &f
		}
	}
	return nil
}

func (c *ContractConfig) IsChainIdAllowed(chainId uint64) bool {
	if c.ChainIds == nil || len(c.ChainIds) == 0 {
		return true
	}

	for _, id := range c.ChainIds {
		if id == chainId {
			return true
		}
	}
	return false
}

type FuncSign struct {
	Sign     string `json:"sign,omitempty" hcl:"sign,omitempty" mapstructure:"sign,omitempty"` // 函数签名
	MaxValue string `json:"max_value" hcl:"max_value" mapstructure:"max_value"`                // 最大值
}

func (s *FuncSign) IsFuncSignAllowed(sign string) bool {
	if s.Sign == "0xFFFFFFFF" {
		return true
	}
	return s.Sign == sign
}

func (s *FuncSign) IsMaxValueAllowed(value *big.Int) bool {
	maxV, ok := big.NewInt(0).SetString(s.MaxValue, 0)
	if !ok {
		return false
	}
	return value.Cmp(maxV) <= 0
}
