package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/database"
	"errors"
	"fmt"
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

	// Bentuk response
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
	// casting user_id dari token (c.Locals)
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

	fmt.Println("SoftDelete Debug -> pekerjaanID:", pekerjaanID, "userID:", userID, "role:", role)

	// kalau role = user biasa → cek kepemilikan
	if role == "user" {
		var count int
		err := database.DB.QueryRow(`
			SELECT COUNT(*)
			FROM pekerjaan_alumni pa
			JOIN alumni a ON pa.alumni_id = a.id
			WHERE pa.id = $1 
			  AND a.user_id = $2 
			  AND pa.is_delete = false
		`, pekerjaanID, userID).Scan(&count)

		if err != nil {
			return err
		}

		fmt.Println("SoftDelete Debug -> hasil cek kepemilikan count:", count)

		if count == 0 {
			return errors.New("tidak diizinkan hapus pekerjaan orang lain")
		}
	}

	// lakukan soft delete
	return repository.SoftDeletePekerjaanRepo(pekerjaanID, userID)
}

func SoftDeletePekerjaanRepo(id int, userRole string, userID int) error {
	if userRole == "admin" {
		// admin boleh hapus siapa saja
		_, err := database.DB.Exec(`
            UPDATE pekerjaan_alumni SET deleted_at=NOW() WHERE id=$1
        `, id)
		return err
	} else {
		// user hanya boleh hapus miliknya sendiri
		_, err := database.DB.Exec(`
            UPDATE pekerjaan_alumni 
            SET deleted_at=NOW() 
            WHERE id=$1 AND alumni_id=$2
        `, id, userID)
		return err
	}
}

// GetTrashPekerjaan - ambil semua pekerjaan yang sudah di-soft delete (HANYA ADMIN YES)
func GetTrashPekerjaan() ([]models.Pekerjaan, error) {
	return repository.GetTrashRepo()
}

// RestorePekerjaan - restore pekerjaan dari trash (admin only)
func RestorePekerjaan(id int) error {
	return repository.RestoreRepo(id)
}

// HardDeletePekerjaan - hapus permanent dari database (admin only)
func HardDeletePekerjaan(id int) error {
	return repository.HardDeleteRepo(id)
}

// GetMyTrashPekerjaan - user lihat trash milik sendiri
func GetMyTrashPekerjaan(userID int) ([]models.Pekerjaan, error) {
	return repository.GetMyTrashRepo(userID)
}

// RestoreMyPekerjaan - user restore pekerjaan milik sendiri
func RestoreMyPekerjaan(id int, userID int) error {
	return repository.RestoreMyRepo(id, userID)
}

// HardDeleteMyPekerjaan - user hard delete pekerjaan milik sendiri
func HardDeleteMyPekerjaan(id int, userID int) error {
	return repository.HardDeleteMyRepo(id, userID)
}

// ================================================================
// 🗑️ TRASH MANAGEMENT SERVICES
// ================================================================

// GetTrashService - service untuk lihat semua trash (ADMIN ONLY)
func GetTrashService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
	}

	trashList, err := GetTrashPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data trash"})
	}
	return c.JSON(fiber.Map{"success": true, "data": trashList})
}

// RestoreService - service untuk restore dari trash (ADMIN ONLY)
func RestoreService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	if err := RestorePekerjaan(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil di-restore"})
}

// HardDeleteService - service untuk hard delete (ADMIN ONLY)
func HardDeleteService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	if err := HardDeletePekerjaan(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanent"})
}

// SoftDeleteService - service untuk soft delete (user/admin)
func SoftDeleteService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	if err := SoftDeletePekerjaan(id, userID, role); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus (soft delete)"})
}

// GetMyTrashService - user service untuk lihat trash milik sendiri
func GetMyTrashService(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	trashList, err := GetMyTrashPekerjaan(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data trash"})
	}
	return c.JSON(fiber.Map{"success": true, "data": trashList})
}

// RestoreMyService - user service untuk restore milik sendiri
func RestoreMyService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)

	if err := RestoreMyPekerjaan(id, userID); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil di-restore"})
}

// HardDeleteMyService - user service untuk hard delete milik sendiri
func HardDeleteMyService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)

	if err := HardDeleteMyPekerjaan(id, userID); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanen"})
}

// ================================================================
// 📋 ADDITIONAL CLEAN SERVICE FUNCTIONS
// ================================================================

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
