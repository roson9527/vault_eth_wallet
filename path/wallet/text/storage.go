package text

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (s *storageEx) put(ctx context.Context, req *logical.Request, payload *modules.WalletExtra) (*modules.WalletExtra, error) {
	if !(payload.Address != "") {
		return nil, errors.New("default address is nil")
	}

	// 获取目标钱包
	oldWallet, err := s.Wallet.Read(ctx, req, doc.NameSpaceGlobal, doc.CryptoTEXT, payload.Address)
	if oldWallet == nil {
		oldWallet = &modules.WalletExtra{
			AddressAlias: make(map[string]string),
		}
	}

	err = s.Wallet.Put(ctx, req, payload)
	if err != nil {
		return nil, err
	}

	srcPath := fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.CryptoTEXT, payload.Address)

	// 更新所有的alias
	for chain, addr := range payload.AddressAlias {
		// 如果地址有变化，需要完全创建新的alias
		tmpNs := oldWallet.NameSpaces
		if oldWallet.AddressAlias[chain] != addr {
			tmpNs = make([]string, 0)
		}
		err = s.Alias.Update(ctx, req, aliasType(chain),
			addr, srcPath, tmpNs, payload.NameSpaces)
		if err != nil {
			return nil, err
		}
	}

	// 移除已经不存在的alias
	for chain, oAddr := range oldWallet.AddressAlias {
		if nAddr, ok := payload.AddressAlias[chain]; !ok || nAddr != oAddr {
			err = s.Alias.Delete(ctx, req, aliasType(chain), oAddr, oldWallet.NameSpaces)
			if err != nil {
				return nil, err
			}
		}
	}

	return payload, nil
}

func (s *storageEx) read(ctx context.Context, req *logical.Request, ns, chain, address string) (*modules.WalletExtra, error) {
	path, _ := s.Alias.ReadSrcPath(ctx, req, ns, aliasType(chain), address)
	w, err := s.Wallet.ReadFromPath(ctx, req, path)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *storageEx) delete(ctx context.Context, req *logical.Request, address string) error {
	// 获取目标钱包
	oldWallet, err := s.Wallet.Read(ctx, req, doc.NameSpaceGlobal, doc.CryptoTEXT, address)
	if err != nil {
		return err
	}

	err = s.Wallet.Delete(ctx, req, doc.CryptoTEXT, oldWallet.Address)
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

func (s *storageEx) list(ctx context.Context, req *logical.Request, ns, chain string) ([]string, error) {
	if ns == doc.NameSpaceGlobal && chain == "" {
		return s.Wallet.List(ctx, req, doc.CryptoTEXT)
	}
	return s.Alias.List(ctx, req, ns, aliasType(chain))
}
