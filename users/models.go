package users

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type UserRead struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type JWTClaims struct {
	Username       string `json:"username"`
	StandardClaims jwt.StandardClaims
}

type LoginResponse struct {
	Token string `json:"token"`
}
