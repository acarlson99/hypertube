package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var database *sql.DB

// UserData contains every field in a user's account
type UserData struct {
	username string
	bio      string
	email    string
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

// DBUserInfo returns a user's profile info
func DBUserInfo(username string) (string, error) {
	var done string
	rows, err := database.Query("SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return "LOL!!!!!!1!!!", err
	}
	rows.Next()
	err = rows.Scan(&done)
	if err != nil {
		return "LOL!L!!!!!!1!!", err
	}
	return done, nil
}
