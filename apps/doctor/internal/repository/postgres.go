package repository

import (
	"context"
	"doctor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, doctor *model.Doctor) error
	GetByID(ctx context.Context, id int64) (*model.Doctor, error)
	Update(ctx context.Context, doctor *model.Doctor) error
	Delete(ctx context.Context, id int64) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, d *model.Doctor) error {
	query := `INSERT INTO doctors (name, specialty, email) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(ctx, query, d.Name, d.Specialty, d.Email).Scan(&d.ID)
}

func (r *repo) GetByID(ctx context.Context, id int64) (*model.Doctor, error) {
	query := `SELECT id, name, specialty, email FROM doctors WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var d model.Doctor
	err := row.Scan(&d.ID, &d.Name, &d.Specialty, &d.Email)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) Update(ctx context.Context, d *model.Doctor) error {
	query := `UPDATE doctors SET name = $1, specialty = $2, email = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, d.Name, d.Specialty, d.Email, d.ID)
	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM doctors WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
