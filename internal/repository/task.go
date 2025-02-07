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

type Data struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
}

//go:generate mockgen -destination=mocks/mock_user_cache.go -package=mocks TaskService/internal/repository TaskServiceInterface

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, taskData Data) (*mongo.InsertOneResult, error)
	GetTaskByID(ctx context.Context, id int64) (Data, error)
	GetAllTasks(ctx context.Context) ([]Data, error)
}

type TaskRepository struct {
	client *mongo.Client
	Task   *TaskServiceInterface
}

func NewTaskRepository(connectionString string) *TaskRepository {
	clientOptions := options.Client().ApplyURI(connectionString)

	newClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)
	}

	return &TaskRepository{client: newClient, Task: new(TaskServiceInterface)}
}

func (t *TaskRepository) CreateTask(ctx context.Context, taskData Data) (*mongo.InsertOneResult, error) {
	if isExist, err := checkUserID(taskData.UserID); err != nil && !isExist {
		return nil, fmt.Errorf("User Not Found %s", err)
	}

	taskData.ID = t.getNextTaskID(ctx)

	insertResult, err := t.getCollection().InsertOne(ctx, taskData)
	if err != nil {
		return nil, fmt.Errorf("Error inserting task %s", err)
	}

	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

	var nums []int
	for _, i := range nums {
		if i%2 != 0 {
		}
	}

	return insertResult, nil
}

func (t *TaskRepository) getNextTaskID(ctx context.Context) int64 {
	var lastTask Data

	err := t.getCollection().FindOne(ctx, bson.M{"id": -1}).Decode(&lastTask)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 1
		}

		log.Print(err)
	}

	return lastTask.ID + 1
}

func (t *TaskRepository) GetTaskByID(ctx context.Context, id int64) (Data, error) {
	result := Data{}

	err := t.getCollection().FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		log.Print(err)
		return Data{}, err
	}

	fmt.Printf("Found document: %+v\n", result)

	return result, nil
}

func (t *TaskRepository) GetAllTasks(ctx context.Context) ([]Data, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var result []Data

	cursor, err := t.getCollection().Find(context.TODO(), bson.M{})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}

func checkUserID(userID int64) (bool, error) {
	return grpc.CheckUserID(userID)
}

func (t *TaskRepository) getCollection() *mongo.Collection {
	return t.client.Database(TasksService).Collection(tasks)
}
