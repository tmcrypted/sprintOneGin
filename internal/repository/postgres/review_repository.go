package postgres

import (
	"context"

	"sprin1/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{pool: pool}
}

func (r *ReviewRepository) Create(ctx context.Context, review *model.Review) error {
	q := `INSERT INTO reviews (deal_id, pvz_id, author_id, target_user_id, rating, body, created_at)
		  VALUES ($1, $2, $3, $4, $5, $6, NOW())
		  RETURNING id, created_at`
	return r.pool.QueryRow(ctx, q,
		review.DealID, review.PvzID, review.AuthorID, review.TargetUserID, review.Rating, review.Body,
	).Scan(&review.ID, &review.CreatedAt)
}

func (r *ReviewRepository) GetAvgRatingByTargetUser(ctx context.Context, targetUserID int64) (float64, error) {
	q := `SELECT COALESCE(AVG(rating)::numeric(3,2), 0) FROM reviews WHERE target_user_id = $1`
	var avg float64
	err := r.pool.QueryRow(ctx, q, targetUserID).Scan(&avg)
	return avg, err
}
