package handler

import (
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type SessionHandler struct {
	sessions *sqlite.SessionsRepo
	catches  *sqlite.CatchesRepo
}

func (h *SessionHandler) Register(r fiber.Router) {
	r.Get("/:id", h.Show)
}

func (h *SessionHandler) List(c *fiber.Ctx) error {
	sessions, err := h.sessions.ListWithCatchCount(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.SessionsList(sessions))
}

func (h *SessionHandler) Show(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	session, err := h.sessions.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	catches, err := h.catches.ListBySession(c.Context(), id)
	if err != nil {
		return err
	}
	locations, err := h.sessions.ListSessionLocations(c.Context(), id)
	if err != nil {
		return err
	}
	return render(c, pages.SessionDetail(session, catches, locations))
}
