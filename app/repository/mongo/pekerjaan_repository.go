package mongo

import (
	"context"
	"crud-alumni/app/models"
	"crud-alumni/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PekerjaanRepository struct {
	collection *mongo.Collection
}

func NewPekerjaanRepository() *PekerjaanRepository {
	return &PekerjaanRepository{
		collection: database.MongoDB.Collection("pekerjaan"),
	}
}

func (r *PekerjaanRepository) Create(pekerjaan *models.Pekerjaan) error {
	pekerjaan.CreatedAt = time.Now()
	pekerjaan.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), bson.M{
		"alumni_id":             pekerjaan.AlumniID,
		"nama_perusahaan":       pekerjaan.NamaPerusahaan,
		"posisi_jabatan":        pekerjaan.PosisiJabatan,
		"bidang_industri":       pekerjaan.BidangIndustri,
		"lokasi_kerja":          pekerjaan.LokasiKerja,
		"gaji_range":            pekerjaan.GajiRange,
		"tanggal_mulai_kerja":   pekerjaan.TanggalMulaiKerja,
		"tanggal_selesai_kerja": pekerjaan.TanggalSelesaiKerja,
		"status_pekerjaan":      pekerjaan.StatusPekerjaan,
		"deskripsi":             pekerjaan.Deskripsi,
		"created_at":            pekerjaan.CreatedAt,
		"updated_at":            pekerjaan.UpdatedAt,
		"is_delete":             pekerjaan.IsDelete,
		"deleted_at":            pekerjaan.DeletedAt,
		"deleted_by":            pekerjaan.DeletedBy,
	})
	return err
}

func (r *PekerjaanRepository) FindAll(limit, offset int64) ([]models.Pekerjaan, int64, error) {
	filter := bson.M{"is_delete": false}

	// Count total documents
	total, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Set options for pagination
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var pekerjaan []models.Pekerjaan
	if err := cursor.All(context.Background(), &pekerjaan); err != nil {
		return nil, 0, err
	}
	return pekerjaan, total, nil
}

func (r *PekerjaanRepository) FindByID(id string) (*models.Pekerjaan, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pekerjaan models.Pekerjaan
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *PekerjaanRepository) FindByAlumniID(alumniID string) ([]models.Pekerjaan, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"alumni_id": alumniID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var pekerjaan []models.Pekerjaan
	if err := cursor.All(context.Background(), &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func (r *PekerjaanRepository) Update(id string, pekerjaan *models.Pekerjaan) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	pekerjaan.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"alumni_id":             pekerjaan.AlumniID,
			"nama_perusahaan":       pekerjaan.NamaPerusahaan,
			"posisi_jabatan":        pekerjaan.PosisiJabatan,
			"bidang_industri":       pekerjaan.BidangIndustri,
			"lokasi_kerja":          pekerjaan.LokasiKerja,
			"gaji_range":            pekerjaan.GajiRange,
			"tanggal_mulai_kerja":   pekerjaan.TanggalMulaiKerja,
			"tanggal_selesai_kerja": pekerjaan.TanggalSelesaiKerja,
			"status_pekerjaan":      pekerjaan.StatusPekerjaan,
			"deskripsi":             pekerjaan.Deskripsi,
			"updated_at":            pekerjaan.UpdatedAt,
			"is_delete":             pekerjaan.IsDelete,
			"deleted_at":            pekerjaan.DeletedAt,
			"deleted_by":            pekerjaan.DeletedBy,
		},
	}

	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	return err
}

// SoftDelete performs a soft delete on a pekerjaan
func (r *PekerjaanRepository) SoftDelete(id string, deletedBy int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_delete":  true,
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		},
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID, "is_delete": false},
		update,
	)
	return err
}

// HardDelete performs a permanent delete
func (r *PekerjaanRepository) HardDelete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

// GetTrash retrieves all soft-deleted items
func (r *PekerjaanRepository) GetTrash() ([]models.Pekerjaan, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"is_delete": true},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var pekerjaan []models.Pekerjaan
	if err := cursor.All(context.Background(), &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

// Restore restores a soft-deleted item
func (r *PekerjaanRepository) Restore(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_delete":  false,
			"deleted_at": nil,
			"deleted_by": nil,
		},
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID, "is_delete": true},
		update,
	)
	return err
}

// Search searches for pekerjaan based on keywords
func (r *PekerjaanRepository) Search(keyword string, limit, offset int64) ([]models.Pekerjaan, int64, error) {
	filter := bson.M{
		"is_delete": false,
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": keyword, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": keyword, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	// Count total matching documents
	total, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Set options for pagination and sorting
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var pekerjaan []models.Pekerjaan
	if err := cursor.All(context.Background(), &pekerjaan); err != nil {
		return nil, 0, err
	}

	return pekerjaan, total, nil
}
