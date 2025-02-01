package server

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"

	"TaskService/internal/env"
	"TaskService/internal/repository"
	"github.com/gorilla/mux"
)

type taskResult struct {
	Message string
	Result  *mongo.InsertOneResult
}

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
		log.Print(err, "Error starting server on port 8080")
		return
	}
}

func getTaskByID(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		log.Print(writer, "id is not liquid")
		return
	}

	intID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Print(writer, "id is not liquid")
		return
	}

	taskData := repository.GetTaskByID(intID)

	encoderErr := json.NewEncoder(writer).Encode(taskData)
	if encoderErr != nil {
		log.Println(encoderErr)
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

	var resultMessage taskResult
	resultMessage.Result = resultData
	resultMessage.Message = "Task Created"

	errEncoder := json.NewEncoder(w).Encode(&resultMessage)
	if errEncoder != nil {
		log.Print(errEncoder)
		return
	}
}

func getAllTasks(writer http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(writer).Encode(repository.GetAllTasks(r.Context()))
	if err != nil {
		log.Print(err)
		return
	}
}
