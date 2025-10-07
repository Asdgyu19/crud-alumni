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
