package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-api/todos"
	"todo-api/users"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestHandleGetAllNoResult(t *testing.T) {
	todoHandler := todos.NewTodoHandler(todos.NewMemoryTodoRepository())

	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	_, token := GetUserAndToken(usersRepo)

	req, _ := http.NewRequest("GET", "/todos", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.AuthMiddleware(todoHandler.HandleListing))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleGetAllWithResult(t *testing.T) {
	repo := todos.NewMemoryTodoRepository()
	todoHandler := todos.NewTodoHandler(repo)
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	user, token := GetUserAndToken(usersRepo)
	todo, _ := todos.NewTodo("Test", user.ID)
	repo.Create(todo)

	req, _ := http.NewRequest("GET", "/todos", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.AuthMiddleware(todoHandler.HandleListing))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":"` + todo.ID.String() + `","user_id":"` + user.ID.String() + `","title":"Test","is_done":false}]`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestHandleGetItemNotFound(t *testing.T) {
	r := mux.NewRouter()
	todoHandler := todos.NewTodoHandler(todos.NewMemoryTodoRepository())
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("GET")

	id := uuid.New().String()
	_, token := GetUserAndToken(usersRepo)

	req, _ := http.NewRequest("GET", "/todos/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleGetItemFound(t *testing.T) {
	r := mux.NewRouter()
	repo := todos.NewMemoryTodoRepository()
	todoHandler := todos.NewTodoHandler(repo)
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("GET")

	user, token := GetUserAndToken(usersRepo)
	todo, _ := todos.NewTodo("Test", user.ID)
	todo, _ = repo.Create(todo)

	req, _ := http.NewRequest("GET", "/todos/"+todo.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":"` + todo.ID.String() + `","user_id":"` + user.ID.String() + `","title":"Test","is_done":false}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleCreate(t *testing.T) {
	todoHandler := todos.NewTodoHandler(todos.NewMemoryTodoRepository())
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	_, token := GetUserAndToken(usersRepo)

	jsonData := `{"title":"Test"}`
	req, _ := http.NewRequest("POST", "/todos", strings.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.AuthMiddleware(todoHandler.Create))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestHandleUpdateNotFound(t *testing.T) {
	todoHandler := todos.NewTodoHandler(todos.NewMemoryTodoRepository())
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r := mux.NewRouter()
	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("PUT")

	id := uuid.New().String()
	_, token := GetUserAndToken(usersRepo)
	jsonData := `{"title":"Test"}`

	req, _ := http.NewRequest("PUT", "/todos/"+id, strings.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleUpdateInvalidBody(t *testing.T) {
	todoHandler := todos.NewTodoHandler(todos.NewMemoryTodoRepository())
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r := mux.NewRouter()
	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("PUT")

	id := uuid.New().String()
	_, token := GetUserAndToken(usersRepo)
	jsonData := `{"title":}`

	req, _ := http.NewRequest("PUT", "/todos/"+id, strings.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandleUpdateNotOwned(t *testing.T) {
	r := mux.NewRouter()
	repo := todos.NewMemoryTodoRepository()
	todoHandler := todos.NewTodoHandler(repo)
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("PUT")

	user, _ := GetUserAndToken(usersRepo)
	todo, _ := todos.NewTodo("Test", user.ID)
	todo, _ = repo.Create(todo)

	_, otherToken := GetUserAndToken(usersRepo)
	jsonData := `{"title":"Updated title"}`
	req, _ := http.NewRequest("PUT", "/todos/"+todo.ID.String(), strings.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+otherToken)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleUpdateFound(t *testing.T) {
	r := mux.NewRouter()
	repo := todos.NewMemoryTodoRepository()
	todoHandler := todos.NewTodoHandler(repo)
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("PUT")

	user, token := GetUserAndToken(usersRepo)
	todo, _ := todos.NewTodo("Test", user.ID)
	todo, _ = repo.Create(todo)
	jsonData := `{"title":"Updated title"}`

	req, _ := http.NewRequest("PUT", "/todos/"+todo.ID.String(), strings.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	updated_todo, _ := repo.GetByID(todo.ID)
	if updated_todo.Title != "Updated title" {
		t.Errorf("handler returned unexpected body: got %v want %v", updated_todo.Title, "Updated title")
	}
}

func TestHandleDeleteNotFound(t *testing.T) {
	r := mux.NewRouter()
	todoHandler := todos.NewTodoHandler(todos.NewMemoryTodoRepository())
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("DELETE")

	id := uuid.New().String()
	_, token := GetUserAndToken(usersRepo)

	req, _ := http.NewRequest("DELETE", "/todos/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleDeleteNotOwned(t *testing.T) {
	r := mux.NewRouter()
	repo := todos.NewMemoryTodoRepository()
	todoHandler := todos.NewTodoHandler(repo)
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("DELETE")

	user, _ := GetUserAndToken(usersRepo)
	todo, _ := todos.NewTodo("Test", user.ID)
	todo, _ = repo.Create(todo)

	_, otherToken := GetUserAndToken(usersRepo)
	req, _ := http.NewRequest("DELETE", "/todos/"+todo.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+otherToken)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleDeleteFound(t *testing.T) {
	r := mux.NewRouter()
	repo := todos.NewMemoryTodoRepository()
	todoHandler := todos.NewTodoHandler(repo)
	usersRepo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(usersRepo)

	r.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(todoHandler.HandleItem)).Methods("DELETE")

	user, token := GetUserAndToken(usersRepo)
	todo, _ := todos.NewTodo("Test", user.ID)
	todo, _ = repo.Create(todo)
	req, _ := http.NewRequest("DELETE", "/todos/"+todo.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	_, err := repo.GetByID(todo.ID)
	if err == nil {
		t.Errorf("handler returned unexpected body: got %v want %v", err, "not nil")
	}
}
