package handler

import (
	"fmt"
	"strconv"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type WaterBodyHandler struct {
	waterBodies     *sqlite.WaterBodiesRepo
	waterBodyTypes  *sqlite.WaterBodyTypesRepo
}

func (h *WaterBodyHandler) Register(r fiber.Router) {
	r.Get("/", h.List)
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id", h.Update)
	r.Post("/:id/delete", h.Delete)
}

func (h *WaterBodyHandler) List(c *fiber.Ctx) error {
	wbs, err := h.waterBodies.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.WaterBodiesList(wbs))
}

func (h *WaterBodyHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/water-bodies")
	}
	types, err := h.waterBodyTypes.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.WaterBodyFormDialog(models.WaterBody{}, types, ""))
}

func (h *WaterBodyHandler) Create(c *fiber.Ctx) error {
	wb, errMsg := waterBodyFromForm(c)
	if errMsg != "" {
		types, _ := h.waterBodyTypes.List(c.Context())
		return render(c, pages.WaterBodyFormDialog(wb, types, errMsg))
	}
	created, err := h.waterBodies.Create(c.Context(), wb)
	if err != nil {
		types, _ := h.waterBodyTypes.List(c.Context())
		return render(c, pages.WaterBodyFormDialog(wb, types, err.Error()))
	}
	full, err := h.waterBodies.Get(c.Context(), created.ID)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", "#water-bodies-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.WaterBodyRowPartial(*full))
}

func (h *WaterBodyHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/water-bodies")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	wb, err := h.waterBodies.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	types, err := h.waterBodyTypes.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.WaterBodyFormDialog(*wb, types, ""))
}

func (h *WaterBodyHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	wb, errMsg := waterBodyFromForm(c)
	wb.ID = id
	if errMsg != "" {
		types, _ := h.waterBodyTypes.List(c.Context())
		return render(c, pages.WaterBodyFormDialog(wb, types, errMsg))
	}
	if err := h.waterBodies.Update(c.Context(), wb); err != nil {
		types, _ := h.waterBodyTypes.List(c.Context())
		return render(c, pages.WaterBodyFormDialog(wb, types, err.Error()))
	}
	full, err := h.waterBodies.Get(c.Context(), id)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", fmt.Sprintf("#wb-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.WaterBodyRowPartial(*full))
}

func (h *WaterBodyHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.waterBodies.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func waterBodyFromForm(c *fiber.Ctx) (models.WaterBody, string) {
	name := c.FormValue("name")
	typeID, _ := strconv.Atoi(c.FormValue("water_body_type_id"))
	if name == "" || typeID == 0 {
		return models.WaterBody{Name: name, WaterBodyTypeID: typeID}, "Назва та тип — обов'язкові поля"
	}
	return models.WaterBody{Name: name, WaterBodyTypeID: typeID}, ""
}
