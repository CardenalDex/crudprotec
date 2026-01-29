package usecase

import (
	"context"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
)

type BusinessRepository interface {
	Create(ctx context.Context, b *entity.Business) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Business, error)
	Update(ctx context.Context, b *entity.Business) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MerchantRepository interface {
	Create(ctx context.Context, m *entity.Merchant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Merchant, error)
	GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]entity.Merchant, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TransactionRepository interface {
	Create(ctx context.Context, t *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	ListByMerchant(ctx context.Context, merchantID uuid.UUID) ([]entity.Transaction, error)
	GetAll(ctx context.Context) ([]entity.Transaction, error)
}

type LogRepository interface {
	Create(ctx context.Context, l *entity.Log) error
	GetByID(ctx context.Context, logID string) (entity.Log, error)
	GetByResource(ctx context.Context, resourceID string) ([]entity.Log, error)
	GetAll(ctx context.Context) ([]entity.Log, error)
}

// Input Ports

type TransactionUseCase interface {
	ProcessTransaction(ctx context.Context, merchantID uuid.UUID, amount int64) (*entity.Transaction, error)
}

type AdminUseCase interface {
	RegisterBusiness(ctx context.Context, commission int64) (*entity.Business, error)
	GetAuditTrail(ctx context.Context, resourceID string) ([]entity.Log, error)
}
