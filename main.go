package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type NewsResponse struct {
	TotalArticles int       `json:"totalArticles"`
	Articles      []Article `json:"articles"`
}

type Article struct {
	Source struct {
		Name string  `json:"name"`
		URL  string  `json:"url"`
	} `json:"source"`
	Author      *string   `json:"author"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
}

func main() {
	apiKey := os.Getenv("GNEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("GNEWS_API_KEY environment variable is not set. Please set it to your Gnews API key.")
	}

	country := os.Getenv("NEWS_COUNTRY")
	if country == "" {
		country = "tw" // default
	}

	categoryEnv := os.Getenv("NEWS_CATEGORY")
	var categories []string
	if categoryEnv != "" {
		for _, c := range strings.Split(categoryEnv, ",") {
			trimmed := strings.TrimSpace(c)
			if trimmed != "" {
				categories = append(categories, trimmed)
			}
		}
	}
	
	if len(categories) == 0 {
		categories = []string{""} // Add an empty category to run the loop once without category filtering
	}

	for _, category := range categories {
		url := fmt.Sprintf("https://gnews.io/api/v4/top-headlines?country=%s&apikey=%s", country, apiKey)
		if category != "" {
			url += fmt.Sprintf("&category=%s", category)
			fmt.Printf("--- Top Headlines for Category: %s ---\n", category)
		} else {
			url += "&category=general"
			fmt.Printf("--- Top Headlines ---\n")
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Error creating request for category %s: %v", category, err)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request for category %s: %v", category, err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}

		var newsResp NewsResponse
		if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
			resp.Body.Close()
			log.Fatalf("Error decoding JSON response: %v", err)
		}
		resp.Body.Close()

		fmt.Printf("Fetched %d articles\n\n", len(newsResp.Articles))

		for i, article := range newsResp.Articles {
			fmt.Printf("%d. %s\n", i+1, article.Title)
			fmt.Printf("   Source: %s\n", article.Source.Name)
			fmt.Printf("   URL: %s\n\n", article.URL)
		}
	}
}
