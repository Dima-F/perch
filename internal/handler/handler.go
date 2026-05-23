package handler

import (
	"perch/internal/repository/sqlite"
)

type Handlers struct {
	Sessions    *SessionHandler
	Catches     *CatchHandler
	Locations   *LocationHandler
	Lures       *LureHandler
	WaterBodies *WaterBodyHandler
	Brands      *BrandHandler
	FishSpecies *FishSpeciesHandler
}

func New(
	sessions *sqlite.SessionsRepo,
	catches *sqlite.CatchesRepo,
	locations *sqlite.LocationsRepo,
	lures *sqlite.LuresRepo,
	techniques *sqlite.TechniquesRepo,
	waterBodies *sqlite.WaterBodiesRepo,
	waterBodyTypes *sqlite.WaterBodyTypesRepo,
	brands *sqlite.BrandsRepo,
	fishSpecies *sqlite.FishSpeciesRepo,
) *Handlers {
	return &Handlers{
		Sessions:    &SessionHandler{sessions: sessions, catches: catches},
		Catches:     &CatchHandler{catches: catches, sessions: sessions, locations: locations, lures: lures},
		Locations:   &LocationHandler{locations: locations, waterBodies: waterBodies},
		Lures:       &LureHandler{lures: lures},
		WaterBodies: &WaterBodyHandler{waterBodies: waterBodies, waterBodyTypes: waterBodyTypes},
		Brands:      &BrandHandler{brands: brands},
		FishSpecies: &FishSpeciesHandler{fishSpecies: fishSpecies},
	}
}
