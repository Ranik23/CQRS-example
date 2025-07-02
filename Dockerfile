# Use the official Go image
FROM golang:1.24.2

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

RUN apt-get update && apt-get install -y postgresql-client && apt-get install -y kafkacat

# Copy all source files including vendor and entrypoint
COPY . .

# Build the binary
RUN go build -mod=vendor -o main ./cmd/main.go

# Make entrypoint executable
RUN chmod +x /app/entrypoint.sh

# Expose application port
EXPOSE 8080

# Set entrypoint and default command
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./main"]
