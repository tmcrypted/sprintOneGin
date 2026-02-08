package model

import "time"

type ListingStatus string

const (
	ListingStatusActive ListingStatus = "active"
	ListingStatusClosed ListingStatus = "closed"
)

type Listing struct {
	ID           int64         `json:"id" db:"id"`
	OwnerID      int64         `json:"owner_id" db:"owner_id"`
	PvzID        int64         `json:"pvz_id" db:"pvz_id"`
	CellsCount   int           `json:"cells_count" db:"cells_count"`
	PayPerShift  int           `json:"pay_per_shift" db:"pay_per_shift"`
	ShiftDate    time.Time     `json:"shift_date" db:"shift_date"`
	Status       ListingStatus `json:"status" db:"status"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}
