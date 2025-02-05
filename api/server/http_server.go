package server

import (
	"encoding/json"
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
	repo repository.TaskServiceInterface
}

func NewTaskHandler(repo repository.TaskServiceInterface) *TaskHandler {
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
		log.Print(err, "Error starting server on port", env.GetEnvString("PORT"))
		return err
	}

	return err
}

func (th *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	taskData := repository.Data{}

	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		log.Print(err)
		return
	}

	resultData, errTask := th.repo.CreateTask(r.Context(), taskData)

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

func (th *TaskHandler) getAllTasks(writer http.ResponseWriter, r *http.Request) {
	result, repoErr := th.repo.GetAllTasks(r.Context())
	if repoErr != nil {
		log.Print(repoErr)
		return
	}
	err := json.NewEncoder(writer).Encode(result)
	if err != nil {
		log.Print(err)
		return
	}
}

func (th *TaskHandler) getTaskByID(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		log.Print(writer, "Failed to get", id)
		return
	}

	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Print(writer, "Failed to parse int", err)
		return
	}

	taskData, err := th.repo.GetTaskByID(intID)
	if err != nil {
		log.Print(err)
		return
	}
	encoderErr := json.NewEncoder(writer).Encode(taskData)
	if encoderErr != nil {
		log.Print(encoderErr)
		return
	}
}
