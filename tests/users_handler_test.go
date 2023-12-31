package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-api/users"
)

func TestRegister(t *testing.T) {
	usersHandler := users.NewUsersHandler(users.NewMemoryUserRepository())

	jsonData := `{"username":"testuser","password":"testpassword"}`
	req, err := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestRegisterWithEmptyUsername(t *testing.T) {
	usersHandler := users.NewUsersHandler(users.NewMemoryUserRepository())

	jsonData := `{"username":"","password":"testpassword"}`
	req, err := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestRegisterWithEmptyPassword(t *testing.T) {
	usersHandler := users.NewUsersHandler(users.NewMemoryUserRepository())

	jsonData := `{"username":"testuser","password":""}`
	req, err := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLogin(t *testing.T) {
	repo := users.NewMemoryUserRepository()
	usersHandler := users.NewUsersHandler(repo)
	user, _ := users.NewUser("testuser", "testpassword")
	user.Password, _ = users.HashPassword(user.Password)
	repo.Create(user)
	jsonData := `{"username":"testuser","password":"testpassword"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoginWithEmptyUsername(t *testing.T) {
	usersHandler := users.NewUsersHandler(users.NewMemoryUserRepository())

	jsonData := `{"username":"","password":"testpassword"}`
	req, err := http.NewRequest("POST", "/login", strings.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLoginWithEmptyPassword(t *testing.T) {
	usersHandler := users.NewUsersHandler(users.NewMemoryUserRepository())

	jsonData := `{"username":"testuser","password":""}`
	req, err := http.NewRequest("POST", "/login", strings.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
