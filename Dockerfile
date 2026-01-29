# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies and swag
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Swagger documentation
# Note: Ensure the path to main.go is correct based on your project structure
RUN swag init -g cmd/app/main.go -d ./ --parseDependency --output docs

# Build the application
# CGO_ENABLED=0 is preferred for alpine to ensure a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/app/main.go

# Stage 2: Final Image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create the data directory for the internal SQLite database
RUN mkdir -p /app/data

# Copy the binary and the generated docs from the builder
COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs

# Expose the default port
EXPOSE 8080

# Run the application
CMD ["./server"]