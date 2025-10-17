package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// --- Configuration ---
const (
	API_ENDPOINT      = "https://article.newsbro.cc/api/v1/rss" 
	AUTH_TOKEN_HEADER = "Bearer xxx"
	YAML_FILE         = "rss.yaml"
)

// --- DTOs for API Communication ---

// RSSCreateRequest for both MetaFeeds and actual RssFeeds
type RSSCreateRequest struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name,omitempty"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
	Parents     []string `json:"parents"`
}
type RSSResponse struct {
	Name string `json:"name"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

type YamlRegistry struct {
	MetaFeeds []MetaFeed `yaml:"rss_meta_feeds"`
	Feeds     []RssFeed  `yaml:"rss_feeds"`
}
type MetaFeed struct {
	Name        string   `yaml:"name"`
	DisplayName string   `yaml:"display_name"`
	Description string   `yaml:"description"`
	Parents     []string `yaml:"parents"`
}
type RssFeed struct {
	Name        string   `yaml:"name"`
	DisplayName string   `yaml:"display_name"`
	Link        string   `yaml:"link"`
	Description string   `yaml:"description"`
	Parents     []string `yaml:"parents"`
}

// --- API Execution Helper ---

func callCreateRSS(client *http.Client, reqBody RSSCreateRequest) error {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body for '%s': %w", reqBody.Name, err)
	}

	req, err := http.NewRequest("POST", API_ENDPOINT, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request for '%s': %w", reqBody.Name, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", AUTH_TOKEN_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API Request failed for '%s': %w", reqBody.Name, err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusCreated {
		var rssResponse RSSResponse
		if err := json.Unmarshal(responseBody, &rssResponse); err != nil {
			log.Printf("✅ Success, but failed to unmarshal response for '%s'. Response: %s", reqBody.Name, string(responseBody))
		} else {
			log.Printf("✅ Created '%s'", rssResponse.Name)
		}
		return nil
	} else if resp.StatusCode == http.StatusConflict {
		var errResponse ErrorResponse
		json.Unmarshal(responseBody, &errResponse)
		log.Printf("⚠️  Skipped '%s' (Status: %s). Reason: %s", reqBody.Name, resp.Status, errResponse.Error)
		return nil
	} else {
		var errResponse ErrorResponse
		if err := json.Unmarshal(responseBody, &errResponse); err != nil {
			return fmt.Errorf("failed (Status: %s). Response: %s", resp.Status, string(responseBody))
		}
		return fmt.Errorf("failed (Status: %s). Reason: %s", resp.Status, errResponse.Error)
	}
}

func main() {
	log.SetFlags(0)
	log.Println("--- Starting RSS Registry API Ingestion Script ---")

	data, err := os.ReadFile(YAML_FILE)
	if err != nil {
		log.Fatalf("❌ Error reading YAML file %s: %v", YAML_FILE, err)
	}

	var registry YamlRegistry
	if err := yaml.Unmarshal(data, &registry); err != nil {
		log.Fatalf("❌ Error unmarshalling YAML: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	log.Println("\n==================================================")
	log.Printf("PASS 1: Creating %d Meta-Feeds (Categories)...", len(registry.MetaFeeds))
	log.Println("==================================================")

	for _, meta := range registry.MetaFeeds {
		reqBody := RSSCreateRequest{
			Name:        meta.Name,
			Description: meta.Description,
			DisplayName: meta.DisplayName,
			Link:        "", 
			Parents:     meta.Parents,
		}
		if err := callCreateRSS(client, reqBody); err != nil {
			log.Printf("❌ Meta-Feed Error for '%s': %v", meta.Name, err)
		}
	}
	
	// =======================================================
	// PASS 2: Create Actual RSS Feeds (Children)
	// =======================================================
	log.Println("\n==================================================")
	log.Printf("PASS 2: Creating %d Actual RSS Feeds...", len(registry.Feeds))
	log.Println("==================================================")
	
	for _, feed := range registry.Feeds {
		reqBody := RSSCreateRequest{
			Name:        feed.Name,
			Description: feed.Description,
			DisplayName: feed.DisplayName,
			Link:        feed.Link, 
			Parents:     feed.Parents,
		}
		if err := callCreateRSS(client, reqBody); err != nil {
			log.Printf("❌ RSS Feed Error for '%s': %v", feed.Name, err)
		}
	}

	log.Println("\n--- Script Finished ---")
}
