package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pekerjaan struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           float64            `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   *time.Time         `bson:"tanggal_mulai_kerja,omitempty" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	Deskripsi           string             `bson:"deskripsi" json:"deskripsi"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
	IsDelete            bool               `bson:"is_delete" json:"is_delete"`
	DeletedAt           *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at"`
	DeletedBy           primitive.ObjectID `bson:"deleted_by,omitempty" json:"deleted_by"`
}

type PekerjaanDetail struct {
	Pekerjaan  `bson:",inline"`
	AlumniNama string `bson:"alumni_nama,omitempty" json:"alumni_nama,omitempty"`
	TrashInfo  *struct {
		DeletedAt   time.Time          `bson:"deleted_at" json:"deleted_at"`
		DeletedBy   primitive.ObjectID `bson:"deleted_by" json:"deleted_by"`
		RestoreNote string             `bson:"restore_note,omitempty" json:"restore_note,omitempty"`
	} `bson:"trash_info,omitempty" json:"trash_info,omitempty"`
}

type PekerjaanResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    *Pekerjaan `json:"data,omitempty"`
}

type PekerjaanListResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    []Pekerjaan `json:"data,omitempty"`
	Meta    struct {
		Total  int64 `json:"total"`
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	} `json:"meta,omitempty"`
}
