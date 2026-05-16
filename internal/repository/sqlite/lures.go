package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type LuresRepo struct{ db *sql.DB }

func NewLures(db *sql.DB) *LuresRepo { return &LuresRepo{db: db} }

func (r *LuresRepo) List(ctx context.Context) ([]models.Lure, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.color, l.size, l.notes, l.weight_g, l.model_id,
		       m.name, m.notes, m.luretype_id, m.brand_id,
		       lt.name, b.name
		FROM lures l
		JOIN models m ON m.id = l.model_id
		JOIN lure_types lt ON lt.id = m.luretype_id
		LEFT JOIN brands b ON b.id = m.brand_id
		ORDER BY m.name, l.color`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Lure
	for rows.Next() {
		var lu models.Lure
		var m models.LureModel
		var lt models.LureType
		var brand models.Brand
		var brandName sql.NullString
		if err := rows.Scan(
			&lu.ID, &lu.Color, &lu.Size, &lu.Notes, &lu.WeightG, &lu.ModelID,
			&m.Name, &m.Notes, &m.LureTypeID, &m.BrandID,
			&lt.Name, &brandName,
		); err != nil {
			return nil, err
		}
		m.ID = lu.ModelID
		m.LureType = &lt
		if brandName.Valid {
			brand.Name = brandName.String
			m.Brand = &brand
		}
		lu.LureModel = &m
		out = append(out, lu)
	}
	return out, rows.Err()
}

func (r *LuresRepo) Get(ctx context.Context, id int) (*models.Lure, error) {
	var lu models.Lure
	var m models.LureModel
	var lt models.LureType
	var brand models.Brand
	var brandName sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT l.id, l.color, l.size, l.notes, l.weight_g, l.model_id,
		       m.name, m.notes, m.luretype_id, m.brand_id, lt.name, b.name
		FROM lures l
		JOIN models m ON m.id = l.model_id
		JOIN lure_types lt ON lt.id = m.luretype_id
		LEFT JOIN brands b ON b.id = m.brand_id
		WHERE l.id = ?`, id).
		Scan(&lu.ID, &lu.Color, &lu.Size, &lu.Notes, &lu.WeightG, &lu.ModelID,
			&m.Name, &m.Notes, &m.LureTypeID, &m.BrandID, &lt.Name, &brandName)
	if err != nil {
		return nil, err
	}
	m.ID = lu.ModelID
	m.LureType = &lt
	if brandName.Valid {
		brand.Name = brandName.String
		m.Brand = &brand
	}
	lu.LureModel = &m
	return &lu, nil
}

func (r *LuresRepo) Create(ctx context.Context, lure models.Lure) (*models.Lure, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO lures (color, size, notes, weight_g, model_id) VALUES (?, ?, ?, ?, ?)`,
		lure.Color, lure.Size, lure.Notes, lure.WeightG, lure.ModelID)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	lure.ID = int(id)
	return &lure, nil
}

func (r *LuresRepo) Update(ctx context.Context, lure models.Lure) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE lures SET color=?, size=?, notes=?, weight_g=?, model_id=? WHERE id=?`,
		lure.Color, lure.Size, lure.Notes, lure.WeightG, lure.ModelID, lure.ID)
	return err
}

func (r *LuresRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lures WHERE id = ?`, id)
	return err
}
