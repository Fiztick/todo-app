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

func GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["boardId"])
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
		`SELECT tasks.id, tasks.title, tasks.completed, tasks.column_id, tasks.position, tasks.created_at 
		FROM tasks 
		JOIN columns ON tasks.column_id = columns.id
		WHERE columns.board_id = $1
		ORDER BY position ASC`,
		boardID,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var columnID sql.NullInt64
		err := rows.Scan(&task.ID, &task.Title, &task.Completed, &columnID, &task.Position, &task.CreatedAt)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to scan task")
			return
		}
		if columnID.Valid {
			task.ColumnID = int(columnID.Int64)
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	respondJson(w, http.StatusOK, tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var body struct {
		Title    string `json:"title"`
		ColumnID int    `json:"column_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT boards.user_id FROM columns
		JOIN boards ON columns.board_id = boards.id
		WHERE columns.id = $1`,
		body.ColumnID,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Column not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this column")
		return
	}

	var task models.Task
	err = db.DB.QueryRow(
		`INSERT INTO tasks (title, column_id, position)
		VALUES ($1, $2, (SELECT COALESCE(MAX(position) + 1, 0) FROM tasks WHERE column_id = $2))
		RETURNING id, title, completed, COALESCE(column_id, 0), position, created_at`,
		body.Title, body.ColumnID,
	).Scan(&task.ID, &task.Title, &task.Completed, &task.ColumnID, &task.Position, &task.CreatedAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	respondJson(w, http.StatusCreated, task)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT boards.user_id FROM tasks
		JOIN columns ON tasks.column_id = columns.id
		JOIN boards ON columns.board_id = boards.id
		WHERE tasks.id = $1`,
		id,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this task")
		return
	}

	var task models.Task
	var columnID sql.NullInt64
	err = db.DB.QueryRow(
		`UPDATE tasks SET completed = NOT completed WHERE id = $1
		RETURNING id, title, completed, column_id, position, created_at`,
		id,
	).Scan(&task.ID, &task.Title, &task.Completed, &columnID, &task.Position, &task.CreatedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}
	if columnID.Valid {
		task.ColumnID = int(columnID.Int64)
	}

	respondJson(w, http.StatusOK, task)
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT boards.user_id FROM tasks
		JOIN columns ON tasks.column_id = columns.id
		JOIN boards ON columns.board_id = boards.id
		WHERE tasks.id = $1`,
		id,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this task")
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

	var task models.Task
	var columnID sql.NullInt64
	err = db.DB.QueryRow(
		`UPDATE tasks SET title = $1 WHERE id = $2
		RETURNING id, title, completed, column_id, position, created_at`,
		body.Title, id,
	).Scan(&task.ID, &task.Title, &task.Completed, &columnID, &task.Position, &task.CreatedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to edit task")
		return
	}
	if columnID.Valid {
		task.ColumnID = int(columnID.Int64)
	}

	respondJson(w, http.StatusOK, task)
}

func MoveTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT boards.user_id FROM tasks
		JOIN columns ON tasks.column_id = columns.id
		JOIN boards ON columns.board_id = boards.id
		WHERE tasks.id = $1`,
		id,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this task")
		return
	}

	var body struct {
		ColumnID int `json:"column_id"`
		Position int `json:"position"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var task models.Task
	err = db.DB.QueryRow(
		`UPDATE tasks SET column_id = $1, position = $2 WHERE id = $3
		RETURNING id, title, completed, COALESCE(column_id, 0), position, created_at`,
		body.ColumnID, body.Position, id,
	).Scan(&task.ID, &task.Title, &task.Completed, &task.ColumnID, &task.Position, &task.CreatedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to move task")
		return
	}

	respondJson(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var ownerID int
	err = db.DB.QueryRow(
		`SELECT boards.user_id FROM tasks
		JOIN columns ON tasks.column_id = columns.id
		JOIN boards ON columns.board_id = boards.id
		WHERE tasks.id = $1`,
		id,
	).Scan(&ownerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if ownerID != userID {
		respondError(w, http.StatusForbidden, "You don't own this task")
		return
	}

	_, err = db.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	respondJson(w, http.StatusOK, map[string]string{"message": "Task deleted"})
}
