# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install git and other build dependencies
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application API
RUN CGO_ENABLED=0 GOOS=linux go build -o innogen-backend ./cmd/main.go
# Build the seeder
RUN CGO_ENABLED=0 GOOS=linux go build -o innogen-seed ./cmd/seed/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy the binaries from the builder stage
COPY --from=builder /app/innogen-backend .
COPY --from=builder /app/innogen-seed .

# Copy environment file
# You will mount this file in docker-compose.yml on the host
# COPY .env .env

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./innogen-backend"]
