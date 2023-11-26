package tests

import (
	"todo-api/users"

	"github.com/google/uuid"
)

func GetUserAndToken(repo users.UsersRepository) (users.User, string) {
	randomUsername := uuid.New().String()
	user, _ := users.NewUser(randomUsername, "test")
	repo.Create(user)
	token, _ := users.GenerateToken(user.Username)
	return user, token
}
