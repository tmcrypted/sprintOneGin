package service

import (
	"context"
	"errors"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo UserRepository
}

// NewUserService возвращает сервис пользователей, работающий с БД через репозиторий.
func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, body dto.CreateUserRequest) (*model.User, error) {
	if body.Email == "" {
		return nil, errors.New("email is required")
	}
	if body.Password == "" {
		return nil, errors.New("password is required")
	}
	if body.FIO == "" {
		return nil, errors.New("fio is required")
	}
	if body.Role == "" {
		body.Role = model.RoleWorker
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        body.Email,
		PasswordHash: string(hash),
		Role:         body.Role,
		FIO:          body.FIO,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAll(ctx)
}
