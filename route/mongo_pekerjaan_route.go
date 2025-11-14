package route

import (
	"crud-alumni/app/service/mongo"
	"crud-alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func MongoPekerjaanRoutes(app fiber.Router) {
	mongoService := mongo.NewPekerjaanService()
	r := app.Group("/mongo/pekerjaan")

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
}
