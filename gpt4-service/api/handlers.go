package main

import (
	"encoding/json"
	"fmt"
	"gpt4-service/data"
	"net/http"
)

type JSONPayload struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (app *Config) WriteArticle(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	title, article, err := callGPT3API(requestPayload.Title, requestPayload.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Error:       false,
		Title:       title,
		Description: article,
	}

	// Convert the response payload to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
