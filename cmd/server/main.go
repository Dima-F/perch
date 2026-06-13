package main

import (
	"log"
	"os"

	"perch/internal/db"
	"perch/internal/handler"
	sqliterepo "perch/internal/repository/sqlite"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/mnt/c/Users/fiial/Dropbox/fishing/perch"
	}

	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer database.Close()

	sessions := sqliterepo.NewSessions(database)
	catches := sqliterepo.NewCatches(database)
	locations := sqliterepo.NewLocations(database)
	lures := sqliterepo.NewLures(database)
	lureModels := sqliterepo.NewLureModels(database)
	lureTypes := sqliterepo.NewLureTypes(database)
	fishingTypes := sqliterepo.NewFishingTypes(database)
	waterBodies := sqliterepo.NewWaterBodies(database)
	waterBodyTypes := sqliterepo.NewWaterBodyTypes(database)
	brands := sqliterepo.NewBrands(database)
	fishSpecies := sqliterepo.NewFishSpecies(database)

	h := handler.New(sessions, catches, locations, lures, lureModels, lureTypes, fishingTypes, waterBodies, waterBodyTypes, brands, fishSpecies)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).SendString(err.Error())
		},
	})

	app.Use(logger.New())
	app.Static("/static", "./static")

	app.Get("/", h.Sessions.List)
	h.Sessions.Register(app.Group("/sessions"))
	h.Catches.Register(app.Group("/catches"))
	h.Locations.Register(app.Group("/locations"))
	h.Lures.Register(app.Group("/lures"))
	h.LureModels.Register(app.Group("/lure-models"))
	h.WaterBodies.Register(app.Group("/water-bodies"))
	h.Brands.Register(app.Group("/brands"))
	h.FishSpecies.Register(app.Group("/fish-species"))

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3000"
	}
	log.Printf("listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
