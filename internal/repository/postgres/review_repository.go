package postgres

import (
	"context"

	"sprin1/internal/model"
	"sprin1/internal/service"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{pool: pool}
}

func (r *ReviewRepository) Create(ctx context.Context, review *model.Review) error {
	builder := sq.
		Insert("reviews").
		Columns(
			"deal_id",
			"pvz_id",
			"author_id",
			"target_user_id",
			"rating",
			"body",
		).
		Values(
			review.DealID,
			review.PvzID,
			review.AuthorID,
			review.TargetUserID,
			review.Rating,
			review.Body,
		).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	return r.pool.QueryRow(ctx, query, args...).Scan(&review.ID, &review.CreatedAt)
}

func (r *ReviewRepository) GetAvgRatingByTargetUser(ctx context.Context, targetUserID int64) (float64, error) {
	var avg float64
	builder := sq.
		Select("COALESCE(AVG(rating)::numeric(3,2), 0)").
		From("reviews").
		Where(sq.Eq{"target_user_id": targetUserID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&avg)
	return avg, err
}

func (r *ReviewRepository) Delete(ctx context.Context, id int64) error {
	builder := sq.
		Delete("reviews").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

// GetAll возвращает страницу отзывов и общее количество записей с учётом фильтра и пагинации.
func (r *ReviewRepository) GetAll(ctx context.Context, filter service.ReviewFilter) ([]*model.Review, int64, error) {
	builder := sq.
		Select(
			"id",
			"deal_id",
			"pvz_id",
			"author_id",
			"target_user_id",
			"rating",
			"body",
			"created_at",
		).
		From("reviews").
		OrderBy("created_at DESC").
		Limit(uint64(filter.Limit)).
		Offset(uint64(filter.Offset)).
		PlaceholderFormat(sq.Dollar)

	if filter.PvzID != nil {
		builder = builder.Where(sq.Eq{"pvz_id": *filter.PvzID})
	}
	if filter.TargetUserID != nil {
		builder = builder.Where(sq.Eq{"target_user_id": *filter.TargetUserID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.Review
	for rows.Next() {
		var rev model.Review
		if err := rows.Scan(
			&rev.ID,
			&rev.DealID,
			&rev.PvzID,
			&rev.AuthorID,
			&rev.TargetUserID,
			&rev.Rating,
			&rev.Body,
			&rev.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, &rev)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Отдельный запрос для общего количества записей без offset/limit.
	countBuilder := sq.
		Select("COUNT(*)").
		From("reviews").
		PlaceholderFormat(sq.Dollar)
	if filter.PvzID != nil {
		countBuilder = countBuilder.Where(sq.Eq{"pvz_id": *filter.PvzID})
	}
	if filter.TargetUserID != nil {
		countBuilder = countBuilder.Where(sq.Eq{"target_user_id": *filter.TargetUserID})
	}

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
