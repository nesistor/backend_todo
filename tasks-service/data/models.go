package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		Task: Task{},
	}
}

type Models struct {
	Task Task
}

type Task struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string    `bson:"user_id" json:"user_id"`
	Title       string    `bson:"name" json:"name"`
	Description string    `bson:"data" json:"data"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (t *Task) InsertTask(task Task) error {
	collection := client.Database("tasks").Collection("tasks")

	_, err := collection.InsertOne(context.TODO(), Task{
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		log.Println("Error inserting task:", err)
		return err
	}

	return nil
}

func (t *Task) UpdateTask(updatedTask Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("tasks").Collection("tasks")

	taskDocID, err := primitive.ObjectIDFromHex(updatedTask.ID)
	if err != nil {
		return err
	}

	userDocID, err := primitive.ObjectIDFromHex(updatedTask.UserID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     taskDocID,
		"user_id": userDocID,
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: updatedTask.Title},
			{Key: "description", Value: updatedTask.Description},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) GetAllByUserID(userID string) ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("tasks").Collection("tasks")

	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Task

	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Error decoding task:", err)

			continue
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *Task) DeleteTask(taskID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("tasks").Collection("tasks")

	taskDocID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		log.Println("Error converting taskID:", err)
		return err
	}

	userDocID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Error converting userID:", err)
		return err
	}

	// Define a filter to match both task ID and user ID
	filter := bson.M{
		"_id":     taskDocID,
		"user_id": userDocID,
	}

	// Delete the task matching the filter
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error deleting task:", err)
		return err
	}

	return nil
}
