package postgres

import (
	"context"
	"log"

	"sprin1/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshSessionRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshSessionRepository(pool *pgxpool.Pool) *RefreshSessionRepository {
	return &RefreshSessionRepository{pool: pool}
}

func (r *RefreshSessionRepository) Create(ctx context.Context, session *model.RefreshSession) error {
	q := `INSERT INTO refresh_sessions (user_id, token_hash, expires_at, created_at)
		  VALUES ($1, $2, $3, NOW())
		  RETURNING id, created_at`
	err := r.pool.QueryRow(ctx, q,
		session.UserID, session.TokenHash, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
	if err != nil {
		log.Printf("refresh_session_repository: Create failed user_id=%d: %v", session.UserID, err)
		return err
	}
	return nil
}

func (r *RefreshSessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshSession, error) {
	q := `SELECT id, user_id, token_hash, expires_at, created_at
		  FROM refresh_sessions WHERE token_hash = $1`
	var s model.RefreshSession
	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(
		&s.ID, &s.UserID, &s.TokenHash, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		log.Printf("refresh_session_repository: GetByTokenHash failed: %v", err)
		return nil, err
	}
	return &s, nil
}

func (r *RefreshSessionRepository) DeleteByID(ctx context.Context, id int64) error {
	q := `DELETE FROM refresh_sessions WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		log.Printf("refresh_session_repository: DeleteByID failed id=%d: %v", id, err)
		return err
	}
	return nil
}

