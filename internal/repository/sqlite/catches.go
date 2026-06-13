package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type CatchesRepo struct{ db *sql.DB }

func NewCatches(db *sql.DB) *CatchesRepo { return &CatchesRepo{db: db} }

const catchSelect = `
	SELECT c.id, c.session_id, c.fish_id, c.lure_id, c.count,
	       c.avg_length_cm, c.max_length_cm, c.notes, c.weight_g,
	       c.jig_weight_g, c.jig_setup,
	       f.name,
	       m.name, b.name, l.color
	FROM catches c
	JOIN fish_species f ON f.id = c.fish_id
	LEFT JOIN lures l ON l.id = c.lure_id
	LEFT JOIN models m ON m.id = l.model_id
	LEFT JOIN brands b ON b.id = m.brand_id`

func scanCatch(row interface{ Scan(...any) error }) (*models.Catch, error) {
	var c models.Catch
	var fish models.FishSpecies
	var lureModel, lureBrand, lureColor sql.NullString
	err := row.Scan(
		&c.ID, &c.SessionID, &c.FishID, &c.LureID, &c.Count,
		&c.AvgLengthCm, &c.MaxLengthCm, &c.Notes, &c.WeightG,
		&c.JigWeightG, &c.JigSetup,
		&fish.Name,
		&lureModel, &lureBrand, &lureColor,
	)
	if err != nil {
		return nil, err
	}
	fish.ID = c.FishID
	c.Fish = &fish
	if lureModel.Valid && c.LureID != nil {
		lm := models.LureModel{Name: lureModel.String}
		if lureBrand.Valid {
			lm.Brand = &models.Brand{Name: lureBrand.String}
		}
		lure := models.Lure{ID: *c.LureID, LureModel: &lm}
		if lureColor.Valid {
			lure.Color = &lureColor.String
		}
		c.Lure = &lure
	}
	return &c, nil
}

func (r *CatchesRepo) List(ctx context.Context) ([]models.Catch, error) {
	rows, err := r.db.QueryContext(ctx, catchSelect+` ORDER BY c.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Catch
	for rows.Next() {
		c, err := scanCatch(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func (r *CatchesRepo) ListBySession(ctx context.Context, sessionID int) ([]models.Catch, error) {
	rows, err := r.db.QueryContext(ctx, catchSelect+` WHERE c.session_id = ? ORDER BY c.id`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Catch
	for rows.Next() {
		c, err := scanCatch(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func (r *CatchesRepo) Get(ctx context.Context, id int) (*models.Catch, error) {
	row := r.db.QueryRowContext(ctx, catchSelect+` WHERE c.id = ?`, id)
	return scanCatch(row)
}

func (r *CatchesRepo) Create(ctx context.Context, c models.Catch) (*models.Catch, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO catches (session_id, fish_id, lure_id, count, avg_length_cm, max_length_cm,
		                     notes, weight_g, jig_weight_g, jig_setup)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.SessionID, c.FishID, c.LureID, c.Count, c.AvgLengthCm, c.MaxLengthCm,
		c.Notes, c.WeightG, c.JigWeightG, c.JigSetup,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	c.ID = int(id)
	return &c, nil
}

func (r *CatchesRepo) Update(ctx context.Context, c models.Catch) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE catches SET fish_id=?, lure_id=?, count=?, avg_length_cm=?, max_length_cm=?,
		                   notes=?, weight_g=?, jig_weight_g=?, jig_setup=?
		WHERE id=?`,
		c.FishID, c.LureID, c.Count, c.AvgLengthCm, c.MaxLengthCm,
		c.Notes, c.WeightG, c.JigWeightG, c.JigSetup, c.ID,
	)
	return err
}

func (r *CatchesRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM catches WHERE id = ?`, id)
	return err
}
