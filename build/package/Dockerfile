# Stage 1: Build the Go application
FROM golang:1.25.2-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app .

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /root

EXPOSE 8080

COPY --from=builder / .

ENTRYPOINT ["./app"]
