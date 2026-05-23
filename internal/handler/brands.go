package handler

import (
	"fmt"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type BrandHandler struct {
	brands *sqlite.BrandsRepo
}

func (h *BrandHandler) Register(r fiber.Router) {
	r.Get("/", h.List)
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id", h.Update)
	r.Post("/:id/delete", h.Delete)
}

func (h *BrandHandler) List(c *fiber.Ctx) error {
	brands, err := h.brands.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.BrandsList(brands))
}

func (h *BrandHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/brands")
	}
	return render(c, pages.BrandFormDialog(models.Brand{}, ""))
}

func (h *BrandHandler) Create(c *fiber.Ctx) error {
	b, errMsg := brandFromForm(c)
	if errMsg != "" {
		return render(c, pages.BrandFormDialog(b, errMsg))
	}
	created, err := h.brands.Create(c.Context(), b)
	if err != nil {
		return render(c, pages.BrandFormDialog(b, err.Error()))
	}
	c.Set("HX-Retarget", "#brands-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.BrandRowPartial(*created))
}

func (h *BrandHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/brands")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	b, err := h.brands.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	return render(c, pages.BrandFormDialog(*b, ""))
}

func (h *BrandHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	b, errMsg := brandFromForm(c)
	b.ID = id
	if errMsg != "" {
		return render(c, pages.BrandFormDialog(b, errMsg))
	}
	if err := h.brands.Update(c.Context(), b); err != nil {
		return render(c, pages.BrandFormDialog(b, err.Error()))
	}
	c.Set("HX-Retarget", fmt.Sprintf("#brand-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.BrandRowPartial(b))
}

func (h *BrandHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.brands.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func brandFromForm(c *fiber.Ctx) (models.Brand, string) {
	name := c.FormValue("name")
	if name == "" {
		return models.Brand{}, "Назва — обов'язкове поле"
	}
	b := models.Brand{Name: name}
	if notes := c.FormValue("notes"); notes != "" {
		b.Notes = &notes
	}
	return b, ""
}
