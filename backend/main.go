package main

import (
	"log"
	"net/http"

	"todo-app/backend/db"
	"todo-app/backend/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlers.CompleteTask).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
