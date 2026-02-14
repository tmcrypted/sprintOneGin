package service

import (
	"context"
	"errors"

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

func (s *userService) CreateUser(ctx context.Context, email, password, fio string, role model.UserRole) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if fio == "" {
		return nil, errors.New("fio is required")
	}
	if role == "" {
		role = model.RoleWorker
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		FIO:          fio,
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
