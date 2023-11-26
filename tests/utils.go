package tests

import "todo-api/users"

func GetUserAndToken(repo users.UsersRepository) (users.User, string) {
	user, _ := users.NewUser("test", "test")
	repo.Create(user)
	token, _ := users.GenerateToken(user.Username)
	return user, token
}
