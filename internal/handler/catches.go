package handler

import (
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type CatchHandler struct {
	catches   *sqlite.CatchesRepo
	sessions  *sqlite.SessionsRepo
	locations *sqlite.LocationsRepo
	lures     *sqlite.LuresRepo
}

func (h *CatchHandler) List(c *fiber.Ctx) error {
	catches, err := h.catches.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.CatchesList(catches))
}
