package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	UUID      uuid.UUID `db:"uuid"`
	Name      string    `db:"account_name"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
