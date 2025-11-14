package route

import (
	"crud-alumni/app/service"
	"crud-alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app fiber.Router) {
	r := app.Group("/pekerjaan")

	//  PEKERJAAN
	r.Get("/", middleware.AuthRequired(), service.GetPekerjaanService)
	r.Get("/:id", middleware.AuthRequired(), service.GetPekerjaanByIDService)
	r.Post("/", middleware.AuthRequired(), service.CreatePekerjaanService)
	r.Put("/:id", middleware.AuthRequired(), service.UpdatePekerjaanService)
	r.Get("/alumni/:alumni_id", middleware.AuthRequired(), service.GetPekerjaanByAlumniService)
	r.Get("/trash", middleware.AuthRequired(), service.GetTrashUnifiedService)
	r.Put("/restore/:id", middleware.AuthRequired(), service.RestoreUnifiedService)
	r.Delete("/hard-delete/:id", middleware.AuthRequired(), service.HardDeleteUnifiedService)
	r.Delete("/soft/:id", middleware.AuthRequired(), service.SoftDeleteService)
}
