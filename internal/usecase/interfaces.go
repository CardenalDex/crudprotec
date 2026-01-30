package usecase

import (
	"context"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
)

type BusinessRepository interface {
	CreateBusiness(ctx context.Context, b *entity.Business) error
	GetBusinessByID(ctx context.Context, id uuid.UUID) (*entity.Business, error)
	UpdateBusiness(ctx context.Context, b *entity.Business) error
	DeleteBusiness(ctx context.Context, id uuid.UUID) error
}

type MerchantRepository interface {
	CreateMerchant(ctx context.Context, m *entity.Merchant) error
	GetMerchantByID(ctx context.Context, id uuid.UUID) (*entity.Merchant, error)
	GetMerchantByBusinessID(ctx context.Context, businessID uuid.UUID) ([]entity.Merchant, error)
	DeleteMerchant(ctx context.Context, id uuid.UUID) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, t *entity.Transaction) error
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	TransactionListByMerchant(ctx context.Context, merchantID uuid.UUID) ([]entity.Transaction, error)
	GetAllTransaction(ctx context.Context) ([]entity.Transaction, error)
}

type LogRepository interface {
	CreateLog(ctx context.Context, l *entity.Log) error
	GetLogByID(ctx context.Context, logID string) (entity.Log, error)
	GetLogByResource(ctx context.Context, resourceID string) ([]entity.Log, error)
	GetAll(ctx context.Context) ([]entity.Log, error)
}

// /////////////////////////////////////////////////////
type TransactionUseCase interface {
	ProcessTransaction(ctx context.Context, actor string, merchantID uuid.UUID, amount int64) (*entity.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetMerchantTransactions(ctx context.Context, merchantID uuid.UUID) ([]entity.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]entity.Transaction, error)
}

type MerchantUseCase interface {
	RegisterMerchant(ctx context.Context, actor string, businessID uuid.UUID) (*entity.Merchant, error)
	GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error)
	GetBusinessMerchants(ctx context.Context, businessID uuid.UUID) ([]entity.Merchant, error)
	RemoveMerchant(ctx context.Context, actor string, id uuid.UUID) error
}

type AdminUseCase interface {
	RegisterBusiness(ctx context.Context, actor string, commission int64) (*entity.Business, error)
	GetBusiness(ctx context.Context, id uuid.UUID) (*entity.Business, error)
	UpdateBusinessCommission(ctx context.Context, actor string, id uuid.UUID, newCommission int64) (*entity.Business, error)
	RemoveBusiness(ctx context.Context, actor string, id uuid.UUID) error

	// Auditing
	GetAuditTrail(ctx context.Context, resourceID string) ([]entity.Log, error)
	GetLogDetails(ctx context.Context, logID string) (entity.Log, error)
	GetAllLogs(ctx context.Context) ([]entity.Log, error)
}
