package todos

import (
	validate "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Todo struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Title  string    `json:"title" validate:"required"`
	IsDone bool      `json:"is_done" default:"false"`
}

func NewTodo(title string, userID uuid.UUID) (Todo, error) {
	todo := Todo{
		ID:     uuid.New(),
		UserID: userID,
		Title:  title,
		IsDone: false,
	}
	validate := validate.New()
	errors := validate.Struct(todo)
	if errors != nil {
		return Todo{}, errors
	}
	return todo, nil
}

type TodoQuery struct {
	ID     uuid.UUID
	UserID uuid.UUID
	IsDone bool
}
