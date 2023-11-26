package api

import "database/sql"

func MigrateUp(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS todos (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			title TEXT NOT NULL,
			is_done BOOLEAN NOT NULL
		);
	`)
	return err
}

func MigrateDown(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS todos;
		DROP TABLE IF EXISTS users;
	`)
	return err
}
