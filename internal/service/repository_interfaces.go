package service

import (
	"context"

	"sprin1/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	UpdateRatingAvg(ctx context.Context, userID int64, avg float64) error
}

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error
	GetAvgRatingByTargetUser(ctx context.Context, targetUserID int64) (float64, error)
}
