package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "TaskService/generated/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func GetUserByID(userID int64) *pb.GetUserResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetUser(ctx, &pb.GetUserRequest{UserID: userID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Printf("user received: ID=%d, Name=%s\n", res.GetUserId(), res.GetName())

	return res
}

func CheckUserID(userID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CheckUser(ctx, &pb.CheckUserRequest{UserID: userID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return res.IsExists
}
