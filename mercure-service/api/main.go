package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dunglas/mercure"
	"github.com/go-chi/chi"
)

const (
	webPort  = "8080"
	mongoURL = "mongodb://mongo:27017"
)

var FirebaseApp *firebase.App

var AuthClient *auth.Client

func main() {
	log.Println("Starting mercure-service")
	hub := mercure.New(
		mercure.WithAllowSubscribe(nil),
	)

	// Inicjalizacja routeru Chi
	r := chi.NewRouter()

	// Endpoint do autentykacji i tworzenia kanału Mercure dla użytkownika
	r.Post("/create-channel/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		// Tworzenie kanału dla użytkownika
		channelName := fmt.Sprintf("user/%s", userID)

		// Autoryzacja użytkownika i utworzenie subskrypcji do jego kanału
		// Tutaj możesz dodać logikę autentykacji i autoryzacji użytkownika

		// Przygotuj wiadomość powitalną
		welcomeMessage := mercure.NewEvent(channelName, []byte("Witaj na swoim kanale Mercure!"))

		// Wyślij wiadomość powitalną do użytkownika
		if err := hub.Publish(welcomeMessage); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Kanał Mercure dla użytkownika %s został utworzony", userID)
	})

	// Endpoint do subskrybowania kanału Mercure przez użytkownika
	r.Get("/subscribe/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		// Tutaj możesz dodać logikę autentykacji i autoryzacji użytkownika

		// Subskrybuj użytkownika do jego kanału
		channelName := fmt.Sprintf("user/%s", userID)
		hub.Subscribe(w, r, channelName)
	})

	// Obsługa subskrypcji klientów
	r.Handle("/updates", hub)

	// Uruchom serwer na porcie 3000
	go func() {
		log.Fatal(http.ListenAndServe(":3000", r))
	}()

	// Komunikat o rozpoczęciu działania serwera
	log.Println("Serwer Mercure jest uruchomiony na porcie 3000")

	// Pozwól serwerowi działać
	select {}
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
