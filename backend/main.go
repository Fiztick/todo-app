package main

import (
	"log"
	"net/http"

	"todo-app/backend/db"
	"todo-app/backend/handlers"
	"todo-app/backend/middleware"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db.Connect()

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Public routes
	r.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")

	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Protected routes
	protected.HandleFunc("/boards", handlers.GetBoards).Methods("GET", "OPTIONS")
	protected.HandleFunc("/boards", handlers.CreateBoard).Methods("POST", "OPTIONS")
	protected.HandleFunc("/boards/{id}", handlers.UpdateBoard).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/boards/{id}", handlers.DeleteBoard).Methods("DELETE", "OPTIONS")

	protected.HandleFunc("/boards/{boardId}/columns", handlers.GetColumns).Methods("GET", "OPTIONS")
	protected.HandleFunc("/boards/{boardId}/columns", handlers.CreateColumn).Methods("POST", "OPTIONS")
	protected.HandleFunc("/columns/{id}", handlers.UpdateColumn).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/columns/{id}", handlers.DeleteColumn).Methods("DELETE", "OPTIONS")

	protected.HandleFunc("/boards/{boardId}/tasks", handlers.GetTasks).Methods("GET", "OPTIONS")
	protected.HandleFunc("/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")
	protected.HandleFunc("/tasks/{id}", handlers.EditTask).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/tasks/{id}/complete", handlers.CompleteTask).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/tasks/{id}/move", handlers.MoveTask).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
