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

// Implement UserRepository interface
func (m *MemoryUserRepository) GetByUsername(username string) (User, error) {
	//return user if exist error otherwise
	if user, ok := m.users[username]; ok {
		return user, nil
	}
	return User{}, errors.New("User not found")
}

func (m *MemoryUserRepository) Create(user User) error {
	_, err := m.GetByUsername(user.Username)
	if err == nil {
		return errors.New("User already exist")
	}
	m.users[user.Username] = user
	return nil
}
