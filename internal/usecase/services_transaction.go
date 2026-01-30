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

func (s *transactionService) ProcessTransaction(ctx context.Context, actor string, mID uuid.UUID, amount int64) (*entity.Transaction, error) {

	merchant, err := s.merchantRepo.GetMerchantByID(ctx, mID)
	if err != nil {
		return nil, errors.New("merchant not found")
	}

	biz, err := s.bizRepo.GetBusinessByID(ctx, merchant.BusinessID)
	if err != nil {
		return nil, errors.New("business configuration missing")
	}

	// If biz.Commission is in basis points (e.g., 550 for 5.5%)
	// Commission = (Amount * 550) / 10000
	commission := (amount * biz.Commission) / 10000

	tx := &entity.Transaction{
		ID:         uuid.New(),
		MerchantID: mID,
		Amount:     amount,
		Commission: biz.Commission,
		Fee:        commission,
		Timestamp:  time.Now(),
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "TRANSACTION_CREATED",
		Actor:      actor,
		ResourceID: tx.ID.String(),
		Timestamp:  time.Now(),
	})

	return tx, nil
}

func (s *transactionService) GetTransaction(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	tx, err := s.repo.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return tx, nil
}

func (s *transactionService) GetMerchantTransactions(ctx context.Context, merchantID uuid.UUID) ([]entity.Transaction, error) {

	if _, err := s.merchantRepo.GetMerchantByID(ctx, merchantID); err != nil {
		return nil, errors.New("merchant not found")
	}

	return s.repo.TransactionListByMerchant(ctx, merchantID)
}

func (s *transactionService) GetAllTransactions(ctx context.Context) ([]entity.Transaction, error) {
	return s.repo.GetAllTransaction(ctx)
}
