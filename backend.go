// Copyright © 2018 Immutability, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/policy"
	"github.com/roson9527/vault_eth_wallet/path/socialid"
	"github.com/roson9527/vault_eth_wallet/path/wallet"
)

// Factory returns the backend
// 入口，插件的返回入口
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := NewBackend()
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, fmt.Errorf("confiuration passed into backend is nil")
	}
	if err = b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func NewBackend() (*modules.EthWalletBackend, error) {
	b := new(modules.EthWalletBackend)

	// TODO
	b.Backend = &framework.Backend{
		Help: "",
		// 响应路由
		Paths: framework.PathAppend(
			wallet.Path(),   // wallet path
			policy.Path(),   // policy path
			socialid.Path(), // socialId path
		),
		// 特殊带权限路径，不能正则，但是可以通配符
		PathsSpecial: &logical.Paths{
			Root:            nil,
			Unauthenticated: nil,
			LocalStorage:    nil,
			SealWrapStorage: nil,
		},
		// 秘密类型列表，简化回调
		Secrets:        []*framework.Secret{},
		RunningVersion: "v0.3.0",
		//// 初始化方法位置
		//InitializeFunc:    nil,
		//// 定时器回调
		//PeriodicFunc:      nil,
		//// minimum age of a WAL
		//WALRollback:       nil,
		//WALRollbackMinAge: 0,
		//// 清理方法
		//Clean:             nil,
		//// 修改键后的调用
		//Invalidate:        nil,
		//// 身份相关
		//AuthRenew:         nil,
		// 用于后端实现的逻辑 secret engine
		BackendType: logical.TypeLogical,
	}

	return b, nil
}
