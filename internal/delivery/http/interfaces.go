package http

import (
	"context"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, body dto.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
}

type ReviewService interface {
	CreateReview(ctx context.Context, body dto.CreateReviewRequest) (*model.Review, error)
}
