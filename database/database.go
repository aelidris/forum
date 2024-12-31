// Package database handles the database connection and schema initialization
package database

import (
	"database/sql"
	"log"

	// SQLite3 driver imported for its side effects only
	_ "github.com/mattn/go-sqlite3"
)

// DB is a global variable that holds the database connection
var DB *sql.DB

// InitDB initializes the SQLite database connection and creates necessary tables
// if they don't already exist. It uses a local file "forum.db" as the database.
func InitDB() {
	var err error
	// Open a connection to the SQLite database
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// SQL statements to create the database schema

	// Users table stores user account information
	// - username and email must be unique
	// - email must contain @ and a domain (e.g., example@domain.com)
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Posts table stores forum posts
	// - linked to users table via user_id foreign key
	// - automatically timestamps post creation
	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	// Comments table stores comments on posts
	// - linked to both users and posts tables via foreign keys
	// - automatically timestamps comment creation
	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    UNIQUE (user_id, post_id)
	);`

	// Execute table creation statements
	// If any statement fails, the program will terminate
	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createPostsTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createCommentsTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createLikesTable)
	if err != nil {
		log.Fatal(err)
	}
}
