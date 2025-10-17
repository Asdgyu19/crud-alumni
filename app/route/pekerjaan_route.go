package route

import (
	"crud-alumni/app/service"
	"crud-alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app fiber.Router) {
	r := app.Group("/pekerjaan")

	//  CORE PEKERJAAN ENDPOINTS
	r.Get("/", middleware.AuthRequired(), service.GetPekerjaanService)
	r.Get("/:id", middleware.AuthRequired(), service.GetPekerjaanByIDService)
	r.Post("/", middleware.AuthRequired(), service.CreatePekerjaanService)
	r.Put("/:id", middleware.AuthRequired(), service.UpdatePekerjaanService)
	r.Get("/alumni/:alumni_id", middleware.AuthRequired(), service.GetPekerjaanByAlumniService)
	r.Get("/trash", middleware.AuthRequired(), service.GetTrashService)
	r.Put("/restore/:id", middleware.AuthRequired(), service.RestoreService)
	r.Delete("/hard-delete/:id", middleware.AuthRequired(), service.HardDeleteService)
	r.Delete("/soft/:id", middleware.AuthRequired(), service.SoftDeleteService)

	// USER ROUTE (Logic di Service)
	r.Get("/my-trash", middleware.AuthRequired(), service.GetMyTrashService)
	r.Put("/my-restore/:id", middleware.AuthRequired(), service.RestoreMyService)
	r.Delete("/my-hard-delete/:id", middleware.AuthRequired(), service.HardDeleteMyService)
}
