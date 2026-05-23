package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type BrandsRepo struct{ db *sql.DB }

func NewBrands(db *sql.DB) *BrandsRepo { return &BrandsRepo{db: db} }

func (r *BrandsRepo) List(ctx context.Context) ([]models.Brand, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, notes FROM brands ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Brand
	for rows.Next() {
		var b models.Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Notes); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BrandsRepo) Get(ctx context.Context, id int) (*models.Brand, error) {
	var b models.Brand
	err := r.db.QueryRowContext(ctx, `SELECT id, name, notes FROM brands WHERE id = ?`, id).
		Scan(&b.ID, &b.Name, &b.Notes)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BrandsRepo) Create(ctx context.Context, b models.Brand) (*models.Brand, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO brands (name, notes) VALUES (?, ?)`, b.Name, b.Notes)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	b.ID = int(id)
	return &b, nil
}

func (r *BrandsRepo) Update(ctx context.Context, b models.Brand) error {
	_, err := r.db.ExecContext(ctx, `UPDATE brands SET name = ?, notes = ? WHERE id = ?`, b.Name, b.Notes, b.ID)
	return err
}

func (r *BrandsRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM brands WHERE id = ?`, id)
	return err
}
