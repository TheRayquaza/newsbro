FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main src/cmd/main.go

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /app/main .
USER nonroot:nonroot

CMD ["./main"]