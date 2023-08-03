package main

import (
	"errors"
	"net/http"

	"github.com/nesistor/backend_todo/task-service/data"
)

type JSONPayload struct {
	TaskID      string `json:"task_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

var tasks = make(map[string]JSONPayload)

func (app *Config) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data
	event := data.Task{
		UserID:      userID,
		ID:          requestPayload.TaskID,
		Title:       requestPayload.Title,
		Description: requestPayload.Description,
	}

	err := app.Models.Task.InsertTask(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "added task succesfully",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *Config) RemoveTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	taskID := r.URL.Query().Get("task_id")

	// Call the DeleteTask method to remove the task
	err := app.Models.Task.DeleteTask(taskID, userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Task removed",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	taskID := r.URL.Query().Get("task_id")

	var requestPayload JSONPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Validate that the taskID in the payload matches the taskID in the query parameters
	if taskID != requestPayload.TaskID {
		app.errorJSON(w, errors.New("task_id mismatch"))
		return
	}

	// Construct the updated task object
	updatedTask := data.Task{
		UserID:      userID,
		ID:          requestPayload.TaskID,
		Title:       requestPayload.Title,
		Description: requestPayload.Description,
	}

	// Call the UpdateTask method to update the task
	err = app.Models.Task.UpdateTask(updatedTask)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Task updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	// Call the GetAllByUserID method to fetch all tasks for the given user ID
	tasks, err := app.Models.Task.GetAllByUserID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Respond with the tasks in the JSON format
	app.writeJSON(w, http.StatusOK, tasks)
}
