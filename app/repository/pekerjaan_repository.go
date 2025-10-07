package repository

import (
	"crud-alumni/app/models"
	"crud-alumni/database"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func GetAllPekerjaanRepo() ([]models.Pekerjaan, error) {
	rows, err := database.DB.Query(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Pekerjaan
	for rows.Next() {
		var p models.Pekerjaan
		var tMulai, tSelesai sql.NullTime
		rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai, &p.StatusPekerjaan, &p.Deskripsi, &p.CreatedAt, &p.UpdatedAt)
		if tMulai.Valid {
			p.TanggalMulaiKerja = tMulai.Time
		}
		if tSelesai.Valid {
			p.TanggalSelesaiKerja = tSelesai.Time
		}
		list = append(list, p)
	}
	return list, nil
}

func GetPekerjaanByIDRepo(id int) (models.Pekerjaan, error) {
	var p models.Pekerjaan
	var tMulai, tSelesai sql.NullTime
	rows, err := database.DB.Query(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni WHERE id=$1`, id)
	if err != nil {
		return models.Pekerjaan{}, err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai, &p.StatusPekerjaan, &p.Deskripsi, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return models.Pekerjaan{}, err
		}
		if tMulai.Valid {
			p.TanggalMulaiKerja = tMulai.Time
		}
		if tSelesai.Valid {
			p.TanggalSelesaiKerja = tSelesai.Time
		}
		return p, nil
	}
	return models.Pekerjaan{}, sql.ErrNoRows
}

func GetPekerjaanByAlumniRepo(alumniID int) ([]models.Pekerjaan, error) {
	rows, err := database.DB.Query(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni WHERE alumni_id=$1`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Pekerjaan
	for rows.Next() {
		var p models.Pekerjaan
		var tMulai, tSelesai sql.NullTime
		rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai, &p.StatusPekerjaan, &p.Deskripsi, &p.CreatedAt, &p.UpdatedAt)
		if tMulai.Valid {
			p.TanggalMulaiKerja = tMulai.Time
		}
		if tSelesai.Valid {
			p.TanggalSelesaiKerja = tSelesai.Time
		}
		list = append(list, p)
	}
	return list, nil
}

func CreatePekerjaanRepo(p *models.Pekerjaan) error {
	rows, err := database.DB.Query(`
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		RETURNING id, created_at, updated_at
	`, p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.Deskripsi)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	}
	return sql.ErrNoRows
}

func UpdatePekerjaanRepo(id int, p *models.Pekerjaan) error {
	_, err := database.DB.Exec(`UPDATE pekerjaan_alumni SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3, lokasi_kerja=$4, gaji_range=$5, tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7, status_pekerjaan=$8, deskripsi_pekerjaan=$9, updated_at=NOW() WHERE id=$10`,
		p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.Deskripsi, id)
	return err
}

func DeletePekerjaanRepo(id int) error {
	_, err := database.DB.Exec("DELETE FROM pekerjaan_alumni WHERE id=$1", id)
	return err
}


// GetPekerjaanWithFilter 
func GetPekerjaanWithFilter(search, sortBy, order string, limit, offset int) ([]models.Pekerjaan, error) {
	query := fmt.Sprintf(`
    SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
           lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
           status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
    FROM pekerjaan_alumni
    WHERE is_delete = false
      AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1)
    ORDER BY %s %s
    LIMIT $2 OFFSET $3
`, sortBy, order)


	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []models.Pekerjaan
	for rows.Next() {
		var p models.Pekerjaan
		err := rows.Scan(
			&p.ID,
			&p.AlumniID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.BidangIndustri,
			&p.LokasiKerja,
			&p.GajiRange,
			&p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja,
			&p.StatusPekerjaan,
			&p.Deskripsi,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}

	return pekerjaanList, nil
}

// CountPekerjaan -> hitung total pekerjaan untuk pagination
func CountPekerjaan(search string) (int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM pekerjaan_alumni
		WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1
	`
	err := database.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func SoftDeletePekerjaan(id int, userID int, role string) error {
	if role == "user" {
		var count int
		err := database.DB.QueryRow(`
			SELECT COUNT(*) 
			FROM pekerjaan_alumni pa
			JOIN alumni a ON pa.alumni_id = a.id
			WHERE pa.id = $1 AND a.user_id = $2 AND pa.is_delete = false
		`, id, userID).Scan(&count)

		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("tidak diizinkan hapus pekerjaan orang lain")
		}
	}

	return SoftDeletePekerjaanRepo(id, userID)
}



func SoftDeletePekerjaanRepo(id int, userID int) error {
	_, err := database.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET is_delete = true, deleted_at = NOW(), deleted_by = $1
		WHERE id = $2
	`, userID, id)
	return err
}