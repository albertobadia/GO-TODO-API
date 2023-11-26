package api

import (
	"database/sql"
	"fmt"
)

func GetPostgresConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)
	return sql.Open("postgres", connectionString)
}
