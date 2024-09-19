package models

import "github.com/google/uuid"

type TransactionDetails struct {
	TargetAccountUUID uuid.UUID
	Amount            int64
}
