package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/utils"
)

func Authenticate(identifier, password string) (models.User, string, error) {
	user, hash, err := repository.GetUserByIdentifier(identifier)
	if err != nil {
		return models.User{}, "", err
	}
	if !utils.CheckPassword(password, hash) {
		return models.User{}, "", err
	}
	token, err := utils.GenerateToken(user)
	return user, token, err
}
