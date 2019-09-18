package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

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
func DBInit() {
	log.Print("Loading database..")
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	database = db
}

// DBGenerateTablesPrompt asks the user to generate tables if they don't already exist
func DBGenerateTablesPrompt() {
	yonp := func(predicate string) bool {
		fmt.Print(predicate + " [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		if input == 'y' {
			fmt.Println("OK!")
			return true
		}
		return false
	}
	if DBHasUserTable() == false {
		if yonp("User table does not exist, create one?") {
			DBCreateUserTable()
		}
		if yonp("Add four test users?") {
			DBGenerateTrash()
		}
	}
	if DBHasVideoTable() == false {
		if yonp("Video table does not exist, create one?") {
			DBCreateVideoTable()
		}
	}
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
	_, err := database.Exec(`CREATE TABLE videos (uuid uuid unique, link text,
		title text, description text, length text, likes numeric, dislikes numeric);
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}
	return nil
}

// DBVideoCreate adds a video to the database
func DBVideoCreate(video VideoData) error {
	command := `INSERT INTO videos (uuid, link, title, description, length, likes, dislikes)
		VALUES (uuid_generate_v1(), $1, $2, $3, $4, $5, $6);`
	_, err := database.Exec(command, video.Link, video.Title,
		video.Description, video.Length, video.Likes, video.Dislikes)
	if err != nil {
		return err
	}
	return nil
}

// DBVideoDelete deletes a video from the database
func DBVideoDelete(UUID string) error {
	command := `DELETE FROM videos WHERE uuid = $1`
	_, err := database.Exec(command, UUID)
	if err != nil {
		return err
	}
	return nil
}

// DBVideoInfo returns a Video's info
func DBVideoInfo(UUID string) (VideoData, error) {
	var done VideoData
	rows, err := database.Query("SELECT * FROM videos WHERE uuid = $1", UUID)
	if err != nil {
		return done, err
	}
	rows.Next()
	err = rows.Scan(&done.UUID, done.Link, done.Title, done.Description, done.Length, done.Likes, done.Dislikes)
	if err != nil {
		return done, err
	}
	return done, nil
}

// DBGenerateTrash generates `random` users to test api calls on
func DBGenerateTrash() {
	command := `INSERT INTO users (username, bio, email) VALUES ('test1', 'test bio 1', 'test@1.1');
		INSERT INTO users (username, bio, email) VALUES ('test2', 'test bio 2', 'test@2.2');
		INSERT INTO users (username, bio, email) VALUES ('test3', 'test bio 3', 'test@3.3');
		INSERT INTO users (username, bio, email) VALUES ('test4', 'test bio 4', 'test@4.4');`
	database.Exec(command)
}
