package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"todo-app/backend/db"
	"todo-app/backend/models"

	"github.com/gorilla/mux"
)

func respondJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, map[string]string{"error": message})
}

// ---------------
// Column Handlers
// ---------------

func GetColumns(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, position, created_at FROM columns ORDER BY position ASC")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch columns")
		return
	}
	defer rows.Close()

	var columns []models.Column
	for rows.Next() {
		var col models.Column
		err := rows.Scan(&col.ID, &col.Title, &col.Position, &col.CreatedAt)
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
	var col models.Column
	err := json.NewDecoder(r.Body).Decode(&col)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = db.DB.QueryRow(
		"INSERT INTO columns (title, position) VALUES ($1, (SELECT COALESCE(MAX(position) + 1, 0) FROM columns)) RETURNING id, title, position, created_at",
		col.Title,
	).Scan(&col.ID, &col.Title, &col.Position, &col.CreatedAt)
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

// -------------
// Task Handlers
// -------------

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, completed, column_id, position, created_at FROM tasks ORDER BY position ASC")
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
	var body struct {
		Title    string `json:"title"`
		ColumnID int    `json:"column_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
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

func MoveTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	_, err = db.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	respondJson(w, http.StatusOK, map[string]string{"message": "Task Deleted"})
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Task ID")
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
	).Scan(&task.ID, &task.Title, &task.Completed, &columnID, &task.Position, task.CreatedAt)
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
