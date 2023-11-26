package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"todo-api/api"
	"todo-api/todos"
	"todo-api/users"

	"github.com/gorilla/mux"
)

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func GetRepositories() (users.UsersRepository, todos.TodoRepository, error) {
	if api.IS_TESTING {
		return users.NewMemoryUserRepository(), todos.NewMemoryTodoRepository(), nil
	}

	postgresConn, err := api.GetPostgresConnection()
	if err != nil {
		return nil, nil, err
	}
	api.MigrateUp(postgresConn)
	postgresUserRepo := users.NewUsersPostgresRepository(postgresConn)
	postgresTodoRepo := todos.NewPostgresTodoRepository(postgresConn)

	return postgresUserRepo, postgresTodoRepo, nil
}

func main() {
	router := mux.NewRouter()
	router.Use(logHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	userRepo, todoRepo, err := GetRepositories()
	usersHandler := users.NewUsersHandler(userRepo)
	todosHandler := todos.NewTodoHandler(todoRepo)
	if err != nil {
		log.Fatal(err)
		server.Shutdown(context.Background())
	}

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})
	})

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
	log.Fatal(server.ListenAndServe())
}
