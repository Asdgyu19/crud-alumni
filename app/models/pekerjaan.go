package models

import "time"


type Pekerjaan struct {
    ID                int        `json:"id"`
    AlumniID          int        `json:"alumni_id"`
    NamaPerusahaan    string     `json:"nama_perusahaan"`
    PosisiJabatan     string     `json:"posisi_jabatan"`
    BidangIndustri    string     `json:"bidang_industri"`
    LokasiKerja       string     `json:"lokasi_kerja"`
    GajiRange         float64    `json:"gaji_range"`
    TanggalMulaiKerja time.Time  `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja time.Time `json:"tanggal_selesai_kerja"`
    StatusPekerjaan   string     `json:"status_pekerjaan"`
    Deskripsi         string     `json:"deskripsi"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
    IsDelete          bool       `json:"is_delete"`
    DeletedAt        *time.Time `json:"deleted_at"`
    DeletedBy         int        `json:"deleted_by"`
}
