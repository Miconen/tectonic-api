# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /api

# Install sqlc CLI tool
        RUN CGO_ENABLED=0 go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Generate go code using sqlc
RUN sqlc generate -f database/sqlc.yaml

# Build the Go application
RUN go build -o main .

# Use scratch to reduce image size
FROM scratch

# Copy executable
COPY --from=builder /api/main /main

# Expose the port that the application runs on
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/main"]
