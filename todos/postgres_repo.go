package todos

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresTodoRepository struct {
	db *sql.DB
}

func NewPostgresTodoRepository(db *sql.DB) *PostgresTodoRepository {
	return &PostgresTodoRepository{
		db: db,
	}
}

func (p *PostgresTodoRepository) Query(query TodoQuery) ([]Todo, error) {
	var todos []Todo
	rows, err := p.db.Query("SELECT id, user_id, title, is_done FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.IsDone)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (p *PostgresTodoRepository) Get(query TodoQuery) (Todo, error) {
	result, err := p.Query(query)
	if err != nil {
		return Todo{}, err
	}
	if len(result) == 0 {
		return Todo{}, nil
	}
	return result[0], nil
}

func (p *PostgresTodoRepository) Create(todo Todo) (Todo, error) {
	_, err := p.db.Exec("INSERT INTO todos (id, user_id, title, is_done) VALUES ($1, $2, $3, $4)", todo.ID, todo.UserID, todo.Title, todo.IsDone)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (p *PostgresTodoRepository) Update(id uuid.UUID, todo Todo) (Todo, error) {
	_, err := p.db.Exec("UPDATE todos SET title = $1, is_done = $2 WHERE id = $3", todo.Title, todo.IsDone, id)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (p *PostgresTodoRepository) Delete(id uuid.UUID) error {
	_, err := p.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
