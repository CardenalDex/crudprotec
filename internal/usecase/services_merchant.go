package usecase

import (
	"context"
	"errors"
	"time"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
)

type merchantService struct {
	repo    MerchantRepository
	bizRepo BusinessRepository
	logRepo LogRepository
}

func NewMerchantService(r MerchantRepository, b BusinessRepository, l LogRepository) MerchantUseCase {
	return &merchantService{
		repo:    r,
		bizRepo: b,
		logRepo: l,
	}
}

func (s *merchantService) RegisterMerchant(ctx context.Context, businessID uuid.UUID) (*entity.Merchant, error) {
	// 1. Validate Business existence
	_, err := s.bizRepo.GetBusinessByID(ctx, businessID)
	if err != nil {
		return nil, errors.New("cannot register merchant: business does not exist")
	}

	now := time.Now()
	m := &entity.Merchant{
		ID:         uuid.New(),
		BusinessID: businessID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateMerchant(ctx, m); err != nil {
		return nil, err
	}

	s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "MERCHANT_REGISTERED",
		Actor:      "admin",
		ResourceID: m.ID.String(),
		Timestamp:  now,
	})

	return m, nil
}

func (s *merchantService) GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	return s.repo.GetMerchantByID(ctx, id)
}

func (s *merchantService) GetBusinessMerchants(ctx context.Context, businessID uuid.UUID) ([]entity.Merchant, error) {
	return s.repo.GetMerchantByBusinessID(ctx, businessID)
}

func (s *merchantService) RemoveMerchant(ctx context.Context, id uuid.UUID) error {
	// Check if exists before deleting for better error handling
	_, err := s.repo.GetMerchantByID(ctx, id)
	if err != nil {
		return errors.New("merchant not found")
	}

	if err := s.repo.DeleteMerchant(ctx, id); err != nil {
		return err
	}

	// Audit Log
	s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "MERCHANT_DELETED",
		Actor:      "admin",
		ResourceID: id.String(),
		Timestamp:  time.Now(),
	})

	return nil
}
