package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	for i := 1; i <= 20; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil {
			err = DB.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("DB not ready, retrying in 5 seconds... (attempt %d/20)", i)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to DB after 20 attempts:", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,
		`CREATE TABLE IF NOT EXISTS boards(
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,
		`CREATE TABLE IF NOT EXISTS columns (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			position INT NOT NULL DEFAULT 0,
			board_id INT REFERENCES boards(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			column_id INT REFERENCES columns(id) ON DELETE RESTRICT,
			position INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal("Failed to create tables: ", err)
		}
	}

	log.Println("All tables are ready!")
}
