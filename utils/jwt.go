package utils

import (
	"errors"
	"os"
	"time"

	"crud-alumni/app/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("mysecretkey") // ganti dengan env

func init() {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		jwtSecret = []byte(s)
	}
}

func GenerateToken(u models.User) (string, error) {
	claims := models.JWTClaims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}
func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
