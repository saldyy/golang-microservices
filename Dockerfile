FROM golang:1.22 as builder
# Use the official Golang image as base

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules dependency file
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

RUN ls -la

# Build the Go binary
RUN CGO_ENABLED=0 go build -o ./dist/main ./cmd/main.go

# Start a new stage from scratch
FROM debian:buster-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /app/dist/main .
COPY --from=builder /app/.env .       

# Expose the port your application runs on
EXPOSE 8080

# Command to run the executable
RUN ls -la
CMD ["./main"]
