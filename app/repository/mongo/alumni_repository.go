package mongo

import (
	"context"
	"crud-alumni/app/models"
	"crud-alumni/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlumniRepository struct {
	collection *mongo.Collection
}

func NewAlumniRepository() *AlumniRepository {
	return &AlumniRepository{
		collection: database.MongoDB.Collection("alumni"),
	}
}

func (r *AlumniRepository) Create(alumni *models.Alumni) error {
	_, err := r.collection.InsertOne(context.Background(), bson.M{
		"nim":         alumni.NIM,
		"nama":        alumni.Nama,
		"jurusan":     alumni.Jurusan,
		"angkatan":    alumni.Angkatan,
		"tahun_lulus": alumni.TahunLulus,
		"email":       alumni.Email,
		"no_telepon":  alumni.NoTelepon,
		"alamat":      alumni.Alamat,
		"created_at":  alumni.CreatedAt,
		"updated_at":  alumni.UpdatedAt,
		"user_id":     alumni.UserID,
	})
	return err
}

func (r *AlumniRepository) FindAll() ([]models.Alumni, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var alumni []models.Alumni
	if err := cursor.All(context.Background(), &alumni); err != nil {
		return nil, err
	}
	return alumni, nil
}

func (r *AlumniRepository) FindByID(id string) (*models.Alumni, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var alumni models.Alumni
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&alumni)
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (r *AlumniRepository) Update(id string, alumni *models.Alumni) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"nim":         alumni.NIM,
			"nama":        alumni.Nama,
			"jurusan":     alumni.Jurusan,
			"angkatan":    alumni.Angkatan,
			"tahun_lulus": alumni.TahunLulus,
			"email":       alumni.Email,
			"no_telepon":  alumni.NoTelepon,
			"alamat":      alumni.Alamat,
			"updated_at":  time.Now(),
			"user_id":     alumni.UserID,
		},
	}

	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	return err
}

func (r *AlumniRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
