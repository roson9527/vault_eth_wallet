package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"
)

type walletMock struct {
}

func (wm *walletMock) Read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	fmt.Println("hello world")
	return nil, errors.New("fx")
}

func TestAutoBind(t *testing.T) {
	wm := walletMock{}
	out := autoBind(&wm)
	f := out[logical.ReadOperation].(*framework.PathOperation)
	fmt.Println(f.Callback(context.Background(), nil, nil))
}
