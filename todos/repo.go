package todos

import "github.com/google/uuid"

type TodoRepository interface {
	Query(query TodoQuery) ([]Todo, error)
	Get(query TodoQuery) (Todo, error)
	Create(todo Todo) (Todo, error)
	Update(id uuid.UUID, todo Todo) (Todo, error)
	Delete(id uuid.UUID) error
}
