package handler

import (
	"fmt"
	"time"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type SessionHandler struct {
	sessions *sqlite.SessionsRepo
	catches  *sqlite.CatchesRepo
}

func (h *SessionHandler) Register(r fiber.Router) {
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id/delete", h.Delete)
	r.Post("/:id", h.Update)
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

func (h *SessionHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/")
	}
	return render(c, pages.SessionFormDialog(models.FishingSession{}, ""))
}

func (h *SessionHandler) Create(c *fiber.Ctx) error {
	s, errMsg := sessionFromForm(c)
	if errMsg != "" {
		return render(c, pages.SessionFormDialog(s, errMsg))
	}
	created, err := h.sessions.Create(c.Context(), s)
	if err != nil {
		return render(c, pages.SessionFormDialog(s, err.Error()))
	}
	c.Set("HX-Retarget", "#sessions-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.SessionRowPartial(struct {
		models.FishingSession
		CatchCount int
	}{*created, 0}))
}

func (h *SessionHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	s, err := h.sessions.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	return render(c, pages.SessionFormDialog(*s, ""))
}

func (h *SessionHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	s, errMsg := sessionFromForm(c)
	s.ID = id
	if errMsg != "" {
		return render(c, pages.SessionFormDialog(s, errMsg))
	}
	if err := h.sessions.Update(c.Context(), s); err != nil {
		return render(c, pages.SessionFormDialog(s, err.Error()))
	}
	c.Set("HX-Retarget", fmt.Sprintf("#session-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	catchCount, _ := h.sessions.GetCatchCount(c.Context(), id)
	return render(c, pages.SessionRowPartial(struct {
		models.FishingSession
		CatchCount int
	}{s, catchCount}))
}

func (h *SessionHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.sessions.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

const datetimeLocal = "2006-01-02T15:04"

func sessionFromForm(c *fiber.Ctx) (models.FishingSession, string) {
	startStr := c.FormValue("start_time")
	if startStr == "" {
		return models.FishingSession{}, "Час початку — обов'язкове поле"
	}
	start, err := time.ParseInLocation(datetimeLocal, startStr, time.Local)
	if err != nil {
		return models.FishingSession{}, "Невірний формат часу початку"
	}
	s := models.FishingSession{StartTime: start}
	endStr := c.FormValue("end_time")
	if endStr == "" {
		return s, "Час кінця — обов'язкове поле"
	}
	end, err := time.ParseInLocation(datetimeLocal, endStr, time.Local)
	if err != nil {
		return s, "Невірний формат часу кінця"
	}
	if end.Before(start) {
		return s, "Час кінця не може бути раніше початку"
	}
	s.EndTime = end
	if notes := c.FormValue("notes"); notes != "" {
		s.Notes = &notes
	}
	return s, ""
}
