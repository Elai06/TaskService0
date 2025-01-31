package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	pb "TaskService/generated/proto"
	"google.golang.org/grpc"
)

var client pb.UserServiceClient

func ConnectGRpc() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

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

	fmt.Printf("user received: ID=%d, Name=%s\n", res.GetUserId(), res.GetName())

	return res
}
