package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var dbTimeout = 3 * time.Second

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		GPTEntry: GPTEntry{},
	}
}

type Models struct {
	GPTEntry GPTEntry
}

type GPTEntry struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string    `bson:"user_id," json:"user_id,"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

func (g *GPTEntry) InsertArticle(entry GPTEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), GPTEntry{
		Title:       entry.Title,
		Description: entry.Description,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs", err)
		return err
	}

	return nil
}

func (g *GPTEntry) GetOneArticle(id string) (*GPTEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := client.Database("articles").Collection("articles")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry GPTEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (g *GPTEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := client.Database("articles").Collection("articles")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil

}

func (g *GPTEntry) DeleteArticle(articleID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := client.Database("articles").Collection("articles")

	articleDocID, err := primitive.ObjectIDFromHex(articleID)
	if err != nil {
		log.Println("Error converting articleID:", err)
		return err
	}

	userDocID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Error converting userID:", err)
		return err
	}

	filter := bson.M{
		"_id":     articleDocID,
		"user_id": userDocID,
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error deleting task:", err)
		return err
	}

	return nil

}
