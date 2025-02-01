package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"TaskService/internal/env"
	"TaskService/internal/repository"
	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/createTask", createTask).Methods(http.MethodPost)
	r.HandleFunc("/getTask", getTaskByID).Methods(http.MethodGet)
	r.HandleFunc("/getTasks", getAllTasks).Methods(http.MethodGet)
	server := &http.Server{
		Addr:         env.GetEnvString("PORT"),
		Handler:      r,
		ReadTimeout:  env.GetTimeDuration("READ_TIMEOUT"),
		WriteTimeout: env.GetTimeDuration("WRITE_TIMEOUT"),
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

	taskData := repository.GetTaskByID(intID)

	encoderErr := json.NewEncoder(writer).Encode(taskData)
	if err != nil {
		log.Fatal(encoderErr)

		return
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	taskData := repository.Data{}
	err := json.NewDecoder(r.Body).Decode(&taskData)

	if err != nil {
		log.Print(err)
		return
	}

	resultData, errTask := repository.CreateTask(r.Context(), taskData)

	if errTask != nil {
		log.Print(errTask)
		return
	}

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

func getAllTasks(writer http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(writer).Encode(repository.GetAllTasks(r.Context()))
	if err != nil {
		log.Fatal(err)
	}
}
