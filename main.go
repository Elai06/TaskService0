package main

import (
	"TaskService/api/grpc"
	"TaskService/api/server"
	"TaskService/internal/env"
	"TaskService/internal/repository"
)

func main() {
	env.LoadEnv()
	repository.ConnectToMongo()
	server.StartServer()
	grpcErr := grpc.ConnectGrpc()

	if grpcErr != nil {
		print(grpcErr)
	}
}
