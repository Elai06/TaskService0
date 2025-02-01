package server

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "TaskService/generated/proto"
	"google.golang.org/grpc"
)

var client pb.UserServiceClient

func ConnectGRpc() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client = pb.NewUserServiceClient(conn)
}

func GetUserById(userId int64) *pb.GetUserResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetUser(ctx, &pb.GetUserRequest{UserId: userId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Print("user received: ID=%s, Name=%s\n", res.GetUserId(), res.GetName())

	return res
}
