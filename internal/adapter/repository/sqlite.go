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

func (r *sqliteRepo) Create(ctx context.Context, b *entity.Business) error {
	model := toBusinessModel(b)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *sqliteRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Business, error) {
	var model BusinessModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *sqliteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&BusinessModel{}, "id = ?", id).Error
}

func (r *sqliteRepo) CreateTransaction(ctx context.Context, t *entity.Transaction) error {
	model := toTransactionModel(t)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *sqliteRepo) ListByMerchant(ctx context.Context, mID uuid.UUID) ([]entity.Transaction, error) {
	var models []TransactionModel
	if err := r.db.WithContext(ctx).Where("merchant_id = ?", mID).Find(&models).Error; err != nil {
		return nil, err
	}

	results := make([]entity.Transaction, len(models))
	for i, m := range models {
		results[i] = *m.toEntity()
	}
	return results, nil
}

func (r *sqliteRepo) CreateLog(ctx context.Context, l *entity.Log) error {
	model := toLogModel(l)
	return r.db.WithContext(ctx).Create(model).Error
}
