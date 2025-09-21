package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAuthEndpoints(t *testing.T) {
	testUser := TestUser{
		Email:    fmt.Sprintf("test_%d@example.com", time.Now().Unix()),
		Password: "testpassword123",
	}

	t.Run("Register User", func(t *testing.T) {
		resp, err := makeRequest("POST", baseURL+"/auth/register", testUser, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 201, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		var authResp AuthResponse
		if err := parseResponse(resp, &authResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if authResp.Token == "" {
			t.Fatal("Expected token in response")
		}
		if authResp.User.Email != testUser.Email {
			t.Fatalf("Expected email %s, got %s", testUser.Email, authResp.User.Email)
		}
	})

	t.Run("Register Duplicate User", func(t *testing.T) {
		resp, err := makeRequest("POST", baseURL+"/auth/register", testUser, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 for duplicate registration, got %d", resp.StatusCode)
		}
	})

	var authToken, refreshToken string

	t.Run("Login User", func(t *testing.T) {
		loginReq := TestUser{
			Email:    testUser.Email,
			Password: testUser.Password,
		}

		resp, err := makeRequest("POST", baseURL+"/auth/login", loginReq, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 200, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		var authResp AuthResponse
		if err := parseResponse(resp, &authResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if authResp.Token == "" {
			t.Fatal("Expected token in response")
		}
		if authResp.RefreshToken == "" {
			t.Fatal("Expected refresh token in response")
		}

		authToken = authResp.Token
		refreshToken = authResp.RefreshToken
	})

	t.Run("Login Invalid Credentials", func(t *testing.T) {
		invalidLogin := TestUser{
			Email:    testUser.Email,
			Password: "wrongpassword",
		}

		resp, err := makeRequest("POST", baseURL+"/auth/login", invalidLogin, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 for invalid credentials, got %d", resp.StatusCode)
		}
	})

	t.Run("Refresh Token", func(t *testing.T) {
		refreshReq := map[string]string{
			"refresh_token": refreshToken,
		}

		resp, err := makeRequest("POST", baseURL+"/auth/refresh", refreshReq, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 200, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		var authResp AuthResponse
		if err := parseResponse(resp, &authResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if authResp.Token == "" {
			t.Fatal("Expected new token in response")
		}

		// Update token for subsequent tests
		authToken = authResp.Token
	})

	t.Run("Refresh Invalid Token", func(t *testing.T) {
		refreshReq := map[string]string{
			"refresh_token": "invalid_token",
		}

		resp, err := makeRequest("POST", baseURL+"/auth/refresh", refreshReq, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 for invalid refresh token, got %d", resp.StatusCode)
		}
	})

	// Run user tests with the auth token
	testUserEndpoints(t, authToken)
}

func TestOAuthEndpoints(t *testing.T) {
	t.Run("OAuth Login", func(t *testing.T) {
		resp, err := makeRequest("GET", baseURL+"/auth/oauth/login", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should return either 200 with auth URL or 503 if not configured
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 200 or 503, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		if resp.StatusCode == http.StatusOK {
			var response map[string]interface{}
			if err := parseResponse(resp, &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if authURL, exists := response["auth_url"]; !exists || authURL == "" {
				t.Fatal("Expected 'auth_url' in response")
			}
		}
	})

	t.Run("OAuth Callback Without Parameters", func(t *testing.T) {
		resp, err := makeRequest("GET", baseURL+"/auth/oauth/callback", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 without parameters, got %d", resp.StatusCode)
		}
	})

	t.Run("OAuth Callback With Invalid State", func(t *testing.T) {
		url := baseURL + "/auth/oauth/callback?code=test&state=invalid"
		resp, err := makeRequest("GET", url, nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 with invalid state, got %d", resp.StatusCode)
		}
	})
}

func TestInputValidation(t *testing.T) {
	t.Run("Register With Missing Fields", func(t *testing.T) {
		invalidUser := map[string]string{
			"email": "test@example.com",
			// missing password
		}

		resp, err := makeRequest("POST", baseURL+"/auth/register", invalidUser, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 for missing fields, got %d", resp.StatusCode)
		}
	})

	t.Run("Register With Invalid Email", func(t *testing.T) {
		invalidUser := TestUser{
			Email:    "invalid-email",
			Password: "password123",
		}

		resp, err := makeRequest("POST", baseURL+"/auth/register", invalidUser, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 for invalid email, got %d", resp.StatusCode)
		}
	})

	t.Run("Login With Empty Body", func(t *testing.T) {
		resp, err := makeRequest("POST", baseURL+"/auth/login", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 for empty body, got %d", resp.StatusCode)
		}
	})
}
