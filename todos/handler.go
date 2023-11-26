package todos

import (
	"encoding/json"
	"net/http"
	"todo-api/api"
	"todo-api/users"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TodoHandler struct {
	todoRepo TodoRepository
}

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
		api.RespondWithError(w, http.StatusInternalServerError, "Error getting all todos")
		return
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
	if err != nil || todo.ID == uuid.Nil {
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
		api.RespondWithError(w, http.StatusInternalServerError, "Error creating todo on repository")
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

	var update_data Todo
	err := json.NewDecoder(r.Body).Decode(&update_data)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	existing_todo, err := h.todoRepo.Get(query)
	if err != nil || existing_todo.ID == uuid.Nil {
		api.RespondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}
	new_todo := updateTodoData(&existing_todo, update_data)

	new_todo, err = h.todoRepo.Update(uuid.MustParse(id), new_todo)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Error updating todo on repository")
		return
	}
	api.RespondWithJSON(w, http.StatusOK, new_todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	user := users.GetUserFromContext(r.Context())
	query := TodoQuery{ID: id, UserID: user.ID}
	result, err := h.todoRepo.Get(query)
	if err != nil || result.ID == uuid.Nil {
		api.RespondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}

	err = h.todoRepo.Delete(id)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Error deleting todo")
		return
	}
	api.RespondWithJSON(w, http.StatusNoContent, nil)
}
