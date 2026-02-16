package dto

import "sprin1/internal/model"

type CreateUserRequest struct {
	Email    string         `json:"email" binding:"required"`
	Password string         `json:"password" binding:"required"`
	FIO      string         `json:"fio" binding:"required"`
	Role     model.UserRole `json:"role"`
}
