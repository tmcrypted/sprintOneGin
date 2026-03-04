package postgres

import (
	"context"

	"sprin1/internal/model"
	"sprin1/internal/service"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PVZRepository struct {
	pool *pgxpool.Pool
}

func NewPVZRepository(pool *pgxpool.Pool) *PVZRepository {
	return &PVZRepository{pool: pool}
}

// GetAll реализует выборку ПВЗ с учётом фильтра и пагинации с помощью squirrel.
func (r *PVZRepository) GetAll(ctx context.Context, filter service.PVZFilter) ([]*model.PVZ, int64, error) {
	builder := sq.
		Select(
			"id",
			"owner_id",
			"status",
			"city",
			"address",
			"company_name",
			"description",
			"contact_phone",
			"contact_telegram",
			"lat",
			"lng",
			"moderated_at",
			"moderated_by",
			"created_at",
			"updated_at",
		).
		From("pvz").
		OrderBy("created_at DESC").
		Limit(uint64(filter.Limit)).
		Offset(uint64(filter.Offset)).
		PlaceholderFormat(sq.Dollar)

	if filter.Status != nil {
		builder = builder.Where(sq.Eq{"status": *filter.Status})
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

	var list []*model.PVZ
	for rows.Next() {
		var p model.PVZ
		if err := rows.Scan(
			&p.ID,
			&p.OwnerID,
			&p.Status,
			&p.City,
			&p.Address,
			&p.CompanyName,
			&p.Description,
			&p.ContactPhone,
			&p.ContactTelegram,
			&p.Lat,
			&p.Lng,
			&p.ModeratedAt,
			&p.ModeratedBy,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Отдельный запрос для общего количества записей без offset/limit.
	countBuilder := sq.
		Select("COUNT(*)").
		From("pvz").
		PlaceholderFormat(sq.Dollar)
	if filter.Status != nil {
		countBuilder = countBuilder.Where(sq.Eq{"status": *filter.Status})
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

func (r *PVZRepository) Create(ctx context.Context, pvz *model.PVZ) error {
	builder := sq.
		Insert("pvz").
		Columns(
			"owner_id",
			"status",
			"city",
			"address",
			"company_name",
			"description",
			"contact_phone",
			"contact_telegram",
			"lat",
			"lng",
		).
		Values(
			pvz.OwnerID,
			pvz.Status,
			pvz.City,
			pvz.Address,
			pvz.CompanyName,
			pvz.Description,
			pvz.ContactPhone,
			pvz.ContactTelegram,
			pvz.Lat,
			pvz.Lng,
		).
		Suffix("RETURNING id, owner_id, status, city, address, company_name, description, contact_phone, contact_telegram, lat, lng, moderated_at, moderated_by, created_at, updated_at").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	return r.pool.QueryRow(ctx, query, args...).Scan(
		&pvz.ID,
		&pvz.OwnerID,
		&pvz.Status,
		&pvz.City,
		&pvz.Address,
		&pvz.CompanyName,
		&pvz.Description,
		&pvz.ContactPhone,
		&pvz.ContactTelegram,
		&pvz.Lat,
		&pvz.Lng,
		&pvz.ModeratedAt,
		&pvz.ModeratedBy,
		&pvz.CreatedAt,
		&pvz.UpdatedAt,
	)
}

func (r *PVZRepository) GetByID(ctx context.Context, id int64) (*model.PVZ, error) {
	builder := sq.
		Select(
			"id",
			"owner_id",
			"status",
			"city",
			"address",
			"company_name",
			"description",
			"contact_phone",
			"contact_telegram",
			"lat",
			"lng",
			"moderated_at",
			"moderated_by",
			"created_at",
			"updated_at",
		).
		From("pvz").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var p model.PVZ
	if err := r.pool.QueryRow(ctx, query, args...).Scan(
		&p.ID,
		&p.OwnerID,
		&p.Status,
		&p.City,
		&p.Address,
		&p.CompanyName,
		&p.Description,
		&p.ContactPhone,
		&p.ContactTelegram,
		&p.Lat,
		&p.Lng,
		&p.ModeratedAt,
		&p.ModeratedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PVZRepository) Moderate(ctx context.Context, id int64, status model.PVZStatus, moderatorID int64) error {
	builder := sq.
		Update("pvz").
		Set("status", status).
		Set("moderated_by", moderatorID).
		Set("moderated_at", sq.Expr("NOW()")).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}