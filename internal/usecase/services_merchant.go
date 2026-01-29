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

// NewMerchantService initializes the merchant use case implementation
func NewMerchantService(r MerchantRepository, b BusinessRepository, l LogRepository) MerchantUseCase {
	return &merchantService{
		repo:    r,
		bizRepo: b,
		logRepo: l,
	}
}

// RegisterMerchant creates a new merchant linked to an existing business
func (s *merchantService) RegisterMerchant(ctx context.Context, businessID uuid.UUID) (*entity.Merchant, error) {
	// 1. Validate Business existence
	_, err := s.bizRepo.GetBusinessByID(ctx, businessID)
	if err != nil {
		return nil, errors.New("cannot register merchant: business does not exist")
	}

	// 2. Create Entity
	now := time.Now()
	m := &entity.Merchant{
		ID:         uuid.New(),
		BusinessID: businessID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 3. Persist
	if err := s.repo.CreateMerchant(ctx, m); err != nil {
		return nil, err
	}

	// 4. Audit Log
	_ = s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "MERCHANT_REGISTERED",
		Actor:      "admin",
		ResourceID: m.ID.String(),
		Timestamp:  now,
	})

	return m, nil
}

// GetMerchant retrieves a single merchant by ID
func (s *merchantService) GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	return s.repo.GetMerchantByID(ctx, id)
}

// GetBusinessMerchants lists all merchants belonging to a specific business
func (s *merchantService) GetBusinessMerchants(ctx context.Context, businessID uuid.UUID) ([]entity.Merchant, error) {
	return s.repo.GetMerchantByBusinessID(ctx, businessID)
}

// RemoveMerchant performs a logic delete on the merchant
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
	_ = s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "MERCHANT_DELETED",
		Actor:      "admin",
		ResourceID: id.String(),
		Timestamp:  time.Now(),
	})

	return nil
}
