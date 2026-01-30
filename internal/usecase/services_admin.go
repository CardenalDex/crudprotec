package usecase

import (
	"context"
	"fmt"
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
	s.logRepo.CreateLog(ctx, &entity.Log{
		ID:             uuid.New(),
		Action:         "CREATE_BUSINESS",
		Actor:          "admin",
		ResourceID:     biz.ID.String(),
		PrevResourceID: "",
		Timestamp:      time.Now(),
	})

	return biz, nil
}

func (s *adminService) GetAuditTrail(ctx context.Context, resID string) ([]entity.Log, error) {
	return s.logRepo.GetLogByResource(ctx, resID)
}

func (s *adminService) GetBusiness(ctx context.Context, id uuid.UUID) (*entity.Business, error) {
	return s.bizRepo.GetBusinessByID(ctx, id)
}

func (s *adminService) UpdateBusinessCommission(ctx context.Context, id uuid.UUID, newCommission int64) (*entity.Business, error) {

	biz, err := s.bizRepo.GetBusinessByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldCommission := biz.Commission
	biz.Commission = newCommission
	biz.UpdatedAt = time.Now()

	if err := s.bizRepo.UpdateBusiness(ctx, biz); err != nil {
		return nil, err
	}

	s.logRepo.CreateLog(ctx, &entity.Log{
		ID:             uuid.New(),
		Action:         "UPDATE_BUSINESS_COMMISSION",
		Actor:          "admin",
		ResourceID:     id.String(),
		PrevResourceID: fmt.Sprintf("old_comm:%d", oldCommission),
		Timestamp:      time.Now(),
	})

	return biz, nil
}

func (s *adminService) RemoveBusiness(ctx context.Context, id uuid.UUID) error {

	if err := s.bizRepo.DeleteBusiness(ctx, id); err != nil {
		return err
	}

	_ = s.logRepo.CreateLog(ctx, &entity.Log{
		ID:         uuid.New(),
		Action:     "DELETE_BUSINESS",
		Actor:      "admin",
		ResourceID: id.String(),
		Timestamp:  time.Now(),
	})

	return nil
}

func (s *adminService) GetLogDetails(ctx context.Context, logID string) (entity.Log, error) {
	return s.logRepo.GetLogByID(ctx, logID)
}

func (s *adminService) GetAllLogs(ctx context.Context) ([]entity.Log, error) {
	return s.logRepo.GetAll(ctx)
}
