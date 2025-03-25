# Use a lightweight base image
FROM golang:1.23.0-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIN_MODE=release

# Set the working directory
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Change to the directory containing the main.go file
WORKDIR /app/cmd

# Build the binary
RUN go build -o /app/main

# Set the working directory back to /app
WORKDIR /app

# Expose the application port
EXPOSE 8443

# Start the application
CMD ["./main"]
