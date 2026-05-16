package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type LocationsRepo struct{ db *sql.DB }

func NewLocations(db *sql.DB) *LocationsRepo { return &LocationsRepo{db: db} }

func (r *LocationsRepo) List(ctx context.Context) ([]models.Location, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.name, l.region, l.notes, l.waterbody_id,
		       w.name, w.type
		FROM locations l
		JOIN water_body w ON w.id = l.waterbody_id
		ORDER BY l.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Location
	for rows.Next() {
		var l models.Location
		var wb models.WaterBody
		if err := rows.Scan(&l.ID, &l.Name, &l.Region, &l.Notes, &l.WaterBodyID, &wb.Name, &wb.Type); err != nil {
			return nil, err
		}
		wb.ID = l.WaterBodyID
		l.WaterBody = &wb
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *LocationsRepo) Get(ctx context.Context, id int) (*models.Location, error) {
	var l models.Location
	var wb models.WaterBody
	err := r.db.QueryRowContext(ctx, `
		SELECT l.id, l.name, l.region, l.notes, l.waterbody_id, w.name, w.type
		FROM locations l JOIN water_body w ON w.id = l.waterbody_id
		WHERE l.id = ?`, id).
		Scan(&l.ID, &l.Name, &l.Region, &l.Notes, &l.WaterBodyID, &wb.Name, &wb.Type)
	if err != nil {
		return nil, err
	}
	wb.ID = l.WaterBodyID
	l.WaterBody = &wb
	return &l, nil
}

func (r *LocationsRepo) Create(ctx context.Context, loc models.Location) (*models.Location, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO locations (name, region, notes, waterbody_id) VALUES (?, ?, ?, ?)`,
		loc.Name, loc.Region, loc.Notes, loc.WaterBodyID)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	loc.ID = int(id)
	return &loc, nil
}

func (r *LocationsRepo) Update(ctx context.Context, loc models.Location) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE locations SET name=?, region=?, notes=?, waterbody_id=? WHERE id=?`,
		loc.Name, loc.Region, loc.Notes, loc.WaterBodyID, loc.ID)
	return err
}

func (r *LocationsRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM locations WHERE id = ?`, id)
	return err
}
