package main

import (
	"fmt"
	"github.com/dunglas/mercure"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	// Tworzymy router chi
	r := chi.NewRouter()

	// Inicjalizujemy serwer Mercure z konfiguracją
	mercureServer := mercure.NewServer(mercure.Config{
		JWTKey: []byte("super-secure-key"), // Klucz JWT do autoryzacji
	})

	// Endpoint Mercure do wysyłania wiadomości
	r.Post("/publish", publishHandler(mercureServer))

	// Endpoint Mercure do subskrybowania kanału
	r.Get("/subscribe", subscribeHandler(mercureServer))

	// Uruchamiamy serwer na porcie 8080
	fmt.Println("Serwer jest uruchomiony na porcie 8080...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Błąd podczas uruchamiania serwera:", err)
	}
}

func publishHandler(mercureServer *mercure.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Tutaj możesz obsłużyć żądanie POST do publikacji wiadomości w serwisie Mercure.
		// Odczytaj dane z żądania, zaimplementuj logikę publikacji i użyj biblioteki Mercure, aby wysłać wiadomość.

		// Przykładowa implementacja:
		// message := r.FormValue("message")
		// mercureServer.Publish("channel-name", []byte(message)) // Publikuj wiadomość na określonym kanale

		// Zwróć odpowiedź
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wiadomość została opublikowana"))
	}
}

func subscribeHandler(mercureServer *mercure.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Tutaj możesz obsłużyć żądanie GET do subskrybowania kanału w serwisie Mercure.
		// Przykładowa implementacja:
		// token, err := mercureServer.GenerateJWT("channel-name")
		// if err != nil {
		//     w.WriteHeader(http.StatusInternalServerError)
		//     w.Write([]byte("Błąd generowania tokenu JWT"))
		//     return
		// }

		// Zwróć token JWT do subskrybowania kanału
		w.WriteHeader(http.StatusOK)
		// w.Write([]byte(token))
	}
}
