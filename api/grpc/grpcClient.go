package grpc

import (
	"context"
	"fmt"
	"time"

	pb "TaskService/generated/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.UserServiceClient

func ConnectGrpc() error {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}

	defer conn.Close()
	client = pb.NewUserServiceClient(conn)

	return nil
}

func GetUserByID(userID int64) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := client.GetUser(ctx, &pb.GetUserRequest{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("could not check user: %v", err)
	}

	fmt.Printf("user received: ID=%d, Name=%s\n", user.GetUserId(), user.GetName())

	return user, nil
}

func CheckUserID(userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := client.CheckUser(ctx, &pb.CheckUserRequest{UserID: userID})
	if err != nil {
		return false, fmt.Errorf("could not greet: %v", err)
	}

	return user.IsExists, nil
}
