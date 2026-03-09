# Gnews API Go Client

A simple Go command-line client for fetching top headlines from Gnews API.

## Prerequisites

- Go 1.26 or later
- A Gnews API key (get one at [gnews.io](https://gnews.io))

## Usage

Set your Gnews API key as an environment variable:

```bash
export GNEWS_API_KEY="your_api_key_here"
export NEWS_COUNTRY="tw" # optional, defaults to 'tw'
export NEWS_CATEGORY="technology,business" # optional, comma-separated list
```

### Windows (PowerShell)
```powershell
$env:GNEWS_API_KEY="your_api_key_here"
$env:NEWS_COUNTRY="tw"
$env:NEWS_CATEGORY="world,technology,business"
```

Run the application:

```bash
go run main.go
```

## Docker

Build the Docker image:

```bash
docker build -t gnewsapi-go-client .
```

Run the Docker image:

```bash
docker run -e GNEWS_API_KEY="your_api_key_here" -e NEWS_COUNTRY="tw" -e NEWS_CATEGORY="technology,business" gnewsapi-go-client
```
