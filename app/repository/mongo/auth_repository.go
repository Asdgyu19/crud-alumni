package mongo

import (
	"context"
	"crud-alumni/app/models"
	"crud-alumni/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		collection: database.MongoDB.Collection("users"),
	}
}

func (r *AuthRepository) ValidateUser(username, password string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{
		"username": username,
		"password": password,
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{
		"username": username,
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
