package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type FishSpeciesRepo struct{ db *sql.DB }

func NewFishSpecies(db *sql.DB) *FishSpeciesRepo { return &FishSpeciesRepo{db: db} }

func (r *FishSpeciesRepo) List(ctx context.Context) ([]models.FishSpecies, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM fish_species ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.FishSpecies
	for rows.Next() {
		var f models.FishSpecies
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *FishSpeciesRepo) Get(ctx context.Context, id int) (*models.FishSpecies, error) {
	var f models.FishSpecies
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM fish_species WHERE id = ?`, id).
		Scan(&f.ID, &f.Name)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FishSpeciesRepo) Create(ctx context.Context, f models.FishSpecies) (*models.FishSpecies, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO fish_species (name) VALUES (?)`, f.Name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	f.ID = int(id)
	return &f, nil
}

func (r *FishSpeciesRepo) Update(ctx context.Context, f models.FishSpecies) error {
	_, err := r.db.ExecContext(ctx, `UPDATE fish_species SET name = ? WHERE id = ?`, f.Name, f.ID)
	return err
}

func (r *FishSpeciesRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fish_species WHERE id = ?`, id)
	return err
}
