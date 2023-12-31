package todos

import (
	"errors"

	"github.com/google/uuid"
)

type MemoryTodoRepository struct {
	todos map[uuid.UUID]Todo
}

func NewMemoryTodoRepository() *MemoryTodoRepository {
	return &MemoryTodoRepository{
		todos: make(map[uuid.UUID]Todo),
	}
}

func (m *MemoryTodoRepository) Query(query TodoQuery) ([]Todo, error) {
	var todos []Todo
	for _, todo := range m.todos {
		if query.ID != uuid.Nil && todo.ID != query.ID {
			continue
		}
		if query.UserID != uuid.Nil && todo.UserID != query.UserID {
			continue
		}
		if query.IsDone && !todo.IsDone {
			continue
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (m *MemoryTodoRepository) Get(query TodoQuery) (Todo, error) {
	result, err := m.Query(query)
	if err != nil {
		return Todo{}, err
	}
	if len(result) == 0 {
		return Todo{}, errors.New("Todo not found")
	}
	return result[0], nil
}

func (m *MemoryTodoRepository) GetByID(id uuid.UUID) (Todo, error) {
	if todo, ok := m.todos[id]; ok {
		return todo, nil
	}
	return Todo{}, errors.New("Todo not found")
}

func (m *MemoryTodoRepository) Create(todo Todo) (Todo, error) {
	m.todos[todo.ID] = todo
	return todo, nil
}

func (m *MemoryTodoRepository) Update(id uuid.UUID, todo Todo) (Todo, error) {
	_, err := m.GetByID(id)
	if err != nil {
		return Todo{}, err
	}
	m.todos[id] = todo
	return todo, nil
}

func (m *MemoryTodoRepository) Delete(id uuid.UUID) error {
	_, err := m.GetByID(id)
	if err != nil {
		return err
	}
	delete(m.todos, id)
	return nil
}
