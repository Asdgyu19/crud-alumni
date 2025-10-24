package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository/mongo"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MongoPekerjaanService struct {
	repo *mongo.PekerjaanRepository
}

func NewMongoPekerjaanService() *MongoPekerjaanService {
	return &MongoPekerjaanService{
		repo: mongo.NewPekerjaanRepository(),
	}
}

// GetAllPekerjaanMongo handles GET /mongo/pekerjaan
func (s *MongoPekerjaanService) GetAllPekerjaanMongo(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Calculate offset
	offset := int64((page - 1) * limit)

	pekerjaan, total, err := s.repo.FindAll(int64(limit), offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch data",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetPekerjaanByIDMongo handles GET /mongo/pekerjaan/:id
func (s *MongoPekerjaanService) GetPekerjaanByIDMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	pekerjaan, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Pekerjaan not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
	})
}

// CreatePekerjaanMongo handles POST /mongo/pekerjaan
func (s *MongoPekerjaanService) CreatePekerjaanMongo(c *fiber.Ctx) error {
	var pekerjaan models.Pekerjaan

	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := s.repo.Create(&pekerjaan); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create pekerjaan",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
		"message": "Pekerjaan created successfully",
	})
}

// UpdatePekerjaanMongo handles PUT /mongo/pekerjaan/:id
func (s *MongoPekerjaanService) UpdatePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")
	var pekerjaan models.Pekerjaan

	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := s.repo.Update(id, &pekerjaan); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update pekerjaan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan updated successfully",
	})
}

// SoftDeletePekerjaanMongo handles DELETE /mongo/pekerjaan/soft/:id
func (s *MongoPekerjaanService) SoftDeletePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(float64)

	if err := s.repo.SoftDelete(id, int(userID)); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to soft delete pekerjaan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan soft deleted successfully",
	})
}

// HardDeletePekerjaanMongo handles DELETE /mongo/pekerjaan/hard/:id
func (s *MongoPekerjaanService) HardDeletePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.HardDelete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to hard delete pekerjaan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan permanently deleted",
	})
}

// GetTrashMongo handles GET /mongo/pekerjaan/trash
func (s *MongoPekerjaanService) GetTrashMongo(c *fiber.Ctx) error {
	pekerjaan, err := s.repo.GetTrash()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch trash",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
	})
}

// RestorePekerjaanMongo handles PUT /mongo/pekerjaan/restore/:id
func (s *MongoPekerjaanService) RestorePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.Restore(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to restore pekerjaan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan restored successfully",
	})
}

// SearchPekerjaanMongo handles GET /mongo/pekerjaan/search
func (s *MongoPekerjaanService) SearchPekerjaanMongo(c *fiber.Ctx) error {
	keyword := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	offset := int64((page - 1) * limit)

	pekerjaan, total, err := s.repo.Search(keyword, int64(limit), offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to search pekerjaan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
		"meta": fiber.Map{
			"page":    page,
			"limit":   limit,
			"total":   total,
			"pages":   (total + int64(limit) - 1) / int64(limit),
			"keyword": keyword,
		},
	})
}
