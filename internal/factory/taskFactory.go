package factory

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var collection *mongo.Collection

type Data struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Id          int64  `json:"id"`
	UserId      int64  `json:"userId"`
}

func ConnectToMongo() {
	url := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!", client)
	collection = client.Database("TasksService").Collection("tasks")
}

func CreateTask(taskData Data) *mongo.InsertOneResult {
	taskData.Id = getNextUserID()
	insertResult, err := collection.InsertOne(context.TODO(), taskData)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

	return insertResult
}

func getNextUserID() int64 {
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})

	var lastTask Data
	err := collection.FindOne(context.TODO(), bson.D{}, opts).Decode(&lastTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1
		}
		log.Fatal(err)
	}

	return lastTask.Id + 1
}

func GetTaskById(id int64) Data {
	filter := map[string]interface{}{"id": id}

	result := Data{}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return Data{}
	}

	fmt.Printf("Found document: %+v\n", result)

	return result
}

func GetAllTasks() []Data {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result []Data

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	return result
}
