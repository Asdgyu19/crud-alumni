package middleware

import (
	"crud-alumni/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token akses diperlukan"})
		}
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Format token tidak valid"})
		}
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid atau expired"})
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		r, ok := c.Locals("role").(string)
		if !ok || r != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
		}
		return c.Next()
	}
}
