package handler

import (
	"fmt"
	"strconv"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type LureModelHandler struct {
	lureModels *sqlite.LureModelsRepo
	lureTypes  *sqlite.LureTypesRepo
	brands     *sqlite.BrandsRepo
}

func (h *LureModelHandler) Register(r fiber.Router) {
	r.Get("/", h.List)
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id", h.Update)
	r.Post("/:id/delete", h.Delete)
}

func (h *LureModelHandler) List(c *fiber.Ctx) error {
	lms, err := h.lureModels.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LureModelsList(lms))
}

func (h *LureModelHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/lure-models")
	}
	lts, err := h.lureTypes.List(c.Context())
	if err != nil {
		return err
	}
	brands, err := h.brands.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LureModelFormDialog(models.LureModel{}, lts, brands, ""))
}

func (h *LureModelHandler) Create(c *fiber.Ctx) error {
	m, errMsg := lureModelFromForm(c)
	if errMsg != "" {
		lts, _ := h.lureTypes.List(c.Context())
		brands, _ := h.brands.List(c.Context())
		return render(c, pages.LureModelFormDialog(m, lts, brands, errMsg))
	}
	created, err := h.lureModels.Create(c.Context(), m)
	if err != nil {
		lts, _ := h.lureTypes.List(c.Context())
		brands, _ := h.brands.List(c.Context())
		return render(c, pages.LureModelFormDialog(m, lts, brands, err.Error()))
	}
	full, err := h.lureModels.Get(c.Context(), created.ID)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", "#lure-models-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.LureModelRowPartial(*full))
}

func (h *LureModelHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/lure-models")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	m, err := h.lureModels.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	lts, err := h.lureTypes.List(c.Context())
	if err != nil {
		return err
	}
	brands, err := h.brands.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LureModelFormDialog(*m, lts, brands, ""))
}

func (h *LureModelHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	m, errMsg := lureModelFromForm(c)
	m.ID = id
	if errMsg != "" {
		lts, _ := h.lureTypes.List(c.Context())
		brands, _ := h.brands.List(c.Context())
		return render(c, pages.LureModelFormDialog(m, lts, brands, errMsg))
	}
	if err := h.lureModels.Update(c.Context(), m); err != nil {
		lts, _ := h.lureTypes.List(c.Context())
		brands, _ := h.brands.List(c.Context())
		return render(c, pages.LureModelFormDialog(m, lts, brands, err.Error()))
	}
	full, err := h.lureModels.Get(c.Context(), id)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", fmt.Sprintf("#lure-model-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.LureModelRowPartial(*full))
}

func (h *LureModelHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.lureModels.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func lureModelFromForm(c *fiber.Ctx) (models.LureModel, string) {
	name := c.FormValue("name")
	lureTypeID, _ := strconv.Atoi(c.FormValue("luretype_id"))
	if name == "" || lureTypeID == 0 {
		return models.LureModel{Name: name, LureTypeID: lureTypeID}, "Назва та тип — обов'язкові поля"
	}
	m := models.LureModel{Name: name, LureTypeID: lureTypeID}
	if bidStr := c.FormValue("brand_id"); bidStr != "" {
		if bid, err := strconv.Atoi(bidStr); err == nil && bid > 0 {
			m.BrandID = &bid
		}
	}
	if notes := c.FormValue("notes"); notes != "" {
		m.Notes = &notes
	}
	return m, ""
}
