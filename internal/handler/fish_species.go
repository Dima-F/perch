package handler

import (
	"fmt"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type FishSpeciesHandler struct {
	fishSpecies *sqlite.FishSpeciesRepo
}

func (h *FishSpeciesHandler) Register(r fiber.Router) {
	r.Get("/", h.List)
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id", h.Update)
	r.Post("/:id/delete", h.Delete)
}

func (h *FishSpeciesHandler) List(c *fiber.Ctx) error {
	species, err := h.fishSpecies.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.FishSpeciesList(species))
}

func (h *FishSpeciesHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/fish-species")
	}
	return render(c, pages.FishSpeciesFormDialog(models.FishSpecies{}, ""))
}

func (h *FishSpeciesHandler) Create(c *fiber.Ctx) error {
	f, errMsg := fishFromForm(c)
	if errMsg != "" {
		return render(c, pages.FishSpeciesFormDialog(f, errMsg))
	}
	created, err := h.fishSpecies.Create(c.Context(), f)
	if err != nil {
		return render(c, pages.FishSpeciesFormDialog(f, err.Error()))
	}
	c.Set("HX-Retarget", "#fish-species-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.FishSpeciesRowPartial(*created))
}

func (h *FishSpeciesHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/fish-species")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	f, err := h.fishSpecies.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	return render(c, pages.FishSpeciesFormDialog(*f, ""))
}

func (h *FishSpeciesHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	f, errMsg := fishFromForm(c)
	f.ID = id
	if errMsg != "" {
		return render(c, pages.FishSpeciesFormDialog(f, errMsg))
	}
	if err := h.fishSpecies.Update(c.Context(), f); err != nil {
		return render(c, pages.FishSpeciesFormDialog(f, err.Error()))
	}
	c.Set("HX-Retarget", fmt.Sprintf("#fish-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.FishSpeciesRowPartial(f))
}

func (h *FishSpeciesHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.fishSpecies.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func fishFromForm(c *fiber.Ctx) (models.FishSpecies, string) {
	name := c.FormValue("name")
	if name == "" {
		return models.FishSpecies{}, "Назва — обов'язкове поле"
	}
	return models.FishSpecies{Name: name}, ""
}
