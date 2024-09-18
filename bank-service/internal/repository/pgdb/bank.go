package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/d1mitrii/money-transfer/bank-service/internal/models"
	"github.com/d1mitrii/money-transfer/bank-service/internal/repository/repoerr"
	"github.com/d1mitrii/money-transfer/bank-service/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
)

type BankRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *BankRepo {
	return &BankRepo{pg}
}

func (b *BankRepo) CreateAccount(ctx context.Context, account models.Account) (uuid.UUID, error) {
	const op = "BankRepo.CreateAccount"

	sql := `INSERT INTO accounts(account_name, balance) VALUES ($1, $2) RETURNING uuid;`

	var accountUUID uuid.UUID

	err := b.Pool.QueryRow(ctx, sql, account.Name, account.Balance).Scan(&accountUUID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return uuid.Nil, repoerr.ErrAlreadyExist
			}
		}
		return uuid.Nil, fmt.Errorf("%s - b.Pool.QueryRow: %w", op, err)
	}

	return accountUUID, nil
}

func (b *BankRepo) GetAccount(ctx context.Context, accountUUID uuid.UUID) (models.Account, error) {
	const op = "BankRepo.GetAccount"

	sql := `SELECT (uuid, account_name, balance, created_at, updated_at) FROM accounts WHERE uuid = $1;`

	var account models.Account

	err := b.Pool.QueryRow(ctx, sql, accountUUID).Scan(&account)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Account{}, repoerr.ErrNotFound
		}
		return models.Account{}, fmt.Errorf("%s - b.Pool.QueryRow: %w", op, err)
	}

	return account, nil
}

func (b *BankRepo) DeleteAccount(ctx context.Context, accountUUID uuid.UUID) error {
	const op = "BankRepo.DeleteAccount"

	sql := `DELETE FROM accounts WHERE uuid = $1;`

	_, err := b.Pool.Exec(ctx, sql, accountUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerr.ErrNotFound
		}
		return fmt.Errorf("%s - b.Pool.Exec: %w", op, err)
	}

	return nil
}

func (b *BankRepo) Deposit(ctx context.Context, details models.TransactionDetails) error {
	const op = "BankRepo.Deposit"

	sql := `UPDATE accounts SET balance = balance + $1 WHERE uuid = $2;`

	tag, err := b.Pool.Exec(ctx, sql, details.Amount, details.TargetAccountUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerr.ErrNotFound
		}
		return fmt.Errorf("%s - b.Pool.Exec: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrNotFound
	}

	return nil
}

func (b *BankRepo) Withdraw(ctx context.Context, details models.TransactionDetails) error {
	const op = "BankRepo.Withdraw"

	sql := `UPDATE accounts SET balance = balance - $1 WHERE uuid = $2;`

	tag, err := b.Pool.Exec(ctx, sql, details.Amount, details.TargetAccountUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerr.ErrNotFound
		}
		return fmt.Errorf("%s - b.Pool.Exec: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrNotFound
	}

	return nil
}

func (b *BankRepo) Refund(ctx context.Context, details models.TransactionDetails) error {
	const op = "BankRepo.Refund"

	sql := `UPDATE accounts SET balance = balance + $1 WHERE uuid = $2;`

	tag, err := b.Pool.Exec(ctx, sql, details.Amount, details.TargetAccountUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerr.ErrNotFound
		}
		return fmt.Errorf("%s - b.Pool.Exec: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrNotFound
	}

	return nil
}
