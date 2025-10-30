## Repo-Feed

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

## Generate swagger
go install github.com/swaggo/swag/cmd/swag@latest
# nix-shell -p go-swag # or this
swag init -g src/cmd/main.go -o docs --parseDependency --parseInternal
```

## Variables

| Name                          | Description                                              | Default Value                            |
|-------------------------------|----------------------------------------------------------|------------------------------------------|
| `PORT`                        | HTTP server port.                                        | `8080`                                   |
| `JWT_SECRET`                  | Secret key for signing JWTs.                             | `your-secret-key`                        |
| `FRONTEND_ORIGIN`             | The origin URL of the frontend application.              | `http://localhost:3000`                  |
| `LOGIN_REDIRECT_URL`          | URL for post-login redirection.                          | `http://localhost:3000/auth/callback`    |
| `REDIS_SENTINELS`             | Comma-separated list of Redis Sentinel addresses.        | `localhost:26379`                        |
| `REDIS_MASTER_NAME`           | Name of the Redis master instance.                       | `mymaster`                               |
| `REDIS_PASSWORD`              | Password for Redis.                                      | *(empty)*                                |
| `REDIS_DB`                    | Redis database index.                                    | `0`                                      |
| `KAFKA_BROKERS`               | Comma-separated list of Kafka broker addresses.          | `localhost:9092`                         |
| `KAFKA_INFERENCE_COMMAND_TOPIC` | Kafka topic for sending inference commands.            | `inference-commands`                     |
| `KAFKA_FEEDBACK_TOPIC`        | Kafka topic for receiving user feedback.                 | `user-feedback`                          |
| `KAFKA_GROUP_ID`              | Consumer group ID for inference commands.                | `repo-feed-group`                        |
| `KAFKA_FEEDBACK_GROUP_ID`     | Consumer group ID for user feedback.                     | `repo-feed-feedback-group`               |
| `DEFAULT_MODEL`               | Default recommendation model to use.                     | `tfidf`                                  |
| `FEEDBACK_EXPIRATION_SECONDS` | Duration (in seconds) to keep user feedback.             | `604800 (7 days)`                        |
| `FEED_RESCORING_ENABLED`      | Enable/disable the feed rescoring background job.        | `false`                                  |
| `FEED_RESCORING_INTERVAL`     | Frequency of the feed rescoring job.                     | `10m`                                    |
| `FEED_RESCORING_BATCH_SIZE`   | Number of feeds to process per batch.                    | `1000`                                   |
| `FEED_RESCORING_WORKERS`      | Number of concurrent workers for rescoring.              | `10`                                     |
| `FEED_RESCORING_THRESHOLD`    | Score below which articles are removed from ZSETs.       | `0.1`                                    |
| `FEED_RESCORING_DECAY_ENABLED` | Enable/disable score decay during rescoring.            | `true`                                   |
| `FEED_RESCORING_DECAY_HALFLIFE` | Half-life duration for score decay (e.g., 24h).        | `60h (2.5 days)`                         |
| `ENVIRONMENT`                 | Application environment (e.g., prod, dev).               | `prod`                                   |
