package handler

import (
	"fmt"
	"strconv"

	"perch/internal/models"
	"perch/internal/repository/sqlite"
	"perch/internal/templates/pages"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	locations   *sqlite.LocationsRepo
	waterBodies *sqlite.WaterBodiesRepo
}

func (h *LocationHandler) Register(r fiber.Router) {
	r.Get("/", h.List)
	r.Get("/new", h.New)
	r.Post("/", h.Create)
	r.Get("/:id/edit", h.Edit)
	r.Post("/:id", h.Update)
	r.Post("/:id/delete", h.Delete)
}

func (h *LocationHandler) List(c *fiber.Ctx) error {
	locs, err := h.locations.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LocationsList(locs))
}

func (h *LocationHandler) New(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/locations")
	}
	wbs, err := h.waterBodies.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LocationFormDialog(models.Location{}, wbs, ""))
}

func (h *LocationHandler) Create(c *fiber.Ctx) error {
	loc, errMsg := locationFromForm(c)
	if errMsg != "" {
		wbs, _ := h.waterBodies.List(c.Context())
		return render(c, pages.LocationFormDialog(loc, wbs, errMsg))
	}
	created, err := h.locations.Create(c.Context(), loc)
	if err != nil {
		wbs, _ := h.waterBodies.List(c.Context())
		return render(c, pages.LocationFormDialog(loc, wbs, err.Error()))
	}
	full, err := h.locations.Get(c.Context(), created.ID)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", "#locations-tbody")
	c.Set("HX-Reswap", "afterbegin")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.LocationRowPartial(*full))
}

func (h *LocationHandler) Edit(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" {
		return c.Redirect("/locations")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	loc, err := h.locations.Get(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	wbs, err := h.waterBodies.List(c.Context())
	if err != nil {
		return err
	}
	return render(c, pages.LocationFormDialog(*loc, wbs, ""))
}

func (h *LocationHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	loc, errMsg := locationFromForm(c)
	loc.ID = id
	if errMsg != "" {
		wbs, _ := h.waterBodies.List(c.Context())
		return render(c, pages.LocationFormDialog(loc, wbs, errMsg))
	}
	if err := h.locations.Update(c.Context(), loc); err != nil {
		wbs, _ := h.waterBodies.List(c.Context())
		return render(c, pages.LocationFormDialog(loc, wbs, err.Error()))
	}
	full, err := h.locations.Get(c.Context(), id)
	if err != nil {
		return err
	}
	c.Set("HX-Retarget", fmt.Sprintf("#loc-%d", id))
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger", "closeDialog")
	return render(c, pages.LocationRowPartial(*full))
}

func (h *LocationHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.locations.Delete(c.Context(), id); err != nil {
		return err
	}
	if c.Get("HX-Request") != "" {
		return c.SendStatus(fiber.StatusOK)
	}
	return c.Redirect("/locations")
}

func locationFromForm(c *fiber.Ctx) (models.Location, string) {
	name := c.FormValue("name")
	region := c.FormValue("region")
	wbID, _ := strconv.Atoi(c.FormValue("waterbody_id"))

	loc := models.Location{Name: name, Region: region, WaterBodyID: wbID}
	if notes := c.FormValue("notes"); notes != "" {
		loc.Notes = &notes
	}

	if name == "" || region == "" || wbID == 0 {
		return loc, "Назва, регіон та водойма — обов'язкові поля"
	}
	return loc, ""
}
