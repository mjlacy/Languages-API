# Stage 1: Build the Go application
FROM golang:1.25.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app .

# Stage 2: Run the unit tests
FROM builder AS tester

WORKDIR /app

RUN go test -v ./...
