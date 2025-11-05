FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main src/cmd/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash=5.2.37-2+b5 \
    ca-certificates=20250419 \
    curl=8.14.1-2 \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m -s /bin/bash nonroot

WORKDIR /app
COPY --from=builder /app/main .

USER nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://${APP_HOST}:${PORT}/health || exit 1

CMD ["./main"]
