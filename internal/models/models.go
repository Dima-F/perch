package models

import "time"

type FishingType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type FishSpecies struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type WaterBodyType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type WaterBody struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	WaterBodyTypeID int    `db:"water_body_type_id"`
	Type            string // populated from JOIN with water_body_type
}

type Location struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Region      string     `db:"region"`
	Notes       *string    `db:"notes"`
	WaterBodyID int        `db:"waterbody_id"`
	WaterBody   *WaterBody `db:"-"`
}

type Brand struct {
	ID    int     `db:"id"`
	Name  string  `db:"name"`
	Notes *string `db:"notes"`
}

type LureType struct {
	ID    int     `db:"id"`
	Name  string  `db:"name"`
	Notes *string `db:"notes"`
}

type LureModel struct {
	ID         int       `db:"id"`
	BrandID    *int      `db:"brand_id"`
	Brand      *Brand    `db:"-"`
	Name       string    `db:"name"`
	Notes      *string   `db:"notes"`
	LureTypeID int       `db:"luretype_id"`
	LureType   *LureType `db:"-"`
}

type Lure struct {
	ID        int        `db:"id"`
	Color     *string    `db:"color"`
	Size      *string    `db:"size"`
	Notes     *string    `db:"notes"`
	WeightG   *float64   `db:"weight_g"`
	ModelID   int        `db:"model_id"`
	LureModel *LureModel `db:"-"`
}

type BlankType struct {
	ID    int     `db:"id"`
	Name  string  `db:"name"`
	Notes *string `db:"notes"`
}

type Blank struct {
	ID          int        `db:"id"`
	BrandID     int        `db:"brand_id"`
	Brand       *Brand     `db:"-"`
	Name        string     `db:"name"`
	Casting     string     `db:"casting"`
	Length      string     `db:"length"`
	Line        *string    `db:"line"`
	BlankTypeID int        `db:"blank_type_id"`
	BlankType   *BlankType `db:"-"`
	Notes       *string    `db:"notes"`
	BoughtAt    *string    `db:"bought_at"`
}

type Reel struct {
	ID           int     `db:"id"`
	BrandID      *int    `db:"brand_id"`
	Brand        *Brand  `db:"-"`
	Name         string  `db:"name"`
	Notes        *string `db:"notes"`
	ReelSize     *int    `db:"reel_size"`
	BearingCount *string `db:"bearing_count"`
	GearRate     *string `db:"gear_rate"`
	BoughtAt     *string `db:"bought_at"`
}

type BraidedLine struct {
	ID        int     `db:"id"`
	BrandID   int     `db:"brand_id"`
	Brand     *Brand  `db:"-"`
	Name      string  `db:"name"`
	Notes     *string `db:"notes"`
	LineWidth string  `db:"line_width"`
	MaxLoad   *float64 `db:"max_load"`
	Color     *string `db:"color"`
	Length    *int    `db:"length"`
}

type Spool struct {
	ID          int     `db:"id"`
	ReelID      int     `db:"reel_id"`
	Reel        *Reel   `db:"-"`
	Notes       *string `db:"notes"`
	SpoolNumber int     `db:"spool_number"`
	Size        int     `db:"size"`
}

type FishingSession struct {
	ID        int        `db:"id"`
	StartTime time.Time  `db:"start_time"`
	EndTime   time.Time  `db:"end_time"`
	Notes     *string    `db:"notes"`
}

var JigSetups = []string{
	"шарнір", "джиг-ріг", "джиг-головка", "відвідний", "токіо-ріг", "мормишка",
}

type Catch struct {
	ID          int      `db:"id"`
	SessionID   int      `db:"session_id"`
	FishID      int      `db:"fish_id"`
	Fish        *FishSpecies `db:"-"`
	LureID      *int     `db:"lure_id"`
	Lure        *Lure    `db:"-"`
	Count       int      `db:"count"`
	AvgLengthCm *float64 `db:"avg_length_cm"`
	MaxLengthCm *float64 `db:"max_length_cm"`
	Notes       *string  `db:"notes"`
	WeightG     *int     `db:"weight_g"`
	JigWeightG  *float64 `db:"jig_weight_g"`
	JigSetup    *string  `db:"jig_setup"`
}
