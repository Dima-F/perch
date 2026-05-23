package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type WaterBodyTypesRepo struct{ db *sql.DB }

func NewWaterBodyTypes(db *sql.DB) *WaterBodyTypesRepo { return &WaterBodyTypesRepo{db: db} }

func (r *WaterBodyTypesRepo) List(ctx context.Context) ([]models.WaterBodyType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM water_body_type ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.WaterBodyType
	for rows.Next() {
		var t models.WaterBodyType
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
