package route

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router) {
	app.Post("/login", func(c *fiber.Ctx) error {
		var req models.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
		}
		user, token, err := service.Authenticate(req.Username, req.Password)
		if err != nil || token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
		}
		resp := models.LoginResponse{User: user, Token: token}
		return c.JSON(fiber.Map{"success": true, "data": resp})
	})
}
