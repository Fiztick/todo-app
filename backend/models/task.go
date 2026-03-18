package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	ColumnID  int       `json:"column_id"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type Column struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}
