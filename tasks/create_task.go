package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"task-manager-backend/db"
	"task-manager-backend/types"
)

type TaskResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

func CreateTask(res http.ResponseWriter, req *http.Request) {
	log.Println("Received create task request")
	reqTask := types.Task{}
	err := json.NewDecoder(req.Body).Decode(&reqTask)
	if err != nil {
		log.Printf("Failed to decode request body: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if reqTask.Title == "" {
		log.Println("Title is required")
		http.Error(res, "Title is required", http.StatusBadRequest)
		return
	}

	id := req.Header.Get("User")
	userId, errn := strconv.Atoi(id)
	if errn != nil {
		log.Println("invalid user id")
		http.Error(res, "Invalid user id", http.StatusBadRequest)
		return
	}

	reqTask.User = userId
	var tasksAction db.DatabaseInterface = db.DB

	if exists, _ := tasksAction.DoesTableExist("tasks"); !exists {
		err = tasksAction.CreateTaskTable()
		if err != nil {
			log.Printf("Failed to create tasks table: %v\n", err)
			http.Error(res, fmt.Sprintf("Failed to create tasks table: %v", err), http.StatusInternalServerError)
			return
		}
	}

	taskId, err := tasksAction.CreateTask(reqTask)

	resTask := TaskResponse{
		Id:      int(taskId),
		Message: "Task created",
	}
	if err != nil {
		log.Printf("Failed to create task: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}
	jsonRes, err := json.Marshal(resTask)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("Task created")

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusCreated)
	res.Write(jsonRes)

}
