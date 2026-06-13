package sqlite

import (
	"context"
	"database/sql"
	"perch/internal/models"
)

type SessionsRepo struct{ db *sql.DB }

func NewSessions(db *sql.DB) *SessionsRepo { return &SessionsRepo{db: db} }

func (r *SessionsRepo) List(ctx context.Context) ([]models.FishingSession, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, start_time, end_time, notes FROM fishing_sessions ORDER BY start_time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.FishingSession
	for rows.Next() {
		var s models.FishingSession
		if err := rows.Scan(&s.ID, &s.StartTime, &s.EndTime, &s.Notes); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SessionsRepo) Get(ctx context.Context, id int) (*models.FishingSession, error) {
	var s models.FishingSession
	err := r.db.QueryRowContext(ctx,
		`SELECT id, start_time, end_time, notes FROM fishing_sessions WHERE id = ?`, id).
		Scan(&s.ID, &s.StartTime, &s.EndTime, &s.Notes)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionsRepo) Create(ctx context.Context, s models.FishingSession) (*models.FishingSession, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO fishing_sessions (start_time, end_time, notes) VALUES (?, ?, ?)`,
		s.StartTime, s.EndTime, s.Notes)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	s.ID = int(id)
	return &s, nil
}

func (r *SessionsRepo) Update(ctx context.Context, s models.FishingSession) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE fishing_sessions SET start_time=?, end_time=?, notes=? WHERE id=?`,
		s.StartTime, s.EndTime, s.Notes, s.ID)
	return err
}

func (r *SessionsRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fishing_sessions WHERE id = ?`, id)
	return err
}

// ListWithCatchCount повертає сесії з кількістю уловів
func (r *SessionsRepo) ListWithCatchCount(ctx context.Context) ([]struct {
	models.FishingSession
	CatchCount int
}, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.start_time, s.end_time, s.notes, COUNT(c.id) as catch_count
		FROM fishing_sessions s
		LEFT JOIN catches c ON c.session_id = s.id
		GROUP BY s.id
		ORDER BY s.start_time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct {
		models.FishingSession
		CatchCount int
	}
	for rows.Next() {
		var row struct {
			models.FishingSession
			CatchCount int
		}
		if err := rows.Scan(&row.ID, &row.StartTime, &row.EndTime, &row.Notes, &row.CatchCount); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// ListSessionLocations повертає локації для сесії
func (r *SessionsRepo) ListSessionLocations(ctx context.Context, sessionID int) ([]models.Location, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.name, l.region, l.notes, l.waterbody_id
		FROM session_locations sl
		JOIN locations l ON l.id = sl.location_id
		WHERE sl.session_id = ?`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.Name, &l.Region, &l.Notes, &l.WaterBodyID); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *SessionsRepo) AddLocation(ctx context.Context, sessionID, locationID int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO session_locations (session_id, location_id) VALUES (?, ?)`,
		sessionID, locationID)
	return err
}

func (r *SessionsRepo) RemoveLocation(ctx context.Context, sessionID, locationID int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM session_locations WHERE session_id=? AND location_id=?`,
		sessionID, locationID)
	return err
}

func (r *SessionsRepo) GetCatchCount(ctx context.Context, sessionID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM catches WHERE session_id = ?`, sessionID).Scan(&count)
	return count, err
}

// DB повертає sql.DB для зовнішнього використання
func (r *SessionsRepo) DB() *sql.DB { return r.db }
