package entity

import (
	"time"

	"github.com/google/uuid"
)

type Business struct {
	ID         uuid.UUID
	Commission int64 // Represented in basis points (e.g., 550 = 5.5%)
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type Merchant struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
