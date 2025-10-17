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
			p.TanggalMulaiKerja = &tMulai.Time
		}
		if tSelesai.Valid {
			p.TanggalSelesaiKerja = &tSelesai.Time
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
			p.TanggalMulaiKerja = &tMulai.Time
		}
		if tSelesai.Valid {
			p.TanggalSelesaiKerja = &tSelesai.Time
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
			p.TanggalMulaiKerja = &tMulai.Time
		}
		if tSelesai.Valid {
			p.TanggalSelesaiKerja = &tSelesai.Time
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
	fmt.Printf("DEBUG: SoftDeletePekerjaanRepo() - ID: %d, UserID: %d\n", id, userID)

	result, err := database.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET is_delete = true, deleted_at = NOW(), deleted_by = $1
		WHERE id = $2
	`, userID, id)

	if err != nil {
		fmt.Printf("DEBUG: Error soft delete: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("DEBUG: Soft delete berhasil. Rows affected: %d\n", rowsAffected)

	return nil
}

// GetTrashRepo - ambil semua pekerjaan yang sudah di-soft delete
func GetTrashRepo() ([]models.Pekerjaan, error) {
	fmt.Println("DEBUG: GetTrashRepo() dipanggil")

	rows, err := database.DB.Query(`
		SELECT pa.id, pa.alumni_id, pa.bidang_industri, pa.nama_perusahaan, 
		       pa.posisi_jabatan, pa.gaji_range, pa.status_pekerjaan, 
		       pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, pa.deleted_at, pa.created_at
		FROM pekerjaan_alumni pa
		WHERE pa.is_delete = true AND pa.deleted_at IS NOT NULL
		ORDER BY pa.deleted_at DESC
	`)
	if err != nil {
		fmt.Printf("DEBUG: Error query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var trashList []models.Pekerjaan
	rowCount := 0
	for rows.Next() {
		rowCount++
		var p models.Pekerjaan
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.BidangIndustri, &p.NamaPerusahaan,
			&p.PosisiJabatan, &p.GajiRange, &p.StatusPekerjaan,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.DeletedAt, &p.CreatedAt,
		)
		if err != nil {
			fmt.Printf("DEBUG: Error scan row %d: %v\n", rowCount, err)
			return nil, err
		}
		fmt.Printf("DEBUG: Berhasil scan row %d - ID: %d, Perusahaan: %s\n", rowCount, p.ID, p.NamaPerusahaan)
		trashList = append(trashList, p)
	}

	fmt.Printf("DEBUG: Total rows found: %d\n", len(trashList))
	return trashList, nil
}

// RestoreRepo - restore pekerjaan dari trash
func RestoreRepo(id int) error {
	result, err := database.DB.Exec(`
		UPDATE pekerjaan_alumni 
		SET is_delete = false, deleted_at = NULL, deleted_by = NULL
		WHERE id = $1 AND is_delete = true
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("pekerjaan tidak ditemukan di trash atau sudah di-restore")
	}

	return nil
}

// HardDeleteRepo - hapus permanent dari database
func HardDeleteRepo(id int) error {
	// Cek apakah sudah di-soft delete dulu
	var isDeleted bool
	err := database.DB.QueryRow(`
		SELECT is_delete FROM pekerjaan_alumni WHERE id = $1
	`, id).Scan(&isDeleted)

	if err != nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	if !isDeleted {
		return errors.New("pekerjaan harus di-soft delete terlebih dahulu sebelum hard delete")
	}

	// Hapus permanent dari database
	result, err := database.DB.Exec(`
		DELETE FROM pekerjaan_alumni WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("gagal menghapus pekerjaan")
	}

	return nil
}

// GetMyTrashRepo - user lihat trash milik sendiri
func GetMyTrashRepo(userID int) ([]models.Pekerjaan, error) {
	fmt.Printf("DEBUG: GetMyTrashRepo() dipanggil untuk userID: %d\n", userID)

	// Query debug untuk cek semua data soft deleted dulu
	var totalDeleted int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM pekerjaan_alumni WHERE is_delete = true
	`).Scan(&totalDeleted)
	if err == nil {
		fmt.Printf("DEBUG: Total semua data soft deleted: %d\n", totalDeleted)
	}

	// Query debug untuk cek data user ini
	var userAlumniID int
	err = database.DB.QueryRow(`
		SELECT id FROM alumni WHERE user_id = $1
	`, userID).Scan(&userAlumniID)
	if err != nil {
		fmt.Printf("DEBUG: Error cari alumni_id untuk user %d: %v\n", userID, err)
		return nil, err
	}
	fmt.Printf("DEBUG: User %d memiliki alumni_id: %d\n", userID, userAlumniID)

	// Cek data soft deleted untuk alumni ini
	var userDeleted int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM pekerjaan_alumni WHERE alumni_id = $1 AND is_delete = true
	`, userAlumniID).Scan(&userDeleted)
	if err == nil {
		fmt.Printf("DEBUG: Data soft deleted untuk alumni_id %d: %d\n", userAlumniID, userDeleted)
	}

	rows, err := database.DB.Query(`
		SELECT pa.id, pa.alumni_id, pa.bidang_industri, pa.nama_perusahaan, 
		       pa.posisi_jabatan, pa.gaji_range, pa.status_pekerjaan, 
		       pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, pa.deleted_at, pa.created_at
		FROM pekerjaan_alumni pa
		JOIN alumni a ON pa.alumni_id = a.id
		WHERE pa.is_delete = true AND pa.deleted_at IS NOT NULL AND a.user_id = $1
		ORDER BY pa.deleted_at DESC
	`, userID)
	if err != nil {
		fmt.Printf("DEBUG: Error query GetMyTrashRepo: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var trashList []models.Pekerjaan
	rowCount := 0
	for rows.Next() {
		rowCount++
		var p models.Pekerjaan
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.BidangIndustri, &p.NamaPerusahaan,
			&p.PosisiJabatan, &p.GajiRange, &p.StatusPekerjaan,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.DeletedAt, &p.CreatedAt,
		)
		if err != nil {
			fmt.Printf("DEBUG: Error scan row %d: %v\n", rowCount, err)
			return nil, err
		}
		fmt.Printf("DEBUG: User trash row %d - ID: %d, Perusahaan: %s\n", rowCount, p.ID, p.NamaPerusahaan)
		trashList = append(trashList, p)
	}

	fmt.Printf("DEBUG: Total user trash found: %d\n", len(trashList))
	return trashList, nil
}

// RestoreMyRepo - user restore pekerjaan milik sendiri
func RestoreMyRepo(id int, userID int) error {
	fmt.Printf("DEBUG: RestoreMyRepo() dipanggil - ID: %d, UserID: %d\n", id, userID)

	// Cek kepemilikan dan apakah di trash
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM pekerjaan_alumni pa
		JOIN alumni a ON pa.alumni_id = a.id
		WHERE pa.id = $1 AND a.user_id = $2 AND pa.is_delete = true
	`, id, userID).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("pekerjaan tidak ditemukan di trash atau bukan milik Anda")
	}

	// Restore data
	result, err := database.DB.Exec(`
		UPDATE pekerjaan_alumni 
		SET is_delete = false, deleted_at = NULL, deleted_by = NULL
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("gagal restore pekerjaan")
	}

	fmt.Printf("DEBUG: Berhasil restore pekerjaan ID: %d\n", id)
	return nil
}

// HardDeleteMyRepo - user hard delete pekerjaan milik sendiri dari trash
func HardDeleteMyRepo(id int, userID int) error {
	fmt.Printf("DEBUG: User %d mencoba hard delete pekerjaan ID: %d\n", userID, id)

	// Cek ownership dan apakah data dalam trash
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM pekerjaan_alumni 
		WHERE id = $1 AND alumni_id = $2 AND is_delete = true
	`, id, userID).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("pekerjaan tidak ditemukan dalam trash atau bukan milik Anda")
	}

	// Hard delete
	result, err := database.DB.Exec(`
		DELETE FROM pekerjaan_alumni 
		WHERE id = $1 AND alumni_id = $2 AND is_delete = true
	`, id, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("gagal menghapus pekerjaan secara permanen")
	}

	fmt.Printf("DEBUG: Berhasil hard delete pekerjaan ID: %d oleh user: %d\n", id, userID)
	return nil
}
