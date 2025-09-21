package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8080/api/v1"
)

type TestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper function to make HTTP requests
func makeRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

// Helper function to parse JSON response
func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
