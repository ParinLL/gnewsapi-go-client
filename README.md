# Gnews API Go Client

A simple Go command-line client for fetching top headlines from Gnews API.

## Use With OpenClaw

This repository includes a publish-ready skill file at [`clawhub-publish/SKILL.md`](./clawhub-publish/SKILL.md).

ClawHub page:

- https://clawhub.ai/ParinLL/gnews

Recommended automated install (registry):

```bash
clawhub install gnews
```

After installation, set:

```bash
export GNEWS_API_KEY="your_api_key_here"
```

## Prerequisites

- Go 1.26 or later
- A Gnews API key (get one at [gnews.io](https://gnews.io))

## Usage

Set your Gnews API key as an environment variable:

```bash
export GNEWS_API_KEY="your_api_key_here"
```

### Windows (PowerShell)
```powershell
$env:GNEWS_API_KEY="your_api_key_here"
```

Run the application:

```bash
go run main.go
```

### Command-Line Flags

The application supports the following command-line flags:

- `--help` or `-h`: Show usage instructions.
- `--debug`: Enable debug mode, which prints the constructed API URLs with `apikey` redacted and raw error responses.
- `--country`: Country code for news (optional, default: `tw`).
- `--category`: Comma-separated list of categories (optional). Example: `world,technology,business`
- `--max`: Maximum number of articles to fetch (optional).

Example:
```bash
go run main.go --country tw --category world,technology --max 10 --debug
```

## Build and Install (Linux)

You can build the executable and move it to a directory in your system's `PATH` (such as `/usr/local/bin`) to run it from anywhere.

```bash
# Build the binary
go build -o gnews-client main.go

# Move the binary to a directory in your PATH
sudo mv gnews-client /usr/local/bin/

# Now you can run it globally
gnews-client
```

## Docker

Build the Docker image:

```bash
docker build -t gnewsapi-go-client .
```

Run the Docker image:

```bash
docker run -e GNEWS_API_KEY="your_api_key_here" gnewsapi-go-client --country tw --category technology,business --max 10
```

## Publish Skill

To publish/update this skill for OpenClaw users:

```bash
git add clawhub-publish/SKILL.md README.md
git commit -m "Update ClawHub skill and README"
git push origin main
clawhub publish ./clawhub-publish --slug gnews --name "GNews CLI" --version <next-version> --tags "latest,news,cli,gnews" --changelog "<your changelog>"
```

After pushing, re-sync/re-publish from your OpenClaw/ClawHub side if your workspace requires an explicit publish step.
