package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	user, hash, err := repository.GetUserByIdentifier(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	if !utils.CheckPassword(req.Password, hash) {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	resp := models.LoginResponse{User: user, Token: token}
	return c.JSON(fiber.Map{"success": true, "data": resp})
}

// Authenticate digunakan oleh internal service
func Authenticate(identifier, password string) (models.User, string, error) {
	user, hash, err := repository.GetUserByIdentifier(identifier)
	if err != nil {
		return models.User{}, "", err
	}
	if !utils.CheckPassword(password, hash) {
		return models.User{}, "", err
	}
	token, err := utils.GenerateToken(user)
	return user, token, err
}
