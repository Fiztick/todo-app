package models

import "time"

type Board struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
