package users

type UsersRepository interface {
	GetByUsername(username string) (User, error)
	Create(user User) (User, error)
}
