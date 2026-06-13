package handler

import (
	"fmt"
	"strconv"
	"time"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type SessionHandler struct {
	sessions     *sqlite.SessionsRepo
	catches      *sqlite.CatchesRepo
	fishingTypes *sqlite.FishingTypesRepo
	locations    *sqlite.LocationsRepo
}

func (h *SessionHandler) Register(r fiber.Router) {
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/locations/picker", h.LocationPicker)
	r.Post("/:id/locations/:loc_id/remove", h.RemoveLocation)
	r.Post("/:id/locations", h.AddLocation)
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
	ftMap, err := h.sessions.ListAllSessionFishingTypes(c.Context())
	if err != nil {
		return err
	}
	for i := range sessions {
		sessions[i].FishingTypes = ftMap[sessions[i].ID]
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
	session.FishingTypes, _ = h.sessions.ListSessionFishingTypes(c.Context(), id)
	return render(c, pages.SessionDetail(session, catches, locations))
}

func (h *SessionHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/")
	}
	allTypes, err := h.fishingTypes.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.SessionFormDialog(models.FishingSession{}, allTypes, ""))
}

func (h *SessionHandler) Create(c *fiber.Ctx) error {
	s, errMsg := sessionFromForm(c)
	allTypes, _ := h.fishingTypes.List(c.Context())
	if errMsg != "" {
		return render(c, pages.SessionFormDialog(s, allTypes, errMsg))
	}
	created, err := h.sessions.Create(c.Context(), s)
	if err != nil {
		return render(c, pages.SessionFormDialog(s, allTypes, err.Error()))
	}
	typeIDs := fishingTypeIDsFromForm(c)
	if err := h.sessions.SetSessionFishingTypes(c.Context(), created.ID, typeIDs); err != nil {
		return err
	}
	created.FishingTypes, _ = h.sessions.ListSessionFishingTypes(c.Context(), created.ID)
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
	s.FishingTypes, _ = h.sessions.ListSessionFishingTypes(c.Context(), id)
	allTypes, err := h.fishingTypes.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.SessionFormDialog(*s, allTypes, ""))
}

func (h *SessionHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	s, errMsg := sessionFromForm(c)
	s.ID = id
	allTypes, _ := h.fishingTypes.List(c.Context())
	if errMsg != "" {
		s.FishingTypes, _ = h.sessions.ListSessionFishingTypes(c.Context(), id)
		return render(c, pages.SessionFormDialog(s, allTypes, errMsg))
	}
	if err := h.sessions.Update(c.Context(), s); err != nil {
		s.FishingTypes, _ = h.sessions.ListSessionFishingTypes(c.Context(), id)
		return render(c, pages.SessionFormDialog(s, allTypes, err.Error()))
	}
	typeIDs := fishingTypeIDsFromForm(c)
	if err := h.sessions.SetSessionFishingTypes(c.Context(), id, typeIDs); err != nil {
		return err
	}
	s.FishingTypes, _ = h.sessions.ListSessionFishingTypes(c.Context(), id)
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


func (h *SessionHandler) LocationPicker(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	all, err := h.locations.List(c.Context())
	if err != nil {
		return err
	}
	linked, err := h.sessions.ListSessionLocations(c.Context(), id)
	if err != nil {
		return err
	}
	linkedSet := make(map[int]bool, len(linked))
	for _, l := range linked {
		linkedSet[l.ID] = true
	}
	var available []models.Location
	for _, l := range all {
		if !linkedSet[l.ID] {
			available = append(available, l)
		}
	}
	return render(c, pages.LocationPickerDialog(id, available))
}

func (h *SessionHandler) AddLocation(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	locID, err := strconv.Atoi(c.FormValue("location_id"))
	if err != nil || locID == 0 {
		return fiber.ErrBadRequest
	}
	if err := h.sessions.AddLocation(c.Context(), id, locID); err != nil {
		return err
	}
	loc, err := h.locations.Get(c.Context(), locID)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", "#locations-tbody")
	c.Set("HX-Reswap", "beforeend")
	c.Set("HX-Trigger", "closeLocationDialog")
	return render(c, pages.SessionLocationRowPartial(id, *loc))
}

func (h *SessionHandler) RemoveLocation(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	locID, err := c.ParamsInt("loc_id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.sessions.RemoveLocation(c.Context(), id, locID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func fishingTypeIDsFromForm(c *fiber.Ctx) []int {
	raw := c.Request().PostArgs().PeekMulti("fishing_type_id")
	ids := make([]int, 0, len(raw))
	for _, b := range raw {
		if id, err := strconv.Atoi(string(b)); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func locationIDsFromForm(c *fiber.Ctx) []int {
	raw := c.Request().PostArgs().PeekMulti("location_id")
	ids := make([]int, 0, len(raw))
	for _, b := range raw {
		if id, err := strconv.Atoi(string(b)); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
