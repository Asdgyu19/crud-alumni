package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        interface{} `json:"id" bson:"_id,omitempty"`
	Username  string      `json:"username" bson:"username"`
	Password  string      `json:"-" bson:"password"`
	Email     string      `json:"email" bson:"email"`
	Role      string      `json:"role" bson:"role"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type JWTClaims struct {
	UserID   interface{} `json:"user_id"`
	Username string      `json:"username"`
	Role     string      `json:"role"`
	jwt.RegisteredClaims
}
