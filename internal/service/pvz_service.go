package service

import (
	"context"
	"errors"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"
)

type pvzService struct {
	repo PVZRepository
}

func NewPVZService(repo PVZRepository) *pvzService {
	return &pvzService{repo: repo}
}

// GetAllPVZ возвращает страницу ПВЗ по фильтру (статус + пагинация)
// и общее количество записей без учёта offset/limit.
func (s *pvzService) GetAllPVZ(ctx context.Context, q dto.GetAllPVZQuery) ([]*model.PVZ, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	offset := (q.Page - 1) * q.Limit

	var status *model.PVZStatus
	if q.Status != "" {
		st := q.Status
		status = &st
	}

	filter := PVZFilter{
		Status: status,
		Offset: offset,
		Limit:  q.Limit,
	}

	return s.repo.GetAll(ctx, filter)
}

func (s *pvzService) CreatePVZ(ctx context.Context, ownerID int64, body dto.CreatePVZRequest) (*model.PVZ, error) {
	pvz := &model.PVZ{
		OwnerID:         ownerID,
		Status:          model.PVZStatusPending,
		City:            body.City,
		Address:         body.Address,
		CompanyName:     body.CompanyName,
		Description:     body.Description,
		ContactPhone:    body.ContactPhone,
		ContactTelegram: body.ContactTelegram,
	}

	if err := s.repo.Create(ctx, pvz); err != nil {
		return nil, err
	}

	return pvz, nil
}

func (s *pvzService) GetPVZ(ctx context.Context, id int64) (*model.PVZ, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *pvzService) ModeratePVZ(ctx context.Context, moderatorID int64, body dto.ModeratePVZRequest) (*model.PVZ, error) {
	if body.Status != model.PVZStatusApproved && body.Status != model.PVZStatusRejected {
		return nil, errors.New("invalid status for moderation")
	}

	if err := s.repo.Moderate(ctx, body.ID, body.Status, moderatorID); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, body.ID)
}
