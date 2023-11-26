package users

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type UsersPostgresRepository struct {
	db *sql.DB
}

func NewUsersPostgresRepository(db *sql.DB) *UsersPostgresRepository {
	return &UsersPostgresRepository{
		db: db,
	}
}

func (u *UsersPostgresRepository) GetByUsername(username string) (User, error) {
	var user User
	err := u.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UsersPostgresRepository) Create(user User) (User, error) {
	_, err := u.db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", user.ID, user.Username, user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
