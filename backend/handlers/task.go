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

var taskList = &models.TaskList{}

func respondJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, map[string]string{"error": message})
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, completed, created_at FROM tasks ORDER BY created_at ASC")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	defer rows.Close()

	taskList = &models.TaskList{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to scan task")
			return
		}
		taskList.Enqueue(task)
	}
	respondJson(w, http.StatusOK, taskList.GetAll())
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = db.DB.QueryRow(
		"INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, completed, created_at",
		task.Title,
	).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	taskList.Enqueue(task)
	respondJson(w, http.StatusCreated, task)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var task models.Task
	err = db.DB.QueryRow(
		"UPDATE tasks SET completed = TRUE WHERE id = $1 RETURNING id, title, completed, created_at",
		id,
	).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	taskList.Remove(id)
	taskList.Enqueue(task)
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

	taskList.Remove(id)
	respondJson(w, http.StatusOK, map[string]string{"message": "Task Deleted"})
}
