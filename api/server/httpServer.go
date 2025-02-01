package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"TaskService/internal/factory"
	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/createTask", createTask).Methods("POST")
	r.HandleFunc("/getTask", getTaskByID).Methods("GET")
	r.HandleFunc("/getTasks", getAllTasks).Methods("GET")

	server := &http.Server{
		Addr:         ":8082",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port 8082")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err, "Error starting server on port 8080")
	}
}

func getTaskByID(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		log.Fatal(writer, "id is not liquid")
	}

	intID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Fatal(writer, "id is not liquid")
	}

	taskData := factory.GetTaskByID(intID)

	encoderErr := json.NewEncoder(writer).Encode(taskData)
	if err != nil {
		log.Fatal(encoderErr)

		return
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	taskData := factory.Data{}
	err := json.NewDecoder(r.Body).Decode(&taskData)

	if err != nil {
		log.Fatal(err)
		return
	}

	resultData := factory.CreateTask(taskData)

	log.Println("TaskCreated")

	errEncoder := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task created successfully",
		"result":  resultData.InsertedID,
	})
	if err != nil {
		log.Fatal(errEncoder)
		return
	}
}

func getAllTasks(writer http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(writer).Encode(factory.GetAllTasks())
	if err != nil {
		log.Fatal(err)
	}
}
