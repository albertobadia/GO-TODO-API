package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"todo-api/api"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	userRepo UsersRepository
}

func NewUsersHandler(userRepo UsersRepository) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
	}
}

func GenerateToken(username string) (string, error) {
	claims := jwt.StandardClaims{
		Subject: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error hashing password")
	}
	return string(hash), nil
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Username == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Username is required")
		return
	}

	if user.Password == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	user.ID = uuid.New()
	user.Password, _ = HashPassword(user.Password)

	err = h.userRepo.Create(user)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user_read := UserRead{
		ID:       user.ID,
		Username: user.Username,
	}

	api.RespondWithJSON(w, http.StatusCreated, user_read)
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	{
		var loginData User
		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if loginData.Username == "" {
			api.RespondWithError(w, http.StatusBadRequest, "Username is required")
			return
		}

		if loginData.Password == "" {
			api.RespondWithError(w, http.StatusBadRequest, "Password is required")
			return
		}

		user, err := h.userRepo.GetByUsername(loginData.Username)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		if !CheckPassword(loginData.Password, user.Password) {
			api.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		token, _ := GenerateToken(user.Username)
		api.RespondWithJSON(w, http.StatusOK, LoginResponse{Token: token})
	}
}

func (h *UsersHandler) GetUserFromToken(token string) (User, error) {
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return User{}, err
	}
	return h.userRepo.GetByUsername(claims.Subject)
}

func (h *UsersHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		token = token[7:]
		user, err := h.GetUserFromToken(token)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		new_context := context.WithValue(r.Context(), UserCtxKey, user)
		next(w, r.WithContext(new_context))
	}
}
