package users

import "errors"

type MemoryUserRepository struct {
	users map[string]User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]User),
	}
}

func (m *MemoryUserRepository) GetByUsername(username string) (User, error) {
	if user, ok := m.users[username]; ok {
		return user, nil
	}
	return User{}, errors.New("User not found")
}

func (m *MemoryUserRepository) Create(user User) (User, error) {
	_, err := m.GetByUsername(user.Username)
	if err == nil {
		return User{}, errors.New("User already exist")
	}
	m.users[user.Username] = user
	return user, nil
}
