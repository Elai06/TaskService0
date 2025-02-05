package main

import (
	"TaskService/api/grpc"
	"TaskService/api/server"
	"TaskService/internal/env"
	"TaskService/internal/repository"
)

func main() {
	env.LoadEnv()

	repo := repository.NewTaskRepository("mongodb://localhost:27017")
	srv := server.NewTaskHandler(*repo.Task)
	srvErr := srv.StartServer()
	if srvErr != nil {
		print(srvErr)
	}

	grpcErr := grpc.ConnectGrpc()
	if grpcErr != nil {
		print(grpcErr)
	}
}
