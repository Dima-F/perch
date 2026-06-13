package handler

import (
	"fmt"
	"strconv"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type CatchHandler struct {
	catches     *sqlite.CatchesRepo
	sessions    *sqlite.SessionsRepo
	locations   *sqlite.LocationsRepo
	lures       *sqlite.LuresRepo
	fishSpecies *sqlite.FishSpeciesRepo
}

func (h *CatchHandler) Register(r fiber.Router) {
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id/delete", h.Delete)
	r.Post("/:id", h.Update)
}

func (h *CatchHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/")
	}
	sessionID, err := strconv.Atoi(c.Query("session_id"))
	if err != nil || sessionID == 0 {
		return fiber.ErrBadRequest
	}
	fish, err := h.fishSpecies.List(c.Context())
	if err != nil {
		return err
	}
	lures, err := h.lures.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.CatchFormDialog(models.Catch{SessionID: sessionID}, fish, lures, ""))
}

func (h *CatchHandler) Create(c *fiber.Ctx) error {
	catch, errMsg := catchFromForm(c)
	fish, _ := h.fishSpecies.List(c.Context())
	lures, _ := h.lures.List(c.Context())
	if errMsg != "" {
		return render(c, pages.CatchFormDialog(catch, fish, lures, errMsg))
	}
	created, err := h.catches.Create(c.Context(), catch)
	if err != nil {
		return render(c, pages.CatchFormDialog(catch, fish, lures, err.Error()))
	}
	full, err := h.catches.Get(c.Context(), created.ID)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", "#catches-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeCatchDialog")
	return render(c, pages.CatchRowPartial(*full))
}

func (h *CatchHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	catch, err := h.catches.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	fish, err := h.fishSpecies.List(c.Context())
	if err != nil {
		return err
	}
	lures, err := h.lures.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.CatchFormDialog(*catch, fish, lures, ""))
}

func (h *CatchHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	catch, errMsg := catchFromForm(c)
	catch.ID = id
	fish, _ := h.fishSpecies.List(c.Context())
	lures, _ := h.lures.List(c.Context())
	if errMsg != "" {
		return render(c, pages.CatchFormDialog(catch, fish, lures, errMsg))
	}
	if err := h.catches.Update(c.Context(), catch); err != nil {
		return render(c, pages.CatchFormDialog(catch, fish, lures, err.Error()))
	}
	full, err := h.catches.Get(c.Context(), id)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", fmt.Sprintf("#catch-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeCatchDialog")
	return render(c, pages.CatchRowPartial(*full))
}

func (h *CatchHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.catches.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func catchFromForm(c *fiber.Ctx) (models.Catch, string) {
	var catch models.Catch

	if sidStr := c.FormValue("session_id"); sidStr != "" {
		if sid, err := strconv.Atoi(sidStr); err == nil {
			catch.SessionID = sid
		}
	}

	fishIDStr := c.FormValue("fish_id")
	if fishIDStr == "" {
		return catch, "Вид риби — обов'язкове поле"
	}
	fishID, err := strconv.Atoi(fishIDStr)
	if err != nil || fishID == 0 {
		return catch, "Невірний вид риби"
	}
	catch.FishID = fishID

	if lureIDStr := c.FormValue("lure_id"); lureIDStr != "" {
		if lureID, err := strconv.Atoi(lureIDStr); err == nil && lureID > 0 {
			catch.LureID = &lureID
		}
	}

	countStr := c.FormValue("count")
	if countStr == "" {
		return catch, "Кількість — обов'язкове поле"
	}
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		return catch, "Кількість має бути більше 0"
	}
	catch.Count = count

	if s := c.FormValue("avg_length_cm"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			catch.AvgLengthCm = &v
		}
	}
	if s := c.FormValue("max_length_cm"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			catch.MaxLengthCm = &v
		}
	}
	if s := c.FormValue("weight_g"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			catch.WeightG = &v
		}
	}
	if s := c.FormValue("jig_setup"); s != "" {
		catch.JigSetup = &s
	}
	if s := c.FormValue("jig_weight_g"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			catch.JigWeightG = &v
		}
	}
	if s := c.FormValue("notes"); s != "" {
		catch.Notes = &s
	}

	return catch, ""
}
