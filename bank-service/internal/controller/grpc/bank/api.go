package bankgrpc

import (
	"context"

	"github.com/d1mitrii/money-transfer/bank-service/internal/models"
	bankv1 "github.com/d1mitrii/money-transfer/bank-service/pkg/grpc/bank/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Bank interface {
	CreateAccount(ctx context.Context, account models.Account) (uuid.UUID, error)
	GetAccount(ctx context.Context, accountUUID uuid.UUID) (models.Account, error)
	DeleteAccount(ctx context.Context, accountUUID uuid.UUID) error
	Deposit(ctx context.Context, details models.TransactionDetails) error
	Withdraw(ctx context.Context, details models.TransactionDetails) error
	Refund(ctx context.Context, details models.TransactionDetails) error
}

type bankAPI struct {
	bankv1.UnimplementedBankServer
	bank Bank
}

func Register(server *grpc.Server, bank Bank) {
	bankv1.RegisterBankServer(server, &bankAPI{bank: bank})
}
