package bank

import (
	"context"
	"log/slog"

	"github.com/d1mitrii/money-transfer/bank-service/internal/models"
)

type (
	AccountProvider interface {
		CreateAccount(ctx context.Context, account models.Account) (string, error)
		GetAccount(ctx context.Context, accountID string) (models.Account, error)
		DeleteAccount(ctx context.Context, accountID string) error
	}

	BalanceProvider interface {
		Deposit(ctx context.Context, accountID string) error
		Withdraw(ctx context.Context, accountID string) error
		Refund(ctx context.Context, accountID string) error
	}

	Bank struct {
		log             *slog.Logger
		accountProvider AccountProvider
		balanceProvider BalanceProvider
	}
)

func New(
	log *slog.Logger,
	accountProvider AccountProvider,
	balanceProvider BalanceProvider,
) *Bank {
	return &Bank{
		log:             log,
		accountProvider: accountProvider,
		balanceProvider: balanceProvider,
	}
}
