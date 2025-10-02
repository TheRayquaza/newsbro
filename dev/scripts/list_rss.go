package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GitHub API response structures
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

// OPML structures
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Body    Body     `xml:"body"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	Text     string    `xml:"text,attr"`
	Title    string    `xml:"title,attr"`
	Type     string    `xml:"type,attr"`
	XMLUrl   string    `xml:"xmlUrl,attr"`
	HTMLUrl  string    `xml:"htmlUrl,attr"`
	Outlines []Outline `xml:"outline"`
}

func main() {
	repo := "plenaryapp/awesome-rss-feeds"
	
	fmt.Println("Fetching OPML files from GitHub repository...")
	opmlFiles, err := findOPMLFiles(repo, "")
	if err != nil {
		fmt.Printf("Error finding OPML files: %v\n", err)
		return
	}

	fmt.Printf("Found %d OPML files\n\n", len(opmlFiles))

	allFeeds := make(map[string]bool)
	
	for _, file := range opmlFiles {
		fmt.Printf("Processing: %s\n", file.Path)
		feeds, err := extractFeedsFromOPML(file.DownloadURL)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		
		fmt.Printf("  Found %d feeds\n", len(feeds))
		for _, feed := range feeds {
			allFeeds[feed] = true
		}
	}

	fmt.Printf("\n=== Total unique RSS feeds: %d ===\n\n", len(allFeeds))
	
	for feed := range allFeeds {
		fmt.Println(feed)
	}
}

func findOPMLFiles(repo, path string) ([]GitHubContent, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var contents []GitHubContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return nil, err
	}

	var opmlFiles []GitHubContent
	
	for _, item := range contents {
		if item.Type == "file" && strings.HasSuffix(strings.ToLower(item.Name), ".opml") {
			opmlFiles = append(opmlFiles, item)
		} else if item.Type == "dir" {
			subFiles, err := findOPMLFiles(repo, item.Path)
			if err != nil {
				fmt.Printf("Warning: couldn't access directory %s: %v\n", item.Path, err)
				continue
			}
			opmlFiles = append(opmlFiles, subFiles...)
		}
	}
	
	return opmlFiles, nil
}

func extractFeedsFromOPML(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var opml OPML
	if err := xml.Unmarshal(body, &opml); err != nil {
		return nil, err
	}

	var feeds []string
	extractFeeds(opml.Body.Outlines, &feeds)
	
	return feeds, nil
}

func extractFeeds(outlines []Outline, feeds *[]string) {
	for _, outline := range outlines {
		if outline.XMLUrl != "" {
			*feeds = append(*feeds, outline.XMLUrl)
		}
		if len(outline.Outlines) > 0 {
			extractFeeds(outline.Outlines, feeds)
		}
	}
}