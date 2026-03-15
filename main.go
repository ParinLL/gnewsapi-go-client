package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	debug := flag.Bool("debug", false, "Enable debug mode (print requested URLs)")
	
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  --debug        Enable debug mode (print requested URLs with redacted API key)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --help, -h     Show this help message\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Environment variables:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  GNEWS_API_KEY  Your Gnews API key (required)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  NEWS_COUNTRY   Country code for news (default: tw)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  NEWS_MAX       Maximum number of articles to fetch\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  NEWS_CATEGORY  Comma-separated list of categories (e.g., world,technology)\n")
	}
	
	flag.Parse()

	apiKey := os.Getenv("GNEWS_API_KEY")
	if apiKey == "" {
		log.Println("GNEWS_API_KEY environment variable is not set. Please set it to your Gnews API key.")
		flag.Usage()
		os.Exit(1)
	}

	country := os.Getenv("NEWS_COUNTRY")
	if country == "" {
		country = "tw" // default
	}

	maxEnv := os.Getenv("NEWS_MAX")

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
		if maxEnv != "" {
			url += fmt.Sprintf("&max=%s", maxEnv)
		}
		if category != "" {
			url += fmt.Sprintf("&category=%s", category)
			fmt.Printf("--- Top Headlines for Category: %s ---\n", category)
		} else {
			url += "&category=general"
			fmt.Printf("--- Top Headlines ---\n")
		}

		if *debug {
			fmt.Printf("[DEBUG] Request URL: %s\n", redactAPIKey(url))
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
			if *debug {
				fmt.Printf("[DEBUG] Raw Error Response: %s\n", string(body))
			}
			log.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}

		if *debug {
			fmt.Printf("[DEBUG] Request successful, status code: %d\n", resp.StatusCode)
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

func redactAPIKey(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	q := parsed.Query()
	if q.Get("apikey") != "" {
		q.Set("apikey", "REDACTED")
		parsed.RawQuery = q.Encode()
	}

	return parsed.String()
}
