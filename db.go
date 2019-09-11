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

// VideoData contains metadata about the video and a link to the video file
type VideoData struct {
	UUID        string
	Link        string
	Title       string
	Description string
	Length      string
	Likes       uint
	Dislikes    uint
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
	_, err := database.Exec(`CREATE TABLE videos (uuid uuid unique, link datalink,
		title text, description text, length numeric, likes numeric, dislikes numeric);
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}
	return nil
}

// DBGenerateTrash generates `random` users to test api calls on
func DBGenerateTrash() {
	command := `INSERT INTO users (username, bio, email) VALUES ('test1', 'test bio 1', 'test@1.1');
		INSERT INTO users (username, bio, email) VALUES ('test2', 'test bio 2', 'test@2.2');
		INSERT INTO users (username, bio, email) VALUES ('test3', 'test bio 3', 'test@3.3');
		INSERT INTO users (username, bio, email) VALUES ('test4', 'test bio 4', 'test@4.4');`
	database.Exec(command)
}
