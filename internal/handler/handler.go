package handler

import (
	"perch/internal/repository/sqlite"
)

type Handlers struct {
	Sessions  *SessionHandler
	Catches   *CatchHandler
	Locations *LocationHandler
	Lures     *LureHandler
}

func New(
	sessions *sqlite.SessionsRepo,
	catches *sqlite.CatchesRepo,
	locations *sqlite.LocationsRepo,
	lures *sqlite.LuresRepo,
	techniques *sqlite.TechniquesRepo,
) *Handlers {
	return &Handlers{
		Sessions:  &SessionHandler{sessions: sessions, catches: catches},
		Catches:   &CatchHandler{catches: catches, sessions: sessions, locations: locations, lures: lures},
		Locations: &LocationHandler{locations: locations},
		Lures:     &LureHandler{lures: lures},
	}
}
