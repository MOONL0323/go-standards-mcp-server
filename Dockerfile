# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates git

# Install golangci-lint
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mcp-server .

# Copy config files
COPY configs/ ./configs/

# Create directories for reports and temp files
RUN mkdir -p /app/reports /app/tmp

# Expose port for HTTP mode
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/mcp-server", "health"] || exit 1

# Run the application
ENTRYPOINT ["/app/mcp-server"]
CMD ["--mode", "http", "--port", "8080"]
