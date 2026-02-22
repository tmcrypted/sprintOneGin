package postgres

import (
	"context"

	"sprin1/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	q := `INSERT INTO users (email, password_hash, role, fio, photo_url, lat, lng, rating_avg, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		  RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		user.Email, user.PasswordHash, user.Role, user.FIO, user.PhotoURL, user.Lat, user.Lng, user.RatingAvg,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	q := `SELECT id, email, password_hash, role, fio, photo_url, lat, lng, rating_avg, created_at, updated_at
		  FROM users WHERE id = $1`
	var u model.User
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FIO, &u.PhotoURL, &u.Lat, &u.Lng, &u.RatingAvg,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	q := `SELECT id, email, password_hash, role, fio, photo_url, lat, lng, rating_avg, created_at, updated_at
		  FROM users WHERE email = $1`
	var u model.User
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FIO, &u.PhotoURL, &u.Lat, &u.Lng, &u.RatingAvg,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	q := `SELECT id, email, password_hash, role, fio, photo_url, lat, lng, rating_avg, created_at, updated_at
		  FROM users ORDER BY id`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FIO, &u.PhotoURL, &u.Lat, &u.Lng, &u.RatingAvg, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *UserRepository) UpdateRatingAvg(ctx context.Context, userID int64, avg float64) error {
	q := `UPDATE users SET rating_avg = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, q, avg, userID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}
