package route

import (
	"crud-alumni/app/service"
	"crud-alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoutes(app fiber.Router) {
	r := app.Group("/alumni")
	alumniService := service.NewAlumniService()

	// GET routes
	r.Get("/", middleware.AuthRequired(), alumniService.GetAllAlumni)
	r.Get("/:id", middleware.AuthRequired(), alumniService.GetAlumniByID)
	r.Get("/tahun/:tahun", middleware.AuthRequired(), alumniService.GetAlumniByTahunAndGaji)

	// POST routes (Admin only)
	r.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), alumniService.CreateAlumni)

	// PUT routes (Admin only)
	r.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), alumniService.UpdateAlumni)

	// DELETE routes (Admin only)
	r.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), alumniService.DeleteAlumni)
}
