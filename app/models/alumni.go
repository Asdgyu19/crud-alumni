package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID            interface{} `json:"id" bson:"_id,omitempty"`
	NIM           string      `json:"nim" bson:"nim"`
	Nama          string      `json:"nama" bson:"nama"`
	Jurusan       string      `json:"jurusan" bson:"jurusan"`
	Angkatan      int         `json:"angkatan" bson:"angkatan"`
	TahunLulus    int         `json:"tahun_lulus" bson:"tahun_lulus"`
	Email         string      `json:"email" bson:"email"`
	NoTelepon     string      `json:"no_telepon" bson:"no_telepon"`
	Alamat        string      `json:"alamat" bson:"alamat"`
	CreatedAt     time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" bson:"updated_at"`
	UserID        int         `json:"user_id" bson:"user_id"`
	PekerjaanID   int         `json:"pekerjaan_id" bson:"pekerjaan_id"`
	NamaPekerjaan string      `json:"nama_pekerjaan" bson:"nama_pekerjaan"`
}

// GetObjectID returns the MongoDB ObjectID
func (a *Alumni) GetObjectID() primitive.ObjectID {
	if oid, ok := a.ID.(primitive.ObjectID); ok {
		return oid
	}
	return primitive.NilObjectID
}
