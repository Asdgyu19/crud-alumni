package repository

import (
	"crud-alumni/app/models"
	"crud-alumni/database"
	"database/sql"
	"errors"
)

func GetUserByIdentifier(identifier string) (models.User, string, error) {
	var u models.User
	var passwordHash string
	row := database.DB.QueryRow(`SELECT id, username, password_hash, role FROM users WHERE username=$1 OR email=$1`, identifier)
	err := row.Scan(&u.ID, &u.Username, &passwordHash, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, "", errors.New("user not found")
		}
		return models.User{}, "", err
	}
	return u, passwordHash, nil
}
