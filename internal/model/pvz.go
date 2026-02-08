package model

import "time"

type PVZStatus string

const (
	PVZStatusPending   PVZStatus = "pending"
	PVZStatusApproved  PVZStatus = "approved"
	PVZStatusRejected  PVZStatus = "rejected"
)

type PVZ struct {
	ID             int64      `json:"id" db:"id"`
	OwnerID        int64      `json:"owner_id" db:"owner_id"`
	Status         PVZStatus  `json:"status" db:"status"`
	City           string     `json:"city" db:"city"`
	Address        string     `json:"address" db:"address"`
	CompanyName    string     `json:"company_name" db:"company_name"`
	Description    *string    `json:"description,omitempty" db:"description"`
	ContactPhone   string     `json:"contact_phone" db:"contact_phone"`
	ContactTelegram *string   `json:"contact_telegram,omitempty" db:"contact_telegram"`
	Lat            *float64   `json:"lat,omitempty" db:"lat"`
	Lng            *float64   `json:"lng,omitempty" db:"lng"`
	ModeratedAt    *time.Time `json:"moderated_at,omitempty" db:"moderated_at"`
	ModeratedBy    *int64     `json:"moderated_by,omitempty" db:"moderated_by"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
