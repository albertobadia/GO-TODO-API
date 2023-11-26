package todos

import (
	"github.com/google/uuid"
)

type Todo struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Title  string    `json:"title"`
	IsDone bool      `json:"is_done"`
}

type TodoQuery struct {
	ID     uuid.UUID
	UserID uuid.UUID
	IsDone bool
}
