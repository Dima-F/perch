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
	techniques := sqliterepo.NewTechniques(database)

	h := handler.New(sessions, catches, locations, lures, techniques)

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
	app.Get("/sessions/:id", h.Sessions.Show)
	app.Get("/catches", h.Catches.List)
	app.Get("/locations", h.Locations.List)
	app.Get("/lures", h.Lures.List)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3000"
	}
	log.Printf("listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
