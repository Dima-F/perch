package handler

import (
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	locations *sqlite.LocationsRepo
}

func (h *LocationHandler) List(c *fiber.Ctx) error {
	locs, err := h.locations.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LocationsList(locs))
}
