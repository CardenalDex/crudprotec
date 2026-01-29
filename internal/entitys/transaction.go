package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID
	MerchantID uuid.UUID
	Amount     int64 // Value in cents (200.00 -> 20000)
	Commission int64 // in cents
	Fee        int64 // in centi% 5.5=550
	Timestamp  time.Time
	DeletedAt  *time.Time
}
