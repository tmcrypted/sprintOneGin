package dto

// CreateReviewRequest — тело запроса POST /reviews
type CreateReviewRequest struct {
	DealID       int64   `json:"deal_id" binding:"required"`
	PvzID        int64   `json:"pvz_id" binding:"required"`
	AuthorID     int64   `json:"author_id" binding:"required"`
	TargetUserID int64   `json:"target_user_id" binding:"required"`
	Rating       int     `json:"rating" binding:"required"`
	Body         *string `json:"body"`
}
