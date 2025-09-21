package tests

import (
	"net/http"
	"testing"
)

func testUserEndpoints(t *testing.T, authToken string) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}

	t.Run("Get User Profile", func(t *testing.T) {
		resp, err := makeRequest("GET", baseURL+"/user/profile", nil, headers)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 200, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		var user User
		if err := parseResponse(resp, &user); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if user.Email == "" {
			t.Fatal("Expected user email in response")
		}
		if user.ID == 0 {
			t.Fatal("Expected user ID in response")
		}
	})

	t.Run("Get Profile Without Auth", func(t *testing.T) {
		resp, err := makeRequest("GET", baseURL+"/user/profile", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 without auth, got %d", resp.StatusCode)
		}
	})

	t.Run("Update User Profile", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"name": "Updated Test User",
		}

		resp, err := makeRequest("PUT", baseURL+"/user/profile", updateReq, headers)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 200, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		var user User
		if err := parseResponse(resp, &user); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if user.Name != "Updated Test User" {
			t.Fatalf("Expected name 'Updated Test User', got '%s'", user.Name)
		}
	})

	t.Run("Update Profile With Invalid Data", func(t *testing.T) {
		updateReq := "invalid json"

		resp, err := makeRequest("PUT", baseURL+"/user/profile", updateReq, headers)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status 400 for invalid JSON, got %d", resp.StatusCode)
		}
	})

	t.Run("Get All Users", func(t *testing.T) {
		resp, err := makeRequest("GET", baseURL+"/users", nil, headers)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// This might return 403 if user is not admin, or 200 if they are
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
			var errorResp ErrorResponse
			parseResponse(resp, &errorResp)
			t.Fatalf("Expected status 200 or 403, got %d. Error: %s", resp.StatusCode, errorResp.Error)
		}

		if resp.StatusCode == http.StatusOK {
			var response map[string]interface{}
			if err := parseResponse(resp, &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if _, exists := response["users"]; !exists {
				t.Fatal("Expected 'users' field in response")
			}
		}
	})

	t.Run("Get Users With Pagination", func(t *testing.T) {
		url := baseURL + "/users?limit=5&offset=0"
		resp, err := makeRequest("GET", url, nil, headers)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
			t.Fatalf("Expected status 200 or 403, got %d", resp.StatusCode)
		}
	})
}
