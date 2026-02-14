package http

import (
	"context"

	"sprin1/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, fio string, role model.UserRole) (*model.User, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
}

type ReviewService interface {
	CreateReview(ctx context.Context, dealID, pvzID, authorID, targetUserID int64, rating int, body *string) (*model.Review, error)
}
