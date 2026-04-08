package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-app/backend/db"
	"todo-app/backend/middleware"
	"todo-app/backend/models"

	"github.com/gorilla/mux"
)

func GetBoards(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	rows, err := db.DB.Query(
		`SELECT id, title, user_id, created_at FROM boards
		WHERE user_id = $1 ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch boards")
		return
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var board models.Board
		err := rows.Scan(&board.ID, &board.Title, &board.UserID, &board.CreatedAt)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to scan board")
			return
		}
		boards = append(boards, board)
	}
	respondJson(w, http.StatusOK, boards)
}

func CreateBoard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var body struct {
		Title string `json:"title"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	var board models.Board
	err = db.DB.QueryRow(
		`INSERT INTO boards (title, user_id)
		VALUES ($1, $2)
		RETURNING id, title, user_id, created_at`,
		body.Title, userID,
	).Scan(&board.ID, &board.Title, &board.UserID, &board.CreatedAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create board")
		return
	}

	defaultColumns := []string{"To Do", "In Progress", "Done"}
	for i, title := range defaultColumns {
		_, err = db.DB.Exec(
			`INSERT INTO columns (title, position, board_id)
			VALUES ($1, $2, $3)`,
			title, i, board.ID,
		)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create default columns")
			return
		}
	}

	respondJson(w, http.StatusCreated, board)
}

func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT user_id FROM boards WHERE id = $1`,
		id,
	).Scan(&ownerID)

	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this board")
		return
	}

	_, err = db.DB.Exec(`DELETE FROM boards WHERE id = $1`, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete board")
		return
	}

	respondJson(w, http.StatusOK, map[string]string{"message": "Board deleted"})
}

func UpdateBoard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid board ID")
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

	var board models.Board
	err = db.DB.QueryRow(
		`UPDATE boards SET title = $1
		WHERE id = $2 AND user_id = $3
		RETURNING id, title, user_id, created_at`,
		body.Title, id, userID,
	).Scan(&board.ID, &board.Title, &board.UserID, &board.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "Board not found or unauthorized")
		return
	}

	respondJson(w, http.StatusOK, board)
}
