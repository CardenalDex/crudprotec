package repository

import (
	"context"
	"os"
	"path/filepath"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitInternalDB() *gorm.DB {
	dbPath := "/app/data/local.db"

	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		panic("failed to create internal database directory: " + err.Error())
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect to internal database: " + err.Error())
	}

	db.AutoMigrate(
		&BusinessModel{},
		&MerchantModel{},
		&TransactionModel{},
		&LogModel{},
	)

	return db
}

type sqliteRepo struct {
	db *gorm.DB
}

func NewSQLiteRepository(db *gorm.DB) *sqliteRepo {
	return &sqliteRepo{db: db}
}

// --- BusinessRepository Implementation ---

func (r *sqliteRepo) CreateBusiness(ctx context.Context, b *entity.Business) error {
	model := toBusinessModel(b)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *sqliteRepo) GetBusinessByID(ctx context.Context, id uuid.UUID) (*entity.Business, error) {
	var model BusinessModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *sqliteRepo) UpdateBusiness(ctx context.Context, b *entity.Business) error {
	model := toBusinessModel(b)

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *sqliteRepo) DeleteBusiness(ctx context.Context, id uuid.UUID) error {

	return r.db.WithContext(ctx).Delete(&BusinessModel{}, "id = ?", id).Error
}

// --- MerchantRepository Implementation ---

func (r *sqliteRepo) CreateMerchant(ctx context.Context, m *entity.Merchant) error {
	model := toMerchantModel(m)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *sqliteRepo) GetMerchantByID(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	var model MerchantModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *sqliteRepo) GetMerchantByBusinessID(ctx context.Context, bizID uuid.UUID) ([]entity.Merchant, error) {
	var models []MerchantModel
	if err := r.db.WithContext(ctx).Where("business_id = ?", bizID).Find(&models).Error; err != nil {
		return nil, err
	}

	merchants := make([]entity.Merchant, len(models))
	for i, m := range models {
		merchants[i] = *m.toEntity()
	}
	return merchants, nil
}

func (r *sqliteRepo) DeleteMerchant(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&MerchantModel{}, "id = ?", id).Error
}

// --- TransactionRepository Implementation ---

func (r *sqliteRepo) CreateTransaction(ctx context.Context, t *entity.Transaction) error {
	model := toTransactionModel(t)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *sqliteRepo) GetTransactionByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	var model TransactionModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *sqliteRepo) TransactionListByMerchant(ctx context.Context, mID uuid.UUID) ([]entity.Transaction, error) {
	var models []TransactionModel
	if err := r.db.WithContext(ctx).Where("merchant_id = ?", mID).Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]entity.Transaction, len(models))
	for i, m := range models {
		transactions[i] = *m.toEntity()
	}
	return transactions, nil
}

func (r *sqliteRepo) GetAllTransaction(ctx context.Context) ([]entity.Transaction, error) {
	var models []TransactionModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]entity.Transaction, len(models))
	for i, m := range models {
		transactions[i] = *m.toEntity()
	}
	return transactions, nil
}

// --- LogRepository Implementation ---

func (r *sqliteRepo) CreateLog(ctx context.Context, l *entity.Log) error {
	model := toLogModel(l)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *sqliteRepo) GetLogByID(ctx context.Context, id string) (entity.Log, error) {
	var model LogModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return entity.Log{}, err
	}
	return *model.toEntity(), nil
}

func (r *sqliteRepo) GetLogByResource(ctx context.Context, resID string) ([]entity.Log, error) {
	var models []LogModel
	if err := r.db.WithContext(ctx).Where("resource_id = ?", resID).Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]entity.Log, len(models))
	for i, m := range models {
		logs[i] = *m.toEntity()
	}
	return logs, nil
}

func (r *sqliteRepo) GetAll(ctx context.Context) ([]entity.Log, error) {
	var models []LogModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]entity.Log, len(models))
	for i, m := range models {
		logs[i] = *m.toEntity()
	}
	return logs, nil
}
