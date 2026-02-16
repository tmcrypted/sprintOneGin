package model

import "time"

type Review struct {
	ID           int64     `json:"id" db:"id"`
	DealID       int64     `json:"deal_id" db:"deal_id"`
	PvzID        int64     `json:"pvz_id" db:"pvz_id"` // для выборки отзывов на странице ПВЗ
	AuthorID     int64     `json:"author_id" db:"author_id"`
	TargetUserID int64     `json:"target_user_id" db:"target_user_id"`
	Rating       int       `json:"rating" db:"rating"` // 1-5
	Body         *string   `json:"body,omitempty" db:"body"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
