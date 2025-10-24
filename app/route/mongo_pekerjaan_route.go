package route

import (
	"crud-alumni/app/service"
	"crud-alumni/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func MongoPekerjaanRoutes(app fiber.Router) {
	mongoService := service.NewMongoPekerjaanService()

	r := app.Group("/mongo/pekerjaan")

	// Migration route (admin only)
	r.Post("/migrate", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"error":   "Only admin can perform migration",
			})
		}

		if err := service.MigratePekerjaanToMongo(); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("Migration failed: %v", err),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Data migration completed successfully",
		})
	})

	// Basic CRUD routes
	r.Get("/", middleware.AuthRequired(), mongoService.GetAllPekerjaanMongo)
	r.Get("/:id", middleware.AuthRequired(), mongoService.GetPekerjaanByIDMongo)
	r.Post("/", middleware.AuthRequired(), mongoService.CreatePekerjaanMongo)
	r.Put("/:id", middleware.AuthRequired(), mongoService.UpdatePekerjaanMongo)

	// Trash management routes
	r.Delete("/soft/:id", middleware.AuthRequired(), mongoService.SoftDeletePekerjaanMongo)
	r.Delete("/hard/:id", middleware.AuthRequired(), mongoService.HardDeletePekerjaanMongo)
	r.Get("/trash", middleware.AuthRequired(), mongoService.GetTrashMongo)
	r.Put("/restore/:id", middleware.AuthRequired(), mongoService.RestorePekerjaanMongo)

	// Search route
	r.Get("/search", middleware.AuthRequired(), mongoService.SearchPekerjaanMongo)

	// Migration route (admin only)
	r.Post("/migrate", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"error":   "Only admin can perform migration",
			})
		}

		if err := service.MigratePekerjaanToMongo(); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("Migration failed: %v", err),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Data migration completed successfully",
		})
	})
}
