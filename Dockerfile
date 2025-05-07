# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set required build flags for SQLite
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE -D_FILE_OFFSET_BITS=64"

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Run swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g main.go

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create directory for database
RUN mkdir -p /data

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy the config file
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Set environment variable for database path
ENV DB_PATH=/data/app.db

# Create a non-root user
RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app /data

# Switch to non-root user
USER appuser

# Command to run the executable
CMD ["./main"] 