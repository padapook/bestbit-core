package model

import (
	"github.com/padapook/bestbit-core/internal/utils/apperr"
)

var (
	ErrWalletNotFound = &apperr.AppErr{
		HTTPCode: 404,
		BusinessCode: 40401,
		Message: "ERR_WALLET_NOT_FOUND",
	}

	ErrInSufficientFunds = &apperr.AppErr{
		HTTPCode: 400,
		BusinessCode: 40001,
		Message: "ERR_INSUFFICIEN_FUNDS",
	}
)

