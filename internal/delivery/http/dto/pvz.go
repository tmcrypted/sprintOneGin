package dto

import "sprin1/internal/model"

type GetAllPVZQuery struct {
	Page   int             `form:"page"`
	Limit  int             `form:"limit"`
	Status model.PVZStatus `form:"status"`
}

type CreatePVZRequest struct {
	City            string  `json:"city" binding:"required"`
	Address         string  `json:"address" binding:"required"`
	CompanyName     string  `json:"company_name" binding:"required"`
	Description     *string `json:"description,omitempty"`
	ContactPhone    string  `json:"contact_phone" binding:"required"`
	ContactTelegram *string `json:"contact_telegram,omitempty"`
}

type ModeratePVZRequest struct {
	ID     int64            `json:"id" binding:"required"`
	Status model.PVZStatus  `json:"status" binding:"required"`
}
