package models

import "time"

type Column struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	BoardID   int       `json:"board_id"`
	CreatedAt time.Time `json:"created_at"`
}
