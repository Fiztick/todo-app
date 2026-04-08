package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-app/backend/db"
	"todo-app/backend/middleware"
	"todo-app/backend/models"

	"github.com/gorilla/mux"
)

func GetColumns(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["boardID"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT user_id FROM boards WHERE id = $1`,
		boardID,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Board not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this board")
		return
	}

	rows, err := db.DB.Query(
		`SELECT id, title, position, board_id, created_at FROM columns
		WHERE board_id = $1 ORDER BY position ASC`,
		boardID,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch columns")
		return
	}
	defer rows.Close()

	var columns []models.Column
	for rows.Next() {
		var col models.Column
		err := rows.Scan(&col.ID, &col.Title, &col.Position, &col.BoardID, &col.CreatedAt)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to scan column")
			return
		}
		columns = append(columns, col)
	}

	if columns == nil {
		columns = []models.Column{}
	}

	respondJson(w, http.StatusOK, columns)
}

func CreateColumn(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["boardID"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT user_id FROM boards WHERE id = $1`,
		boardID,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Board not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this board")
		return
	}

	var body struct {
		Title string `json:"title"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	var col models.Column
	err = db.DB.QueryRow(
		`INSERT INTO columns (title, position, board_id)
		VALUES ($1, (SELECT COALESCE(MAX(position) + 1, 0) FROM columns WHERE board_id = $2), $2)
		RETURNING id, title, position, board_id, created_at`,
		body.Title, boardID,
	).Scan(&col.ID, &col.Title, &col.Position, &col.BoardID, &col.CreatedAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create column")
		return
	}

	respondJson(w, http.StatusCreated, col)
}

func UpdateColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid column ID")
		return
	}

	var col models.Column
	err = json.NewDecoder(r.Body).Decode(&col)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = db.DB.QueryRow(
		"UPDATE columns SET title = $1 WHERE id = $2 RETURNING id, title, position, created_at",
		col.Title, id,
	).Scan(&col.ID, &col.Title, &col.Position, &col.CreatedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Column not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update column")
		return
	}

	respondJson(w, http.StatusOK, col)
}

func DeleteColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid column ID")
		return
	}

	_, err = db.DB.Exec("DELETE FROM columns WHERE id = $1", id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete column")
		return
	}

	respondJson(w, http.StatusOK, map[string]string{"message": "Column deleted"})
}
