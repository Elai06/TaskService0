package main

import (
	"TaskService/api/server"
	"TaskService/internal/factory"
)

func main() {
	factory.ConnectToMongo()
	server.StartServer()
	server.ConnectGRpc()
}
