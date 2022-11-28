package storage

import "errors"

var (
	ErrWalletNotFound  = errors.New("wallet not found")
	ErrCryptoTypeIsNil = errors.New("crypto type is nil")
)
