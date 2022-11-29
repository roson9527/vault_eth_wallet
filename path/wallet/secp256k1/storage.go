package secp256k1

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type storageEx struct {
	*storage.Core
}

func (s *storageEx) create(ctx context.Context, req *logical.Request, payload *modules.WalletExtra) (*modules.WalletExtra, error) {
	if !(payload.Address != "" && payload.Address == payload.AddressAlias[doc.ChainETH]) {
		return nil, errors.New("default address is nil")
	}

	w, err := s.Wallet.Create(ctx, req, payload)
	if err != nil {
		return nil, err
	}

	srcPath := fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.CryptoSECP256K1, payload.Address)
	hclog.Default().Info("secp256k1:create", "path.src", srcPath)
	err = s.Alias.Create(ctx, req, defAliasType, payload.Address, srcPath, w.NameSpaces)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *storageEx) update(ctx context.Context, req *logical.Request, payload *modules.WalletExtra) (*modules.WalletExtra, error) {
	if !(payload.Address != "" && payload.Address == payload.AddressAlias[doc.ChainETH]) {
		return nil, errors.New("default address is nil")
	}

	// 获取目标钱包
	oldWallet, err := s.Wallet.Read(ctx, req, doc.NameSpaceGlobal, doc.CryptoSECP256K1, payload.Address)
	if err != nil {
		return nil, err
	}

	newWallet, err := s.Wallet.Update(ctx, req, payload)
	if err != nil {
		return nil, err
	}

	srcPath := fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.CryptoSECP256K1, payload.Address)

	// 更新所有的alias
	for chain, addr := range payload.AddressAlias {
		// 如果地址有变化，需要完全创建新的alias
		tmpNs := oldWallet.NameSpaces
		if oldWallet.AddressAlias[chain] != addr {
			tmpNs = make([]string, 0)
		}
		err = s.Alias.Update(ctx, req, aliasType(chain),
			addr, srcPath, tmpNs, newWallet.NameSpaces)
		if err != nil {
			return nil, err
		}
	}

	// 移除已经不存在的alias
	for chain, oAddr := range oldWallet.AddressAlias {
		if nAddr, ok := newWallet.AddressAlias[chain]; !ok || nAddr != oAddr {
			err = s.Alias.Delete(ctx, req, aliasType(chain), oAddr, oldWallet.NameSpaces)
			if err != nil {
				return nil, err
			}
		}
	}

	return newWallet, nil
}

func (s *storageEx) delete(ctx context.Context, req *logical.Request, address string) error {
	// 获取目标钱包
	oldWallet, err := s.Wallet.Read(ctx, req, doc.NameSpaceGlobal, doc.CryptoSECP256K1, address)
	if err != nil {
		return err
	}

	err = s.Wallet.Delete(ctx, req, doc.CryptoSECP256K1, oldWallet.Address)
	if err != nil {
		return err
	}

	for chain, addr := range oldWallet.AddressAlias {
		err = s.Alias.Delete(ctx, req, aliasType(chain), addr, oldWallet.NameSpaces)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *storageEx) read(ctx context.Context, req *logical.Request, ns, chain, address string) (*modules.WalletExtra, error) {
	path, _ := s.Alias.ReadSrcPath(ctx, req, ns, aliasType(chain), address)
	w, err := s.Wallet.ReadFromPath(ctx, req, path)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *storageEx) list(ctx context.Context, req *logical.Request, ns, chain string) ([]string, error) {
	if ns == doc.NameSpaceGlobal && chain == "" {
		return s.Wallet.List(ctx, req, doc.CryptoSECP256K1)
	}
	return s.Alias.List(ctx, req, ns, aliasType(chain))
}
