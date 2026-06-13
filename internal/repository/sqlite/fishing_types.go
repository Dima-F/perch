package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type FishingTypesRepo struct{ db *sql.DB }

func NewFishingTypes(db *sql.DB) *FishingTypesRepo { return &FishingTypesRepo{db: db} }

func (r *FishingTypesRepo) List(ctx context.Context) ([]models.FishingType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM fishing_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.FishingType
	for rows.Next() {
		var t models.FishingType
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *FishingTypesRepo) Get(ctx context.Context, id int) (*models.FishingType, error) {
	var t models.FishingType
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM fishing_types WHERE id = ?`, id).
		Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *FishingTypesRepo) Create(ctx context.Context, name string) (*models.FishingType, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO fishing_types (name) VALUES (?)`, name)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.FishingType{ID: int(id), Name: name}, nil
}

func (r *FishingTypesRepo) Update(ctx context.Context, id int, name string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE fishing_types SET name = ? WHERE id = ?`, name, id)
	return err
}

func (r *FishingTypesRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fishing_types WHERE id = ?`, id)
	return err
}
