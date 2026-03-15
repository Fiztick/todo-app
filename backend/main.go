package main

import (
	"log"
	"net/http"

	"todo-app/backend/db"
	"todo-app/backend/handlers"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks/{id}", handlers.CompleteTask).Methods("PATCH", "OPTIONS")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
