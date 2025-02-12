package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"TaskService/internal/env"
	"TaskService/internal/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskResult struct {
	Message string
	Result  *mongo.InsertOneResult
}

type TaskHandler struct {
	repo repository.ITaskService
}

func NewTaskHandler(repo repository.ITaskService) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (th *TaskHandler) StartServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/createTask", th.createTask).Methods(http.MethodPost)
	r.HandleFunc("/getTask", th.getTaskByID).Methods(http.MethodGet)
	r.HandleFunc("/getTasks", th.getAllTasks).Methods(http.MethodGet)
	server := &http.Server{
		Addr:         env.GetEnvString("PORT"),
		Handler:      r,
		ReadTimeout:  env.GetTimeDuration("READ_TIMEOUT"),
		WriteTimeout: env.GetTimeDuration("WRITE_TIMEOUT"),
	}
	log.Print("Starting server on port")
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server on port: %s", env.GetEnvString("PORT"))
	}
	return nil
}

func (th *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	taskData := repository.Data{}
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	resultData, errTask := th.repo.CreateTask(r.Context(), taskData)
	if errTask != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errTask)
		return
	}
	var resultMessage taskResult
	resultMessage.Result = resultData
	resultMessage.Message = "Task Created"
	errEncoder := json.NewEncoder(w).Encode(&resultMessage)
	if errEncoder != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errEncoder)
		return
	}
}

func (th *TaskHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, repoErr := th.repo.GetAllTasks(ctx)
	if repoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(repoErr)
		return
	}
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (th *TaskHandler) getTaskByID(w http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(w, "Failed to get", id)
		return
	}
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(w, "Failed to parse int", err)
		return
	}
	taskData, err := th.repo.GetTaskByID(request.Context(), intID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	encoderErr := json.NewEncoder(w).Encode(taskData)
	if encoderErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(encoderErr)
		return
	}
}
