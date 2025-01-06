# Use the official Go image for aarch64 (arm64)
FROM --platform=linux/arm64 golang:1.23.4 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

# Create a working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final build
FROM --platform=linux/arm64 alpine:latest

# Set up a working directory
WORKDIR /app

# Install necessary tools (optional, for debugging or running)
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file into the image
COPY .env .env

# Expose the application port (adjust as needed)
EXPOSE 8080

# Run the application
CMD ["./main"]
