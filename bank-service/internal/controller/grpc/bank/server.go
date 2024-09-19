package bankgrpc

import (
	"context"
	"errors"

	"github.com/d1mitrii/money-transfer/bank-service/internal/controller/grpc/grpcerr"
	"github.com/d1mitrii/money-transfer/bank-service/internal/models"
	"github.com/d1mitrii/money-transfer/bank-service/internal/services/servicerr"
	bankv1 "github.com/d1mitrii/money-transfer/bank-service/pkg/grpc/bank/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (b *bankAPI) CreateAccount(ctx context.Context, in *bankv1.CreateAccountRequest) (*bankv1.CreateAccountResponse, error) {
	if len(in.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty account name")
	}

	if in.GetBalance() < 0 {
		return nil, status.Error(codes.InvalidArgument, "negative balance forbidden")
	}

	accountUUID, err := b.bank.CreateAccount(ctx, models.Account{
		Name:    in.GetName(),
		Balance: in.GetBalance(),
	})

	if err != nil {
		if errors.Is(err, servicerr.ErrAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, "account already exist")
		}
		return nil, grpcerr.ErrServiceLayer
	}

	return &bankv1.CreateAccountResponse{AccountUUID: accountUUID.String()}, nil
}

func (b *bankAPI) GetAccount(ctx context.Context, in *bankv1.GetAccountRequest) (*bankv1.GetAccountResponse, error) {
	accountUUID, err := uuid.Parse(in.GetAccountUUID())
	if err != nil {
		return nil, grpcerr.ErrParseUUID
	}

	account, err := b.bank.GetAccount(ctx, accountUUID)
	if err != nil {
		if errors.Is(err, servicerr.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, grpcerr.ErrServiceLayer
	}

	return &bankv1.GetAccountResponse{
		AccountUUID: account.UUID.String(),
		Name:        account.Name,
		Balance:     account.Balance,
	}, nil
}

func (b *bankAPI) DeleteAccount(ctx context.Context, in *bankv1.DeleteAccountRequest) (*emptypb.Empty, error) {
	accountUUID, err := uuid.Parse(in.GetAccountUUID())
	if err != nil {
		return &emptypb.Empty{}, grpcerr.ErrParseUUID
	}

	if err := b.bank.DeleteAccount(ctx, accountUUID); err != nil {
		if errors.Is(err, servicerr.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, "account not found")
		}
		return &emptypb.Empty{}, grpcerr.ErrServiceLayer
	}

	return &emptypb.Empty{}, nil
}

func (b *bankAPI) Deposit(ctx context.Context, in *bankv1.DepositRequest) (*emptypb.Empty, error) {
	accountUUID, err := uuid.Parse(in.GetAccountUUID())
	if err != nil {
		return &emptypb.Empty{}, grpcerr.ErrParseUUID
	}

	if in.GetAmount() <= 0 {
		return &emptypb.Empty{}, grpcerr.ErrIncorrectAmount
	}

	if err := b.bank.Deposit(ctx, models.TransactionDetails{
		TargetAccountUUID: accountUUID,
		Amount:            in.GetAmount(),
	}); err != nil {
		if errors.Is(err, servicerr.ErrNotFound) {
			return &emptypb.Empty{}, grpcerr.ErrAccountNotFound
		}
		return &emptypb.Empty{}, grpcerr.ErrServiceLayer
	}

	return &emptypb.Empty{}, nil
}

func (b *bankAPI) Withdraw(ctx context.Context, in *bankv1.WithdrawRequest) (*emptypb.Empty, error) {
	accountUUID, err := uuid.Parse(in.GetAccountUUID())
	if err != nil {
		return &emptypb.Empty{}, grpcerr.ErrParseUUID
	}

	if in.GetAmount() <= 0 {
		return &emptypb.Empty{}, grpcerr.ErrIncorrectAmount
	}

	if err := b.bank.Withdraw(ctx, models.TransactionDetails{
		TargetAccountUUID: accountUUID,
		Amount:            in.GetAmount(),
	}); err != nil {
		if errors.Is(err, servicerr.ErrNotFound) {
			return &emptypb.Empty{}, grpcerr.ErrAccountNotFound
		}
		return &emptypb.Empty{}, grpcerr.ErrServiceLayer
	}

	return &emptypb.Empty{}, nil
}

func (b *bankAPI) Refund(ctx context.Context, in *bankv1.RefundRequest) (*emptypb.Empty, error) {
	accountUUID, err := uuid.Parse(in.GetAccountUUID())
	if err != nil {
		return &emptypb.Empty{}, grpcerr.ErrParseUUID
	}

	if in.GetAmount() <= 0 {
		return &emptypb.Empty{}, grpcerr.ErrIncorrectAmount
	}

	if err := b.bank.Refund(ctx, models.TransactionDetails{
		TargetAccountUUID: accountUUID,
		Amount:            in.GetAmount(),
	}); err != nil {
		if errors.Is(err, servicerr.ErrNotFound) {
			return &emptypb.Empty{}, grpcerr.ErrAccountNotFound
		}
		return &emptypb.Empty{}, grpcerr.ErrServiceLayer
	}

	return &emptypb.Empty{}, nil
}
