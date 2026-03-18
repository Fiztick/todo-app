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
	for i := 1; i <= 10; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil {
			err = DB.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("DB not ready, retrying in 3 seconds... (attempt %d/10)", i)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to DB after 10 attempts:", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS columns (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			position INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			column_id INT REFERENCES columns(id) ON DELETE SET NULL,
			position INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
		`INSERT INTO columns (title, position)
		SELECT 'To Do', 0
		WHERE NOT EXISTS (SELECT 1 FROM columns)`,
		`INSERT INTO columns (title, position)
		SELECT 'In Progress', 1
		WHERE NOT EXISTS (SELECT 1 FROM columns OFFSET 1)`,
		`INSERT INTO columns (title, position)
		SELECT 'Done', 2
		WHERE NOT EXISTS (SELECT 1 FROM columns OFFSET 2)`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal("Failed to create tasks table: ", err)
		}
	}

	log.Println("Tasks table is ready!")
}
