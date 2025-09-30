## Srvc-Scrapping

### Build

```bash
# Build the application
go build -o bin/api src/cmd/main.go

# Build with specific OS/Architecture
GOOS=linux GOARCH=amd64 go build -o bin/api-linux src/cmd/main.go

# Build for production (optimized)
go build -ldflags="-w -s" -o bin/api src/cmd/main.go
```

### Dev

```bash
# Run in development mode
go run src/cmd/main.go

# Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Download and tidy dependencies
go mod tidy
go mod download
```

## Test

```bash
## Test with race condition and coverage
go test -v -cover -race  -coverprofile=coverage.out ./...

## Analyze coverage
go tool cover -html=coverage.out
```

## Tools

```bash
# Install linting tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Security audit
go install github.com/securecodewarrior/govulncheck@latest
govulncheck ./...

# Generate documentation
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
```

## Variables

| Name                 | Value                                 |
| -------------------- | ------------------------------------- |
| DATABASE_USERNAME    | `username`                            |
| DATABASE_PASSWORD    | `password`                            |
| DATABASE_HOST        | `localhost:5432`                      |
| DATABASE_NAME        | `repo_article`                        |
| KAFKA_BROKERS        | `localhost:9092`                      |
| KAFKA_TOPIC          | `new-article-command`                 |
| RSS_FEED_URLS        | `link1,link2`                         |
