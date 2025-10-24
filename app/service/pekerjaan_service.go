package service

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
	// Ambil parameter query
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	// Hitung offset
	offset := (page - 1) * limit

	// Validasi kolom yang boleh dipakai sorting
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

	// Validasi order (asc/desc)
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Ambil data dari repository
	pekerjaan, err := repository.GetPekerjaanWithFilter(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pekerjaan"})
	}

	// Hitung total data
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

func SoftDeletePekerjaan(pekerjaanID int, uid interface{}, role string) error {
	var userID int
	switch v := uid.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case string:
		id, _ := strconv.Atoi(v)
		userID = id
	default:
		userID = 0
	}

	return repository.SoftDeletePekerjaan(pekerjaanID, userID, role)
}

func GetTrashPekerjaan() ([]models.Pekerjaan, error) {
	return repository.GetTrashRepo()
}

func RestorePekerjaan(id int) error {
	return repository.RestoreRepo(id)
}

func HardDeletePekerjaan(id int) error {
	return repository.HardDeleteRepo(id)
}

func GetMyTrashPekerjaan(userID int) ([]models.Pekerjaan, error) {
	return repository.GetMyTrashRepo(userID)
}

func RestoreMyPekerjaan(id int, userID int) error {
	return repository.RestoreMyRepo(id, userID)
}

func HardDeleteMyPekerjaan(id int, userID int) error {
	return repository.HardDeleteMyRepo(id, userID)
}

// UNIFIED TRASH MANAGEMENT SERVICES

func GetTrashUnifiedService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var trashList []models.Pekerjaan
	var err error

	if role == "admin" {
		trashList, err = GetTrashPekerjaan()
	} else {
		trashList, err = GetMyTrashPekerjaan(userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data trash"})
	}
	return c.JSON(fiber.Map{"success": true, "data": trashList})
}

func RestoreUnifiedService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var err error
	if role == "admin" {
		err = RestorePekerjaan(id)
	} else {
		err = RestoreMyPekerjaan(id, userID)
	}

	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil di-restore"})
}

func HardDeleteUnifiedService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userIDRaw := c.Locals("user_id")

	// Handle konversi float64 ke int
	var userID int
	switch v := userIDRaw.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		return c.Status(500).JSON(fiber.Map{"error": "Invalid user ID type"})
	}

	var err error
	if role == "admin" {
		err = HardDeletePekerjaan(id)
	} else {
		err = HardDeleteMyPekerjaan(id, userID)
	}

	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanen"})
}

// SoftDeleteService - service untuk soft delete (unified untuk admin dan user)
func SoftDeleteService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userIDRaw := c.Locals("user_id")
	role := c.Locals("role").(string)

	// Handle konversi float64 ke int
	var userID int
	switch v := userIDRaw.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		return c.Status(500).JSON(fiber.Map{"error": "Invalid user ID type"})
	}

	if err := SoftDeletePekerjaan(id, userID, role); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus (soft delete)"})
}

// GetPekerjaanByIDService - service untuk get pekerjaan by ID
func GetPekerjaanByIDService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	p, err := GetPekerjaanByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": p})
}

// GetPekerjaanByAlumniService - service untuk get pekerjaan by alumni ID (ADMIN ONLY)
func GetPekerjaanByAlumniService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
	}

	alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
	list, err := GetPekerjaanByAlumni(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data"})
	}
	return c.JSON(fiber.Map{"success": true, "data": list})
}

// CreatePekerjaanService - service untuk create pekerjaan (ADMIN ONLY)
func CreatePekerjaanService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
	}

	var p models.Pekerjaan
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}
	if err := CreatePekerjaan(&p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal insert pekerjaan"})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": p})
}

// UpdatePekerjaanService - service untuk update pekerjaan (ADMIN ONLY)
func UpdatePekerjaanService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	var p models.Pekerjaan
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}
	if err := UpdatePekerjaan(id, &p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update pekerjaan"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data pekerjaan berhasil diupdate"})
}
