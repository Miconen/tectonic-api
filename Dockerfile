# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /api

# Install sqlc CLI tool
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Generate go code using sqlc
RUN sqlc generate -f database/sqlc.yaml

# Build the Go application
RUN go build -o main .

# Expose the port that the application runs on
EXPOSE 8080

# Command to run the application
RUN ["chmod", "+x", "./wait-for-it.sh"]
CMD ./wait-for-it.sh 0.0.0.0:5432 -- ./main
