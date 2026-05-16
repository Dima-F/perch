package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type TechniquesRepo struct{ db *sql.DB }

func NewTechniques(db *sql.DB) *TechniquesRepo { return &TechniquesRepo{db: db} }

func (r *TechniquesRepo) List(ctx context.Context) ([]models.Technique, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM techniques ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Technique
	for rows.Next() {
		var t models.Technique
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TechniquesRepo) Get(ctx context.Context, id int) (*models.Technique, error) {
	var t models.Technique
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM techniques WHERE id = ?`, id).
		Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TechniquesRepo) Create(ctx context.Context, name string) (*models.Technique, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO techniques (name) VALUES (?)`, name)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.Technique{ID: int(id), Name: name}, nil
}

func (r *TechniquesRepo) Update(ctx context.Context, id int, name string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE techniques SET name = ? WHERE id = ?`, name, id)
	return err
}

func (r *TechniquesRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM techniques WHERE id = ?`, id)
	return err
}
