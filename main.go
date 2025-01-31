package main

import (
	"TaskService/api/grpc"
	"TaskService/api/server"
	"TaskService/internal/factory"
)

func main() {
	factory.ConnectToMongo()
	server.StartServer()
	grpc.ConnectGRpc()
}
