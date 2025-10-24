package route

import (
	"crud-alumni/app/service"
	"crud-alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(router fiber.Router) {
	fileService := service.NewFileService()

	files := router.Group("/files")
	files.Use(middleware.AuthRequired())

	files.Post("/upload", fileService.UploadFile)
	files.Get("/", fileService.GetAllFiles)
	files.Get("/:id", fileService.GetFileByID)
	files.Delete("/:id", fileService.DeleteFile)
}



