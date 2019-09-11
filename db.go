package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var database *sql.DB

// UserData contains every field in a user's account
type UserData struct {
	Name  string
	Bio   string
	Email string
}

// DBInit logs in to the database
func DBInit() error {
	fmt.Println("Loading database..")
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	database = db
	return nil
}

// DBHasUserTable returns true if the database is loaded and contains a user table
func DBHasUserTable() bool {
	_, err := database.Query("SELECT * FROM users")
	if err != nil {
		return false
	}
	return true
}

// DBCreateUserTable creates a user table in the database
func DBCreateUserTable() error {
	_, err := database.Exec("CREATE TABLE users (username text unique, bio text, email text);")
	if err != nil {
		return err
	}
	return nil
}

// DBHasVideoTable returns true if the database is loaded and contains a user table
func DBHasVideoTable() bool {
	_, err := database.Query("SELECT * FROM videos")
	if err != nil {
		return false
	}
	return true
}

// DBCreateVideoTable creates a video table in the database
func DBCreateVideoTable() error {
	_, err := database.Exec("CREATE TABLE videos (id uuid unique, title text);")
	if err != nil {
		return err
	}
	return nil
}

// DBUserCreate adds a user into the database
func DBUserCreate(user UserData) error {
	command := `INSERT INTO users (username, bio, email) VALUES ($1, $2, $3);`
	_, err := database.Exec(command, user.Name, user.Bio, user.Email)
	if err != nil {
		return err
	}
	return nil
}

// DBUserDelete removes a user from the database
func DBUserDelete(username string) error {
	command := `DELETE FROM users WHERE username = $1;`
	_, err := database.Exec(command, username)
	if err != nil {
		return err
	}
	return nil
}

// DBUserInfo returns a user's profile info
func DBUserInfo(username string) (UserData, error) {
	var done UserData
	rows, err := database.Query("SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return done, err
	}
	rows.Next()
	err = rows.Scan(&done.Name, &done.Bio, &done.Email)
	if err != nil {
		return done, err
	}
	return done, nil
}
