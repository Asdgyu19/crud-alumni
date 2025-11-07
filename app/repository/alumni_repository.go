package repository

import (
	"crud-alumni/app/models"
	"crud-alumni/database"
	"errors"
	"fmt"
)

func GetAllAlumniRepo() ([]models.Alumni, error) {
	rows, err := database.DB.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Alumni
	for rows.Next() {
		var a models.Alumni
		rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		list = append(list, a)
	}
	return list, nil
}

func GetAlumniByIDRepo(id int) (models.Alumni, error) {
	var a models.Alumni
	// Gunakan Query dan Next agar Scan dapat menerima semua kolom yang dipilih
	rows, err := database.DB.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id=$1`, id)
	if err != nil {
		return models.Alumni{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return models.Alumni{}, err
		}
		return a, nil
	}
	return models.Alumni{}, errors.New("not found")
}
func CreateAlumniRepo(a *models.Alumni) error {
	rows, err := database.DB.Query(`
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
		RETURNING id, created_at, updated_at
	`, a.NIM, a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.Email, a.NoTelepon, a.Alamat)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return err
		}
		return nil
	}
	return errors.New("failed to insert alumni")
}

func UpdateAlumniRepo(id int, a *models.Alumni) error {
	_, err := database.DB.Exec(`UPDATE alumni SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, email=$5, no_telepon=$6, alamat=$7, updated_at=NOW() WHERE id=$8`,
		a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.Email, a.NoTelepon, a.Alamat, id)
	return err
}

func DeleteAlumniRepo(id int) error {
	res, err := database.DB.Exec("DELETE FROM alumni WHERE id=$1", id)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetAlumniWithFilter(search, sortBy, order string, limit, offset int) ([]models.Alumni, error) {
	query := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE nama ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Alumni
	for rows.Next() {
		var a models.Alumni
		rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		list = append(list, a)
	}
	return list, nil
}

func CountAlumni(search string) (int, error) {
	var total int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1`, "%"+search+"%").Scan(&total)
	return total, err
}

func GetAlumniByTahunAndGajiRepo(tahun int) ([]models.Result, error) {
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
		return nil, err
	}
	defer rows.Close()

	var data []models.Result
	for rows.Next() {
		var r models.Result
		if err := rows.Scan(&r.ID, &r.Nama, &r.Jurusan, &r.TahunLulus,
			&r.BidangIndustri, &r.NamaPerusahaan, &r.PosisiJabatan, &r.RangeGaji); err != nil {
			return nil, err
		}
		data = append(data, r)
	}

	return data, nil
}
