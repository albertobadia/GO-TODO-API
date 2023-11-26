package users

import (
	"github.com/google/uuid"

	validate "github.com/go-playground/validator/v10"
)

type User struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Username string    `json:"username" validate:"required"`
	Password string    `json:"password" validate:"required"`
}

func NewUser(username, password string) (User, error) {
	user := User{
		ID:       uuid.New(),
		Username: username,
		Password: password,
	}
	errors := validate.New().Struct(user)
	if errors != nil {
		return User{}, errors
	}
	return user, nil
}

type UserRead struct {
	ID       uuid.UUID
	Username string
}

type LoginResponse struct {
	Token string
}
