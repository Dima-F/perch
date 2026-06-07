package handler

import (
	"fmt"
	"strconv"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type LureHandler struct {
	lures      *sqlite.LuresRepo
	lureModels *sqlite.LureModelsRepo
}

func (h *LureHandler) Register(r fiber.Router) {
	r.Get("/", h.List)
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id", h.Update)
	r.Post("/:id/delete", h.Delete)
}

func (h *LureHandler) List(c *fiber.Ctx) error {
	lures, err := h.lures.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LuresList(lures))
}

func (h *LureHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/lures")
	}
	lms, err := h.lureModels.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LureFormDialog(models.Lure{}, lms, ""))
}

func (h *LureHandler) Create(c *fiber.Ctx) error {
	lu, errMsg := lureFromForm(c)
	if errMsg != "" {
		lms, _ := h.lureModels.List(c.Context())
		return render(c, pages.LureFormDialog(lu, lms, errMsg))
	}
	created, err := h.lures.Create(c.Context(), lu)
	if err != nil {
		lms, _ := h.lureModels.List(c.Context())
		return render(c, pages.LureFormDialog(lu, lms, err.Error()))
	}
	full, err := h.lures.Get(c.Context(), created.ID)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", "#lures-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.LureRowPartial(*full))
}

func (h *LureHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/lures")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	lu, err := h.lures.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	lms, err := h.lureModels.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LureFormDialog(*lu, lms, ""))
}

func (h *LureHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	lu, errMsg := lureFromForm(c)
	lu.ID = id
	if errMsg != "" {
		lms, _ := h.lureModels.List(c.Context())
		return render(c, pages.LureFormDialog(lu, lms, errMsg))
	}
	if err := h.lures.Update(c.Context(), lu); err != nil {
		lms, _ := h.lureModels.List(c.Context())
		return render(c, pages.LureFormDialog(lu, lms, err.Error()))
	}
	full, err := h.lures.Get(c.Context(), id)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", fmt.Sprintf("#lure-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.LureRowPartial(*full))
}

func (h *LureHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.lures.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func lureFromForm(c *fiber.Ctx) (models.Lure, string) {
	modelID, _ := strconv.Atoi(c.FormValue("model_id"))
	if modelID == 0 {
		return models.Lure{}, "Модель — обов'язкове поле"
	}
	lu := models.Lure{ModelID: modelID}
	if color := c.FormValue("color"); color != "" {
		lu.Color = &color
	}
	if size := c.FormValue("size"); size != "" {
		lu.Size = &size
	}
	if weightStr := c.FormValue("weight_g"); weightStr != "" {
		if w, err := strconv.ParseFloat(weightStr, 64); err == nil {
			lu.WeightG = &w
		}
	}
	if notes := c.FormValue("notes"); notes != "" {
		lu.Notes = &notes
	}
	return lu, ""
}
