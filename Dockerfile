# Development stage
FROM golang:1.24-alpine3.21 AS dev

# Install Git, swag and Air
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/air-verse/air@v1.61.7

# set working dir
WORKDIR /app

# Copy go.mod and download deps (for cache)
COPY go.mod go.sum ./
RUN go mod download

# copy the source code
COPY . .

# Generate Swagger docs
RUN swag init

# export the port
EXPOSE 8000

# Run the executable
CMD ["sh", "-c", "swag init && air"]


# Builder stage for production
FROM golang:1.24-alpine3.21 AS builder

# Install Git, curl and swag (required for swagger generation during build)
RUN apk add --no-cache git curl
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Download golang-migrate CLI tool for database migrations in production
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# set working dir
WORKDIR /app

# Copy go.mod and download deps (for cache)
COPY go.mod go.sum ./
RUN go mod download

# copy the source code
COPY . .

# Generate Swagger docs
RUN swag init

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Final runtime production stage
FROM alpine:3.21 AS release

# Install ca-certificates (useful for https requests in production) and bash
RUN apk add --no-cache ca-certificates bash

WORKDIR /app

# Copy application binary, migrations, docs, and migrate CLI
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Export port
EXPOSE 8000

# Run binary by default
CMD ["./main"]
