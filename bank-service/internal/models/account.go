package models

import "time"

type Account struct {
	ID        string
	Name      string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
