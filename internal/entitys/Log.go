package entity

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID             uuid.UUID
	Action         string //"UPDATE_TRANSACTION"
	Actor          string // made by
	ResourceID     string // ID of the affected resource
	PrevResourceID string // Previous state
	Timestamp      time.Time
}
