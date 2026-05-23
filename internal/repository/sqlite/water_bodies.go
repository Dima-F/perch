package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type WaterBodiesRepo struct{ db *sql.DB }

func NewWaterBodies(db *sql.DB) *WaterBodiesRepo { return &WaterBodiesRepo{db: db} }

func (r *WaterBodiesRepo) List(ctx context.Context) ([]models.WaterBody, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT wb.id, wb.name, wb.water_body_type_id, wbt.name
		FROM water_body wb
		JOIN water_body_type wbt ON wbt.id = wb.water_body_type_id
		ORDER BY wb.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.WaterBody
	for rows.Next() {
		var wb models.WaterBody
		if err := rows.Scan(&wb.ID, &wb.Name, &wb.WaterBodyTypeID, &wb.Type); err != nil {
			return nil, err
		}
		out = append(out, wb)
	}
	return out, rows.Err()
}

func (r *WaterBodiesRepo) Get(ctx context.Context, id int) (*models.WaterBody, error) {
	var wb models.WaterBody
	err := r.db.QueryRowContext(ctx, `
		SELECT wb.id, wb.name, wb.water_body_type_id, wbt.name
		FROM water_body wb
		JOIN water_body_type wbt ON wbt.id = wb.water_body_type_id
		WHERE wb.id = ?`, id).
		Scan(&wb.ID, &wb.Name, &wb.WaterBodyTypeID, &wb.Type)
	if err != nil {
		return nil, err
	}
	return &wb, nil
}

func (r *WaterBodiesRepo) Create(ctx context.Context, wb models.WaterBody) (*models.WaterBody, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO water_body (name, water_body_type_id) VALUES (?, ?)`,
		wb.Name, wb.WaterBodyTypeID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	wb.ID = int(id)
	return &wb, nil
}

func (r *WaterBodiesRepo) Update(ctx context.Context, wb models.WaterBody) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE water_body SET name = ?, water_body_type_id = ? WHERE id = ?`,
		wb.Name, wb.WaterBodyTypeID, wb.ID)
	return err
}

func (r *WaterBodiesRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM water_body WHERE id = ?`, id)
	return err
}
