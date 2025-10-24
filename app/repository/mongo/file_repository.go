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

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository() *FileRepository {
	return &FileRepository{
		collection: database.MongoDB.Collection("files"),
	}
}

func (r *FileRepository) Create(file *models.File) error {
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), file)
	if err != nil {
		return err
	}

	file.ID = result.InsertedID
	return nil
}

func (r *FileRepository) FindByID(id string) (*models.File, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file models.File
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) FindAll() ([]models.File, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var files []models.File
	if err := cursor.All(context.Background(), &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) FindByAlumniID(alumniID string) ([]models.File, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"alumni_id": alumniID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var files []models.File
	if err := cursor.All(context.Background(), &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
