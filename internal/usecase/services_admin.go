package usecase

import (
	"context"
	"time"

	entity "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/google/uuid"
)

type adminService struct {
	bizRepo BusinessRepository
	logRepo LogRepository
}

func NewAdminService(br BusinessRepository, lr LogRepository) AdminUseCase {
	return &adminService{br, lr}
}

func (s *adminService) RegisterBusiness(ctx context.Context, commission int64) (*entity.Business, error) {
	biz := &entity.Business{
		ID:         uuid.New(),
		Commission: commission,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.bizRepo.CreateBusiness(ctx, biz); err != nil {
		return nil, err
	}

	return biz, nil
}

func (s *adminService) GetAuditTrail(ctx context.Context, resID string) ([]entity.Log, error) {
	return s.logRepo.GetLogByResource(ctx, resID)
}
