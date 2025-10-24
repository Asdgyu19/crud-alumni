package models

import (
	"encoding/json"
	"time"
)

type Pekerjaan struct {
	ID                  int        `json:"id" gorm:"primaryKey"`
	AlumniID            int        `json:"alumni_id" gorm:"not null"`
	NamaPerusahaan      string     `json:"nama_perusahaan" gorm:"type:varchar(255);not null"`
	PosisiJabatan       string     `json:"posisi_jabatan" gorm:"type:varchar(255);not null"`
	BidangIndustri      string     `json:"bidang_industri" gorm:"type:varchar(255)"`
	LokasiKerja         string     `json:"lokasi_kerja" gorm:"type:varchar(255)"`
	GajiRange           float64    `json:"gaji_range"`
	TanggalMulaiKerja   *time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `json:"status_pekerjaan" gorm:"type:varchar(50)"`
	Deskripsi           string     `json:"deskripsi" gorm:"type:text"`

	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null"`
	IsDelete  bool       `json:"is_delete" gorm:"default:false"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy int        `json:"deleted_by"`
}


type PekerjaanTrash struct {
	ID            int             `json:"id" gorm:"primaryKey"`
	PekerjaanID   int             `json:"pekerjaan_id" gorm:"not null"`
	PekerjaanData json.RawMessage `json:"pekerjaan_data" gorm:"type:jsonb;not null"` 
	DeletedBy     int             `json:"deleted_by" gorm:"not null"`
	DeletedAt     time.Time       `json:"deleted_at" gorm:"not null"`
	RestoredAt    *time.Time      `json:"restored_at"`
	RestoredBy    *int            `json:"restored_by"`
	RestoreNote   string          `json:"restore_note" gorm:"type:text"`

	
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}


type PekerjaanDetail struct {
	*Pekerjaan
	AlumniNama string `json:"alumni_nama,omitempty"`
	TrashInfo  *struct {
		DeletedAt   time.Time `json:"deleted_at"`
		DeletedBy   string    `json:"deleted_by"`
		RestoreNote string    `json:"restore_note,omitempty"`
	} `json:"trash_info,omitempty"`
}

type PekerjaanListResponse struct {
	Data []PekerjaanDetail `json:"data"`
	Meta MetaInfo          `json:"meta"`
}

type PekerjaanTrashListResponse struct {
	Data []PekerjaanTrash `json:"data"`
	Meta MetaInfo         `json:"meta"`
}

type SinglePekerjaanResponse struct {
	Data PekerjaanDetail `json:"data"`
}
