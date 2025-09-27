package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"strconv"
	"bytes"
)

type Article struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Category string `json:"category"`
	SubCat   string `json:"subcategory"`
	Link     string `json:"link"`
}

func main() {
	// ---- Config ----
	tsvFile := "../../apps/data/news.tsv"
	apiURL := "http://localhost:8080/api/v1/articles" // replace with your endpoint
	token := "Bearer xxx"

	// ---- Open TSV ----
	file, err := os.Open(tsvFile)
	if err != nil {
		log.Fatalf("Failed to open TSV: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip header
		if lineNum == 1 {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 6 {
			log.Printf("Skipping line %d: not enough fields", lineNum)
			continue
		}

		// Remove "N" prefix and convert to uint
		idStr := strings.TrimPrefix(fields[0], "N")
		idInt, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Printf("Invalid ID on line %d: %v", lineNum, err)
			continue
		}

		article := Article{
			ID:       idInt,
			Category: fields[1],
			SubCat:   fields[2],
			Title:    fields[3],
			Abstract: fields[4],
			Link:     fields[5],
		}

		// ---- Encode JSON ----
		data, err := json.Marshal(article)
		if err != nil {
			log.Printf("Error marshaling JSON for line %d: %v", lineNum, err)
			continue
		}

		// ---- POST request ----
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error creating request for line %d: %v", lineNum, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Printf("Error dumping request for line %d: %v", lineNum, err)
			continue
		}

		fmt.Printf("REQUEST:\n%s", string(reqDump))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error posting line %d: %v", lineNum, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			log.Printf("Failed to insert line %d: HTTP %s", lineNum, resp.Status)
		} else {
			fmt.Printf("Inserted line %d: %s\n", lineNum, article.Title)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading TSV: %v", err)
	}
}
