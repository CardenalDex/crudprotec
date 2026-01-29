package repository

import (
	"time"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Commission int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (BusinessModel) TableName() string { return "businesses" }

type MerchantModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	BusinessID uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (MerchantModel) TableName() string { return "merchants" }

type TransactionModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	MerchantID uuid.UUID `gorm:"type:uuid;index"`
	Amount     int64     // Stored in cents
	Commission int64     // Calculated cents
	Fee        int64
	Timestamp  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (TransactionModel) TableName() string { return "transactions" }

type LogModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Action         string
	Actor          string
	ResourceID     string `gorm:"index"`
	PrevResourceID string
	Timestamp      time.Time `gorm:"index"`
}

func (LogModel) TableName() string { return "audit_logs" }

//////////////////////////////////////////////////////////////////////////

func toBusinessModel(e *entity.Business) *BusinessModel {
	return &BusinessModel{
		ID:         e.ID,
		Commission: e.Commission,
	}
}

func (m *BusinessModel) toEntity() *entity.Business {
	return &entity.Business{
		ID:         m.ID,
		Commission: m.Commission,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func toMerchantModel(e *entity.Merchant) *MerchantModel {
	return &MerchantModel{
		ID:         e.ID,
		BusinessID: e.BusinessID,
	}
}

func (m *MerchantModel) toEntity() *entity.Merchant {
	return &entity.Merchant{
		ID:         m.ID,
		BusinessID: m.BusinessID,
		CreatedAt:  m.CreatedAt,
		DeletedAt:  &m.DeletedAt.Time,
	}
}

func toTransactionModel(e *entity.Transaction) *TransactionModel {
	return &TransactionModel{
		ID:         e.ID,
		MerchantID: e.MerchantID,
		Amount:     e.Amount,
		Commission: e.Commission,
		Fee:        e.Fee,
		Timestamp:  e.Timestamp,
	}
}

func (m *TransactionModel) toEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:         m.ID,
		MerchantID: m.MerchantID,
		Amount:     m.Amount,
		Commission: m.Commission,
		Fee:        m.Fee,
		Timestamp:  m.Timestamp,
	}
}

func toLogModel(e *entity.Log) *LogModel {
	return &LogModel{
		ID:             e.ID,
		Action:         e.Action,
		Actor:          e.Actor,
		ResourceID:     e.ResourceID,
		PrevResourceID: e.PrevResourceID,
		Timestamp:      e.Timestamp,
	}
}
