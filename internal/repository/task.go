package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"TaskService/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tasks = "tasks"
const TasksService = "TasksService"

var collection *mongo.Collection

type Data struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
}

func ConnectToMongo() {
	url := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database(TasksService).Collection(tasks)
}

func CreateTask(ctx context.Context, taskData Data) (*mongo.InsertOneResult, error) {
	if isExist, err := checkUserID(taskData.UserID); err != nil && !isExist {
		return nil, fmt.Errorf("User Not Found", err)
	}

	taskData.ID = getNextTaskID()

	insertResult, err := collection.InsertOne(ctx, taskData)
	if err != nil {
		return nil, fmt.Errorf("Error inserting task", err)
	}

	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

	return insertResult, err
}

func getNextTaskID() int64 {
	var lastTask Data

	err := collection.FindOne(context.TODO(), bson.M{"id": -1}).Decode(&lastTask)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 1
		}

		log.Print(err)
	}

	return lastTask.ID + 1
}

func GetTaskByID(id int64) Data {
	result := Data{}

	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	if err != nil {
		log.Print(err)
		return Data{}
	}

	fmt.Printf("Found document: %+v\n", result)

	return result
}

func GetAllTasks(ctx context.Context) []Data {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var result []Data

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Print(err)
		return nil
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		log.Print(err)
		return nil
	}

	return result
}

func checkUserID(userID int64) (bool, error) {
	return grpc.CheckUserID(userID)
}
