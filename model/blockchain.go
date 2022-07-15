package model

import (
	"time"
)

type UserAddress struct {
	ID        uint64
	Address   string
	CreatedAt time.Time
}
