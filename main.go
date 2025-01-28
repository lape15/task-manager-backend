package main

import (
	"fmt"
	"task-manager-backend/db"
	"task-manager-backend/middleware"
	"task-manager-backend/tasks"
	"task-manager-backend/user"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func handleMainRequest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte("Welcome to the API"))

	default:
		fmt.Print("Default route\n")

		http.Error(res, "Method not allowed there", http.StatusMethodNotAllowed)
	}
}

func handleTasksRequest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		tasks.CreateTask(res, req)
	default:
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	db.ConnectDb()
	route := mux.NewRouter()

	route.HandleFunc("/signup", user.SignUp).Methods("POST")
	route.HandleFunc("/login", user.Login).Methods("POST")
	route.HandleFunc("/tasks", middleware.AuthMiddleware(handleTasksRequest))
	route.HandleFunc("/", handleMainRequest)
	http.ListenAndServe(":8080", route)
	defer db.DB.Close()
}
