package repository

import (
	"context"
	"perch/internal/models"
)

type FishingTypes interface {
	List(ctx context.Context) ([]models.FishingType, error)
	Get(ctx context.Context, id int) (*models.FishingType, error)
	Create(ctx context.Context, name string) (*models.FishingType, error)
	Update(ctx context.Context, id int, name string) error
	Delete(ctx context.Context, id int) error
}

type FishSpecies interface {
	List(ctx context.Context) ([]models.FishSpecies, error)
	Get(ctx context.Context, id int) (*models.FishSpecies, error)
	Create(ctx context.Context, name string) (*models.FishSpecies, error)
	Update(ctx context.Context, id int, name string) error
	Delete(ctx context.Context, id int) error
}

type WaterBodies interface {
	List(ctx context.Context) ([]models.WaterBody, error)
	Get(ctx context.Context, id int) (*models.WaterBody, error)
	Create(ctx context.Context, name, wbType string) (*models.WaterBody, error)
	Update(ctx context.Context, wb models.WaterBody) error
	Delete(ctx context.Context, id int) error
}

type Locations interface {
	List(ctx context.Context) ([]models.Location, error)
	Get(ctx context.Context, id int) (*models.Location, error)
	Create(ctx context.Context, loc models.Location) (*models.Location, error)
	Update(ctx context.Context, loc models.Location) error
	Delete(ctx context.Context, id int) error
}

type Brands interface {
	List(ctx context.Context) ([]models.Brand, error)
	Get(ctx context.Context, id int) (*models.Brand, error)
	Create(ctx context.Context, name string, notes *string) (*models.Brand, error)
	Update(ctx context.Context, b models.Brand) error
	Delete(ctx context.Context, id int) error
}

type Lures interface {
	List(ctx context.Context) ([]models.Lure, error)
	Get(ctx context.Context, id int) (*models.Lure, error)
	Create(ctx context.Context, lure models.Lure) (*models.Lure, error)
	Update(ctx context.Context, lure models.Lure) error
	Delete(ctx context.Context, id int) error
}

type LureModels interface {
	List(ctx context.Context) ([]models.LureModel, error)
	Get(ctx context.Context, id int) (*models.LureModel, error)
	Create(ctx context.Context, m models.LureModel) (*models.LureModel, error)
	Update(ctx context.Context, m models.LureModel) error
	Delete(ctx context.Context, id int) error
}

type Blanks interface {
	List(ctx context.Context) ([]models.Blank, error)
	Get(ctx context.Context, id int) (*models.Blank, error)
	Create(ctx context.Context, b models.Blank) (*models.Blank, error)
	Update(ctx context.Context, b models.Blank) error
	Delete(ctx context.Context, id int) error
}

type Reels interface {
	List(ctx context.Context) ([]models.Reel, error)
	Get(ctx context.Context, id int) (*models.Reel, error)
	Create(ctx context.Context, r models.Reel) (*models.Reel, error)
	Update(ctx context.Context, r models.Reel) error
	Delete(ctx context.Context, id int) error
}

type FishingSessions interface {
	List(ctx context.Context) ([]models.FishingSession, error)
	Get(ctx context.Context, id int) (*models.FishingSession, error)
	Create(ctx context.Context, s models.FishingSession) (*models.FishingSession, error)
	Update(ctx context.Context, s models.FishingSession) error
	Delete(ctx context.Context, id int) error
}

type Catches interface {
	List(ctx context.Context) ([]models.Catch, error)
	ListBySession(ctx context.Context, sessionID int) ([]models.Catch, error)
	Get(ctx context.Context, id int) (*models.Catch, error)
	Create(ctx context.Context, c models.Catch) (*models.Catch, error)
	Update(ctx context.Context, c models.Catch) error
	Delete(ctx context.Context, id int) error
}
