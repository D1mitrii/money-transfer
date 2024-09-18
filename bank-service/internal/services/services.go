package services

import (
	"github.com/d1mitrii/money-transfer/bank-service/internal/services/bank"
)

type Services struct {
	Bank *bank.Bank
}

func New(bank *bank.Bank) *Services {
	return &Services{
		Bank: bank,
	}
}
