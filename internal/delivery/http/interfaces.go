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
	DeleteUser(ctx context.Context, id int64) error
}

type ReviewService interface {
	CreateReview(ctx context.Context, body dto.CreateReviewRequest) (*model.Review, error)
}

type AuthService interface {
	Login(ctx context.Context, body dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, body dto.CreateUserRequest) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
}

type PVZService interface {
	// GetAllPVZ возвращает страницу ПВЗ и общее количество записей.
	GetAllPVZ(ctx context.Context, q dto.GetAllPVZQuery) ([]*model.PVZ, int64, error)
	CreatePVZ(ctx context.Context, ownerID int64, body dto.CreatePVZRequest) (*model.PVZ, error)
	GetPVZ(ctx context.Context, id int64) (*model.PVZ, error)
	ModeratePVZ(ctx context.Context, moderatorID int64, body dto.ModeratePVZRequest) (*model.PVZ, error)
}
