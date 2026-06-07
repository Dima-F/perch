package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type LureModelsRepo struct{ db *sql.DB }

func NewLureModels(db *sql.DB) *LureModelsRepo { return &LureModelsRepo{db: db} }

const lureModelSelect = `
	SELECT m.id, m.name, m.brand_id, m.luretype_id, b.name, lt.name
	FROM models m
	JOIN lure_types lt ON lt.id = m.luretype_id
	LEFT JOIN brands b ON b.id = m.brand_id`

func scanLureModel(row interface{ Scan(...any) error }) (*models.LureModel, error) {
	var m models.LureModel
	var brand models.Brand
	var lt models.LureType
	var brandName sql.NullString
	if err := row.Scan(&m.ID, &m.Name, &m.BrandID, &m.LureTypeID, &brandName, &lt.Name); err != nil {
		return nil, err
	}
	m.LureType = &lt
	if brandName.Valid {
		brand.Name = brandName.String
		m.Brand = &brand
	}
	return &m, nil
}

func (r *LureModelsRepo) List(ctx context.Context) ([]models.LureModel, error) {
	rows, err := r.db.QueryContext(ctx, lureModelSelect+` ORDER BY b.name, m.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.LureModel
	for rows.Next() {
		m, err := scanLureModel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *m)
	}
	return out, rows.Err()
}

func (r *LureModelsRepo) Get(ctx context.Context, id int) (*models.LureModel, error) {
	row := r.db.QueryRowContext(ctx, lureModelSelect+` WHERE m.id = ?`, id)
	return scanLureModel(row)
}

func (r *LureModelsRepo) Create(ctx context.Context, m models.LureModel) (*models.LureModel, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO models (name, brand_id, luretype_id, notes) VALUES (?, ?, ?, ?)`,
		m.Name, m.BrandID, m.LureTypeID, m.Notes)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	m.ID = int(id)
	return &m, nil
}

func (r *LureModelsRepo) Update(ctx context.Context, m models.LureModel) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE models SET name=?, brand_id=?, luretype_id=?, notes=? WHERE id=?`,
		m.Name, m.BrandID, m.LureTypeID, m.Notes, m.ID)
	return err
}

func (r *LureModelsRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM models WHERE id = ?`, id)
	return err
}
