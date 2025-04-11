# Use the official Golang image as the base image
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /api

# Install CLI tools
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install golang.org/x/tools/cmd/stringer@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Generate go code using sqlc
RUN sqlc generate -f database/sqlc.yaml
RUN go generate ./models

# Build the Go application
RUN CGO_ENABLED=0 go build -o main .

# Use scratch to reduce image size
FROM scratch

# Copy certificates and timezone info
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy executable
COPY --from=builder /api/main /main

# Expose the port that the application runs on
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/main"]
