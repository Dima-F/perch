package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type LureTypesRepo struct{ db *sql.DB }

func NewLureTypes(db *sql.DB) *LureTypesRepo { return &LureTypesRepo{db: db} }

func (r *LureTypesRepo) List(ctx context.Context) ([]models.LureType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, notes FROM lure_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.LureType
	for rows.Next() {
		var t models.LureType
		if err := rows.Scan(&t.ID, &t.Name, &t.Notes); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
