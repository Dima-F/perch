package handler

import (
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type LureHandler struct {
	lures *sqlite.LuresRepo
}

func (h *LureHandler) List(c *fiber.Ctx) error {
	lures, err := h.lures.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LuresList(lures))
}
