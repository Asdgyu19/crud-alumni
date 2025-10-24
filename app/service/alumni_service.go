package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/app/repository/mongo"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	useMongoDb   bool
	postgresRepo repository.AlumniRepository
	mongoRepo    *mongo.AlumniRepository
)

func init() {
	// Check if we should use MongoDB
	useMongoDb = os.Getenv("USE_MONGODB") == "true"

	if useMongoDb {
		mongoRepo = mongo.NewAlumniRepository()
	}
}

// PostgreSQL functions
func GetAllAlumni() ([]models.Alumni, error) {
	if useMongoDb {
		return mongoRepo.FindAll()
	}
	// Original PostgreSQL code
	return repository.GetAllAlumniRepo()
}

func GetAlumniByID(id interface{}) (*models.Alumni, error) {
	if useMongoDb {
		stringID, ok := id.(string)
		if !ok {
			numID, ok := id.(int)
			if ok {
				stringID = strconv.Itoa(numID)
			}
		}
		return mongoRepo.FindByID(stringID)
	}
	// Original PostgreSQL code
	numID, ok := id.(int)
	if !ok {
		strID, ok := id.(string)
		if ok {
			numID, _ = strconv.Atoi(strID)
		}
	}
	alumni, err := repository.GetAlumniByIDRepo(numID)
	return &alumni, err
}

func CreateAlumni(a *models.Alumni) error {
	if useMongoDb {
		return mongoRepo.Create(a)
	}
	// Original PostgreSQL code
	return repository.CreateAlumniRepo(a)
}

func UpdateAlumni(id interface{}, a *models.Alumni) error {
	if useMongoDb {
		stringID, ok := id.(string)
		if !ok {
			numID, ok := id.(int)
			if ok {
				stringID = strconv.Itoa(numID)
			}
		}
		return mongoRepo.Update(stringID, a)
	}
	// Original PostgreSQL code
	numID, ok := id.(int)
	if !ok {
		strID, ok := id.(string)
		if ok {
			numID, _ = strconv.Atoi(strID)
		}
	}
	return repository.UpdateAlumniRepo(numID, a)
}

func DeleteAlumni(id interface{}) error {
	if useMongoDb {
		stringID, ok := id.(string)
		if !ok {
			numID, ok := id.(int)
			if ok {
				stringID = strconv.Itoa(numID)
			}
		}
		return mongoRepo.Delete(stringID)
	}
	// Original PostgreSQL code
	numID, ok := id.(int)
	if !ok {
		strID, ok := id.(string)
		if ok {
			numID, _ = strconv.Atoi(strID)
		}
	}
	return repository.DeleteAlumniRepo(numID)
}

func GetAlumniService(c *fiber.Ctx) error {
	if useMongoDb {
		// MongoDB implementation
		alumni, err := mongoRepo.FindAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch alumni"})
		}
		return c.JSON(fiber.Map{
			"data": alumni,
			"meta": fiber.Map{
				"total": len(alumni),
			},
		})
	}

	// Original PostgreSQL implementation
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	whitelist := map[string]bool{"id": true, "nama": true, "angkatan": true, "tahun_lulus": true}
	if !whitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	data, err := repository.GetAlumniWithFilter(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch alumni"})
	}

	total, _ := repository.CountAlumni(search)

	response := models.AlumniResponse{
		Data: data,
		Meta: models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}
