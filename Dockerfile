# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy only go.mod first (for better caching)
COPY go.mod ./

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o redis-server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy just the binary from builder
COPY --from=builder /app/redis-server .

# Create directory for AOF file
RUN mkdir -p /app/data

# Expose Redis port
EXPOSE 6379

CMD ["./redis-server"]