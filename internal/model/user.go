package model

import "time"

type UserRole string

const (
	RoleOwner    UserRole = "owner"
	RoleWorker   UserRole = "worker"
	RoleModerator UserRole = "moderator"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	FIO          string    `json:"fio" db:"fio"`
	PhotoURL     *string   `json:"photo_url,omitempty" db:"photo_url"`
	Lat          *float64  `json:"lat,omitempty" db:"lat"`
	Lng          *float64  `json:"lng,omitempty" db:"lng"`
	RatingAvg    float64   `json:"rating_avg" db:"rating_avg"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
