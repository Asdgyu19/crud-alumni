package route

import (
	"crud-alumni/app/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router) {
	authService := service.NewAuthService()
	app.Post("/login", authService.Login)
}
