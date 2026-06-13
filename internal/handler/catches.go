package handler

import (
	"perch/internal/repository/sqlite"

	"github.com/gofiber/fiber/v2"
)

type CatchHandler struct {
	catches   *sqlite.CatchesRepo
	sessions  *sqlite.SessionsRepo
	locations *sqlite.LocationsRepo
	lures     *sqlite.LuresRepo
}

func (h *CatchHandler) Register(r fiber.Router) {
}
