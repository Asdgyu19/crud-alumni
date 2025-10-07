package route

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"
	"crud-alumni/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app fiber.Router) {
	r := app.Group("/pekerjaan")

	// GET semua pekerjaan (sudah pakai pagination, sorting, search)
	r.Get("/", middleware.AuthRequired(), service.GetPekerjaanService)

	// GET pekerjaan by ID
	r.Get("/:id", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		p, err := service.GetPekerjaanByID(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.JSON(fiber.Map{"success": true, "data": p})
	})

	// GET pekerjaan by Alumni ID (khusus admin)
	r.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.AdminOnly(), func(c *fiber.Ctx) error {
		alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
		list, err := service.GetPekerjaanByAlumni(alumniID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data"})
		}
		return c.JSON(fiber.Map{"success": true, "data": list})
	})

	// POST pekerjaan baru
	r.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), func(c *fiber.Ctx) error {
		var p models.Pekerjaan
		if err := c.BodyParser(&p); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
		}
		if err := service.CreatePekerjaan(&p); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal insert pekerjaan"})
		}
		return c.Status(201).JSON(fiber.Map{"success": true, "data": p})
	})

	// UPDATE pekerjaan
	r.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var p models.Pekerjaan
		if err := c.BodyParser(&p); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
		}
		if err := service.UpdatePekerjaan(id, &p); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update pekerjaan"})
		}
		return c.JSON(fiber.Map{"success": true, "message": "Data pekerjaan diupdate"})
	})

	// DELETE pekerjaan
	r.Delete("/:id", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))

		// ambil info user dari JWT middleware
		userID := c.Locals("user_id").(int)
		role := c.Locals("role").(string)

		if err := service.SoftDeletePekerjaan(id, userID, role); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus (soft delete)"})
	})

}
