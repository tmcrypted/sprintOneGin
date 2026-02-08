package model

import "time"

type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusAccepted  ApplicationStatus = "accepted"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
)

// Application — заявка на объявление. status=accepted => сделка (deal).
type Application struct {
	ID          int64              `json:"id" db:"id"`
	ListingID   int64              `json:"listing_id" db:"listing_id"`
	ApplicantID int64              `json:"applicant_id" db:"applicant_id"`
	Status      ApplicationStatus  `json:"status" db:"status"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}
