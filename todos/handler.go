package todos

import (
	"encoding/json"
	"net/http"
	"todo-api/api"
	"todo-api/users"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// TodoHandler handles HTTP requests for Todo
type TodoHandler struct {
	todoRepo TodoRepository
}

// UserHandler handles HTTP requests for User

// NewTodoHandler creates a new instance of TodoHandler
func NewTodoHandler(todoRepo TodoRepository) *TodoHandler {
	return &TodoHandler{
		todoRepo: todoRepo,
	}
}

// Implement HTTP handlers for Todo and User

func (h *TodoHandler) HandleListing(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	}
}

func (h *TodoHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	}
}

func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	user := users.GetUserFromContext(r.Context())
	todos, err := h.todoRepo.Query(TodoQuery{UserID: user.ID})
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if todos == nil {
		todos = []Todo{}
	}
	api.RespondWithJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := users.GetUserFromContext(r.Context())
	query := TodoQuery{
		ID:     uuid.MustParse(id),
		UserID: user.ID,
	}
	todo, err := h.todoRepo.Get(query)
	if err != nil {
		api.RespondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}
	api.RespondWithJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	context := r.Context()
	user := context.Value(users.UserCtxKey).(users.User)

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	todo, err = NewTodo(
		todo.Title,
		user.ID,
	)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	new_todo, err := h.todoRepo.Create(todo)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	api.RespondWithJSON(w, http.StatusCreated, new_todo)
}

func updateTodoData(previous *Todo, newData Todo) Todo {
	if newData.Title != "" {
		previous.Title = newData.Title
	}
	if newData.IsDone != previous.IsDone {
		previous.IsDone = newData.IsDone
	}
	return *previous
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := users.GetUserFromContext(r.Context())
	query := TodoQuery{ID: uuid.MustParse(id), UserID: user.ID}

	existing_todo, err := h.todoRepo.Get(query)
	if err != nil {
		api.RespondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}

	var update_data Todo
	err = json.NewDecoder(r.Body).Decode(&update_data)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	new_todo := updateTodoData(&existing_todo, update_data)

	err = h.todoRepo.Update(uuid.MustParse(id), new_todo)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Todo updated successfully"})
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := users.GetUserFromContext(r.Context())
	query := TodoQuery{ID: uuid.MustParse(id), UserID: user.ID}
	_, err := h.todoRepo.Get(query)
	if err != nil {
		api.RespondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}

	err = h.todoRepo.Delete(uuid.MustParse(id))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	api.RespondWithJSON(w, http.StatusNoContent, map[string]string{"message": "Todo deleted successfully"})
}
