package model

import "time"

// Message — сообщение в рамках сделки (deal = application со status=accepted).
type Message struct {
	ID        int64     `json:"id" db:"id"`
	DealID    int64     `json:"deal_id" db:"deal_id"`
	SenderID  int64     `json:"sender_id" db:"sender_id"`
	Body      string    `json:"body" db:"body"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
