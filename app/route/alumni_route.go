package route

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"
	"crud-alumni/middleware"
	"crud-alumni/database"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AlumniRoutes(app fiber.Router) {
	r := app.Group("/alumni")

	// GET semua alumni
	r.Get("/", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		alumni, err := service.GetAllAlumni()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data"})
		}
		return c.JSON(fiber.Map{"success": true, "data": alumni})
	})

	// GET alumni by ID
	r.Get("/:id", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		a, err := service.GetAlumniByID(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
		}
		return c.JSON(fiber.Map{"success": true, "data": a})
	})

	// POST alumni baru (Admin only)
	r.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), func(c *fiber.Ctx) error {
		var a models.Alumni
		if err := c.BodyParser(&a); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
		}
		if err := service.CreateAlumni(&a); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal insert data"})
		}
		return c.Status(201).JSON(fiber.Map{"success": true, "data": a})
	})

	// UPDATE alumni (Admin only)
	r.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var a models.Alumni
		if err := c.BodyParser(&a); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
		}
		if err := service.UpdateAlumni(id, &a); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update"})
		}
		return c.JSON(fiber.Map{"success": true, "message": "Data alumni diupdate"})
	})

	// DELETE alumni (Admin only)
	r.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		if err := service.DeleteAlumni(id); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus"})
		}
		return c.JSON(fiber.Map{"success": true, "message": "Alumni dihapus"})
	})

	// GET alumni by tahun lulus & gaji >= 4jt
r.Get("/tahun/:tahun", middleware.AuthRequired(), func(c *fiber.Ctx) error {
    // konversi param tahun ke integer
    tahun, err := strconv.Atoi(c.Params("tahun"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Parameter tahun harus angka"})
    }

    rows, err := database.DB.Query(`
        SELECT a.id, a.nama, a.jurusan, a.tahun_lulus,
               p.bidang_industri, p.nama_perusahaan, p.posisi_jabatan, p.gaji_range
        FROM alumni a
        JOIN pekerjaan_alumni p ON a.id = p.alumni_id
        WHERE a.tahun_lulus = $1 
          AND p.gaji_range >= 4000000 
          AND p.status_pekerjaan = 'aktif'
    `, tahun)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
    }
    defer rows.Close()

    type Result struct {
        ID             int     `json:"id"`
        Nama           string  `json:"nama"`
        Jurusan        string  `json:"jurusan"`
        TahunLulus     int     `json:"tahun_lulus"`
        BidangIndustri string  `json:"bidang_industri"`
        NamaPerusahaan string  `json:"nama_perusahaan"`
        PosisiJabatan  string  `json:"posisi_jabatan"`
        RangeGaji      float64 `json:"range_gaji"`
    }

    var data []Result
    for rows.Next() {
        var r2 Result
        if err := rows.Scan(&r2.ID, &r2.Nama, &r2.Jurusan, &r2.TahunLulus,
            &r2.BidangIndustri, &r2.NamaPerusahaan, &r2.PosisiJabatan, &r2.RangeGaji); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca data"})
        }
        data = append(data, r2)
    }

    count := len(data)
    return c.JSON(fiber.Map{"success": true, "count": count, "data": data})
})

}