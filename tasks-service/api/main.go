package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nesistor/backend_todo/tasks-service/data"
)

const (
	webPort  = "8080"
	mongoURL = "mongodb://mongo:27017"
)

var FirebaseApp *firebase.App

var AuthClient *auth.Client
var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting task-service")

	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	credentialsFilePath := os.Getenv("FIREBASE_CREDENTIALS")
	if credentialsFilePath == "" {
		log.Fatal("FIREBASECREDENTIALS environment variable not set.")
	}

	opt := option.WithCredentialsFile(credentialsFilePath)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}
	FirebaseApp = firebaseApp

	authClient, err := FirebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Auth client: %v", err)
	}
	AuthClient = authClient

	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

// Middleware to check if the request is authenticated.
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idToken := r.Header.Get("Authorization")
		if idToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		idToken = strings.TrimPrefix(idToken, "Bearer ")

		token, err := AuthClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", token.UID)

		next(w, r)
	}
}
func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
