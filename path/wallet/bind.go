package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/utils"
	"reflect"
)

var (
	bindSlice = []logical.Operation{
		logical.CreateOperation,
		logical.ReadOperation,
		logical.UpdateOperation,
		logical.PatchOperation,
		logical.DeleteOperation,
		logical.ListOperation,
		logical.HelpOperation,
		// The operations below are called globally, the path is less relevant.
		logical.RevokeOperation,
		logical.RenewOperation,
		logical.RollbackOperation,
	}
)

func autoBind(self any) map[logical.Operation]framework.OperationHandler {
	out := make(map[logical.Operation]framework.OperationHandler)

	typ := reflect.TypeOf(self)
	if typ.Kind() != reflect.Pointer {
		panic("self must be a pointer")
	}

	value := reflect.ValueOf(self)
	for _, l := range bindSlice {
		name := utils.FirstUpper(string(l))
		_, existed := typ.MethodByName(name)
		if !existed {
			continue
		}
		ff := (value.MethodByName(name).Interface()).(func(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error))
		out[l] = &framework.PathOperation{
			Callback: ff,
		}
	}

	return out
}
