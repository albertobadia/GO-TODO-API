package main

import (
	"log"
	"net/http"
	"os"
	"todo-api/todos"
	"todo-api/users"

	"github.com/gorilla/mux"
)

func main() {
	memoryUserRepo := users.NewMemoryUserRepository()
	memoryTodoRepo := todos.NewMemoryTodoRepository()

	usersHandler := users.NewUsersHandler(memoryUserRepo)
	todosHandler := todos.NewTodoHandler(memoryTodoRepo)

	router := mux.NewRouter()

	router.HandleFunc("/register", usersHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", usersHandler.Login).Methods(http.MethodPost)

	router.HandleFunc("/todos", usersHandler.AuthMiddleware(
		todosHandler.HandleListing)).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/todos/{id}", usersHandler.AuthMiddleware(
		todosHandler.HandleItem)).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
