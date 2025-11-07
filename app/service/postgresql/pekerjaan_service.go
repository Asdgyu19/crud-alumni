package postgresql

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaan() ([]models.Pekerjaan, error) {
	return repository.GetAllPekerjaanRepo()
}

func GetPekerjaanByID(id int) (models.Pekerjaan, error) {
	return repository.GetPekerjaanByIDRepo(id)
}

func GetPekerjaanByAlumni(alumniID int) ([]models.Pekerjaan, error) {
	return repository.GetPekerjaanByAlumniRepo(alumniID)
}

func CreatePekerjaan(p *models.Pekerjaan) error {
	return repository.CreatePekerjaanRepo(p)
}

func UpdatePekerjaan(id int, p *models.Pekerjaan) error {
	return repository.UpdatePekerjaanRepo(id, p)
}

func DeletePekerjaan(id int) error {
	return repository.DeletePekerjaanRepo(id)
}

// GetPekerjaanService -> ambil data pekerjaan alumni dengan pagination, sorting, dan search
func GetPekerjaanService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{
		"id":                  true,
		"nama_perusahaan":     true,
		"posisi_jabatan":      true,
		"bidang_industri":     true,
		"lokasi_kerja":        true,
		"gaji_range":          true,
		"tanggal_mulai_kerja": true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	pekerjaan, err := repository.GetPekerjaanWithFilter(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pekerjaan"})
	}

	total, err := repository.CountPekerjaan(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghitung data pekerjaan"})
	}

	response := models.PekerjaanResponse{
		Data: pekerjaan,
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

// ... [other functions remain the same]
