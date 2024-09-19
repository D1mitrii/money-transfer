package bank

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/d1mitrii/money-transfer/bank-service/internal/models"
	"github.com/d1mitrii/money-transfer/bank-service/internal/repository/repoerr"
	"github.com/d1mitrii/money-transfer/bank-service/internal/services/servicerr"
	"github.com/google/uuid"
)

type (
	AccountProvider interface {
		CreateAccount(ctx context.Context, account models.Account) (uuid.UUID, error)
		GetAccount(ctx context.Context, accountUUID uuid.UUID) (models.Account, error)
		DeleteAccount(ctx context.Context, accountUUID uuid.UUID) error
	}

	BalanceProvider interface {
		Deposit(ctx context.Context, details models.TransactionDetails) error
		Withdraw(ctx context.Context, details models.TransactionDetails) error
		Refund(ctx context.Context, details models.TransactionDetails) error
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

func (b *Bank) CreateAccount(ctx context.Context, account models.Account) (uuid.UUID, error) {
	const op = "Bank.CreateAccount"
	log := b.log.With(
		slog.String("op", op),
		slog.String("account name", account.Name),
	)
	if account.Balance < 0 {
		return uuid.Nil, servicerr.ErrInvalidArgument
	}

	id, err := b.accountProvider.CreateAccount(ctx, account)
	if err != nil {
		if errors.Is(err, repoerr.ErrAlreadyExist) {
			log.Error("account already exist", slog.Any("err", err))
			return uuid.Nil, servicerr.ErrAlreadyExist
		}
		log.Error("failed to create account", slog.Any("err", err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (b *Bank) GetAccount(ctx context.Context, accountUUID uuid.UUID) (models.Account, error) {
	const op = "Bank.GetAccount"
	log := b.log.With(
		slog.String("op", op),
		slog.String("accountUUID", accountUUID.String()),
	)

	account, err := b.accountProvider.GetAccount(ctx, accountUUID)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			log.Error("account not found", slog.Any("err", err))
		}
		log.Error("failed to get account", slog.Any("err", err))
		return models.Account{}, fmt.Errorf("%s: %w", op, err)
	}
	return account, nil
}

func (b *Bank) DeleteAccount(ctx context.Context, accountUUID uuid.UUID) error {
	const op = "Bank.GetAccount"
	log := b.log.With(
		slog.String("op", op),
		slog.String("accountUUID", accountUUID.String()),
	)

	err := b.accountProvider.DeleteAccount(ctx, accountUUID)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			log.Error("account not found", slog.Any("err", err))
			return servicerr.ErrNotFound
		}
		log.Error("failed to delete account", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (b *Bank) Deposit(ctx context.Context, details models.TransactionDetails) error {
	const op = "Bank.Deposit"
	log := b.log.With(
		slog.String("op", op),
		slog.String("accountUUID", details.TargetAccountUUID.String()),
	)

	if details.Amount < 0 {
		log.Error("incorrect ammount")
		return servicerr.ErrInvalidArgument
	}

	if err := b.balanceProvider.Deposit(ctx, details); err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			log.Error("account not found", slog.Any("err", err))
			return servicerr.ErrNotFound
		}
		log.Error("deposit failed", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (b *Bank) Withdraw(ctx context.Context, details models.TransactionDetails) error {
	const op = "Bank.Withdraw"
	log := b.log.With(
		slog.String("op", op),
		slog.String("accountUUID", details.TargetAccountUUID.String()),
	)

	if details.Amount < 0 {
		log.Error("incorrect ammount")
		return servicerr.ErrInvalidArgument
	}

	if err := b.balanceProvider.Withdraw(ctx, details); err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			log.Error("account not found", slog.Any("err", err))
			return servicerr.ErrNotFound
		}
		log.Error("withdraw failed", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (b *Bank) Refund(ctx context.Context, details models.TransactionDetails) error {
	const op = "Bank.Refund"
	log := b.log.With(
		slog.String("op", op),
		slog.String("accountUUID", details.TargetAccountUUID.String()),
	)

	if details.Amount < 0 {
		log.Error("incorrect ammount")
		return servicerr.ErrInvalidArgument
	}

	if err := b.balanceProvider.Refund(ctx, details); err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			log.Error("account not found", slog.Any("err", err))
			return servicerr.ErrNotFound
		}
		log.Error("refund failed", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
