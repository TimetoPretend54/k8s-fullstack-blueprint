# Stage 1: Build stage
FROM golang:1.22-alpine AS builder

# Install git (required for go mod download in some cases)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
# This layer will only rebuild when dependencies change
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
# - CGO_ENABLED=0: Disable CGO for static binary
# - GOOS=linux: Target Linux OS
# - GOARCH=amd64: Target amd64 architecture
# -ldflags="-w -s": Strip debug info to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a \
    -installsuffix cgo \
    -o server .

# Stage 2: Runtime stage
FROM alpine:3.19

# Install ca-certificates for TLS
RUN apk add --no-cache ca-certificates

# Create a non-root user for security
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Set the working directory
WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/server .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]
