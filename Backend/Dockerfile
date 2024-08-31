# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application, specifying the path to the main file
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/backend

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/public_key.pem .

# Expose the application port
EXPOSE 8000

# Run the binary
CMD ["./main"]
