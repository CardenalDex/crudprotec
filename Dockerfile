# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

# Install build dependencies: 
# build-base contains gcc, make, and libc-dev required for CGO
RUN apk add --no-cache git build-base

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/app/main.go -d ./ --parseDependency --output docs

# Build the application with CGO enabled
# We use -tags musl if needed, but standard build usually works on alpine with build-base
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/app/main.go

# Stage 2: Final Image
FROM alpine:latest
# We need to install the shared libraries that the binary was linked against (libc)
RUN apk --no-cache add ca-certificates tzdata libc6-compat

WORKDIR /app

# Create the data directory for the internal SQLite database
RUN mkdir -p /app/data

# Copy the binary and docs from builder
COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./server"]