## Repo-Account

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

| Name                       | Value                                     |
| -------------------------- | ----------------------------------------- |
| DATABASE_USERNAME          | `username`                                |
| DATABASE_PASSWORD          | `password`                                |
| DATABASE_HOST              | `localhost:5432`                          |
| DATABASE_NAME              | `repo_account`                            |
| JWT_SECRET                 | `your-secret-key`                         |
| OIDC_ISSUER_URL            | `http://localhost:8080`                   |
| OIDC_CLIENT_ID             | *(empty)*                                 |
| OIDC_CLIENT_SECRET         | *(empty)*                                 |
| OIDC_REDIRECT_URL          | `http://localhost:8080/auth/callback`     |
| OIDC_REDIRECT_FRONTEND_URL |                                           |
| COOKIE_DOMAIN              | ` `                                       |
| LOGIN_REDIRECT_URL         | `http://localhost:8080/api/v1/auth/login` |
