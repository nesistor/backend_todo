package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type jsonResponse struct {
	Error       bool   `json:"error"`
	Title       string `json:"message"`
	Description string `json:"description"`
	Data        any    `json:"data,omitempty"`
}

type GPTResponse struct {
	Text string `json:"text"`
}

func callGPT3API(taskTitle, taskDescription string) (string, string, error) {
	apiKey := "GPT3_API_KEY"
	url := "https://api.openai.com/v1/engines/davinci/completions"

	titlePrompt := fmt.Sprintf(`Generete me short title for my article 
	write only 5 words which will be the title for article how i can do better 
	my task which is title: %s, description: %s`, taskTitle, taskDescription)
	titlePayload := fmt.Sprintf(`{
		"prompt": "%s",
		"max_tokens": 20 
	}`, titlePrompt)

	titleReq, err := http.NewRequest("POST", url, strings.NewReader(titlePayload))
	if err != nil {
		return "", "", err
	}

	titleReq.Header.Set("Content-Type", "application/json")
	titleReq.Header.Set("Authorization", "Bearer "+apiKey)

	titleResp, err := http.DefaultClient.Do(titleReq)
	if err != nil {
		return "", "", err
	}
	defer titleResp.Body.Close()

	var titleResponse GPTResponse
	if err := json.NewDecoder(titleResp.Body).Decode(&titleResponse); err != nil {
		return "", "", err
	}

	if len(titleResponse.Text) == 0 {
		return "", "", fmt.Errorf("No title response from GPT-3")
	}

	articlePrompt := fmt.Sprintf(`Write me an article about my task and tell me how I can do it smarter, better, and cheaper.
Task title is: %s, Description of the task is: %s
Title for the article is: %s, so think more about writing about the theme of the Title article`, taskTitle, taskDescription, titleResponse.Text)

	articlePayload := fmt.Sprintf(`{
		"prompt": "%s",
		"max_tokens": 500 
	}`, articlePrompt)

	articleReq, err := http.NewRequest("POST", url, strings.NewReader(articlePayload))
	if err != nil {
		return "", "", err
	}

	articleReq.Header.Set("Content-Type", "application/json")
	articleReq.Header.Set("Authorization", "Bearer "+apiKey)

	articleResp, err := http.DefaultClient.Do(articleReq)
	if err != nil {
		return "", "", err
	}
	defer articleResp.Body.Close()

	var articleResponse GPTResponse
	if err := json.NewDecoder(articleResp.Body).Decode(&articleResponse); err != nil {
		return "", "", err
	}

	if len(articleResponse.Text) == 0 {
		return "", "", fmt.Errorf("No article response from GPT-3")
	}

	return titleResponse.Text, articleResponse.Text, nil
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Data = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
