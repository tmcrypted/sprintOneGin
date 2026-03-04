package dto

// CreateReviewRequest — тело запроса POST /reviews
type CreateReviewRequest struct {
	DealID       int64   `json:"deal_id" binding:"required"`
	PvzID        int64   `json:"pvz_id" binding:"required"`
	// AuthorID заполняется на бэке из контекста пользователя.
	AuthorID     int64   `json:"author_id"`
	TargetUserID int64   `json:"target_user_id" binding:"required"`
	Rating       int     `json:"rating" binding:"required"`
	Body         *string `json:"body"`
}

// GetReviewsQuery — параметры запроса GET /reviews с пагинацией и фильтрами.
// Можно фильтровать по ПВЗ и/или по пользователю, которому оставлен отзыв.
type GetReviewsQuery struct {
	Page         int   `form:"page"`
	Limit        int   `form:"limit"`
	PvzID        int64 `form:"pvz_id"`
	TargetUserID int64 `form:"target_user_id"`
}
