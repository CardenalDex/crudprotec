package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
)

type transactionService struct {
	repo         TransactionRepository
	merchantRepo MerchantRepository
	bizRepo      BusinessRepository
	logRepo      LogRepository
}

func NewTransactionService(tr TransactionRepository, mr MerchantRepository, br BusinessRepository, lr LogRepository) TransactionUseCase {
	return &transactionService{tr, mr, br, lr}
}

func (s *transactionService) ProcessTransaction(ctx context.Context, mID uuid.UUID, amount int64) (*entity.Transaction, error) {
	// 1. Validate Merchant
	merchant, err := s.merchantRepo.GetMerchantByID(ctx, mID)
	if err != nil {
		return nil, errors.New("merchant not found")
	}

	// 2. Get Business to find commission rate
	biz, err := s.bizRepo.GetBusinessByID(ctx, merchant.BusinessID)
	if err != nil {
		return nil, errors.New("business configuration missing")
	}

	// 3. Calculate Commission (Cents Math)
	// If biz.Commission is in basis points (e.g., 550 for 5.5%)
	// Commission = (Amount * 550) / 10000
	commission := (amount * biz.Commission) / 10000

	// Assume a flat fee of $0.50 (50 cents) for this example
	var flatFee int64 = 50

	// 4. Create Transaction Entity
	tx := &entity.Transaction{
		ID:         uuid.New(),
		MerchantID: mID,
		Amount:     amount,
		Commission: commission,
		Fee:        flatFee,
		Timestamp:  time.Now(),
	}

	// 5. Persist
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	// 6. Audit Log (Async or Sync depending on requirements)
	_ = s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "TRANSACTION_CREATED",
		Actor:      "system",
		ResourceID: tx.ID.String(),
		Timestamp:  time.Now(),
	})

	return tx, nil
}

// GetTransaction retrieves a single transaction by its UUID
func (s *transactionService) GetTransaction(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	tx, err := s.repo.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return tx, nil
}

// GetMerchantTransactions retrieves all transactions associated with a specific merchant
func (s *transactionService) GetMerchantTransactions(ctx context.Context, merchantID uuid.UUID) ([]entity.Transaction, error) {
	// Optional: Validate merchant existence first if strict referential integrity is needed
	if _, err := s.merchantRepo.GetMerchantByID(ctx, merchantID); err != nil {
		return nil, errors.New("merchant not found")
	}

	return s.repo.TransactionListByMerchant(ctx, merchantID)
}

// GetAllTransactions retrieves every transaction in the database
func (s *transactionService) GetAllTransactions(ctx context.Context) ([]entity.Transaction, error) {
	return s.repo.GetAllTransaction(ctx)
}
