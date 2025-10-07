package middleware

import "github.com/gofiber/fiber/v2"

func RoleAllowed(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role").(string)
        for _, r := range roles {
            if role == r {
                return c.Next()
            }
        }
        return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
    }
}
