package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/app/repository/mongo"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	useMongoDb   bool
	postgresRepo repository.AlumniRepository
	mongoRepo    *mongo.AlumniRepository
}

func NewAlumniService() *AlumniService {
	s := &AlumniService{
		useMongoDb: os.Getenv("USE_MONGODB") == "true",
	}

	if s.useMongoDb {
		s.mongoRepo = mongo.NewAlumniRepository()
	}

	return s
}

// Handler methods
func (s *AlumniService) GetAllAlumni(c *fiber.Ctx) error {
	var alumni []models.Alumni
	var err error

	if s.useMongoDb {
		alumni, err = s.mongoRepo.FindAll()
	} else {
		alumni, err = repository.GetAllAlumniRepo()
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data"})
	}
	return c.JSON(fiber.Map{"success": true, "data": alumni})
}

func (s *AlumniService) GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var alumni *models.Alumni
	if s.useMongoDb {
		alumni, err = s.mongoRepo.FindByID(strconv.Itoa(id))
	} else {
		var a models.Alumni
		a, err = repository.GetAlumniByIDRepo(id)
		alumni = &a
	}

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": alumni})
}

func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	var err error
	if s.useMongoDb {
		err = s.mongoRepo.Create(&alumni)
	} else {
		err = repository.CreateAlumniRepo(&alumni)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal insert data"})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": alumni})
}

func (s *AlumniService) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	if s.useMongoDb {
		err = s.mongoRepo.Update(strconv.Itoa(id), &alumni)
	} else {
		err = repository.UpdateAlumniRepo(id, &alumni)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data alumni diupdate"})
}

func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	if s.useMongoDb {
		err = s.mongoRepo.Delete(strconv.Itoa(id))
	} else {
		err = repository.DeleteAlumniRepo(id)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Alumni dihapus"})
}

func (s *AlumniService) GetAlumniByTahunAndGaji(c *fiber.Ctx) error {
	tahun, err := strconv.Atoi(c.Params("tahun"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Parameter tahun harus angka"})
	}

	data, err := repository.GetAlumniByTahunAndGajiRepo(tahun)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(data),
		"data":    data,
	})
}
