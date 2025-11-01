#!/bin/bash

# Quick start script for Go Standards MCP Server
# This script helps you get started quickly

set -e

echo "🚀 Go Standards MCP Server - Quick Start"
echo "========================================"
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Check golangci-lint (optional)
if command -v golangci-lint &> /dev/null; then
    LINT_VERSION=$(golangci-lint --version | head -n1)
    echo "✅ golangci-lint: $LINT_VERSION"
else
    echo "⚠️  golangci-lint not found. Installing..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    echo "✅ golangci-lint installed"
fi

echo ""
echo "📦 Installing dependencies..."
go mod download
go mod tidy

echo ""
echo "🔨 Building server..."
mkdir -p bin
go build -o bin/mcp-server ./cmd/server

echo ""
echo "✅ Build complete!"
echo ""
echo "📖 Quick Start Options:"
echo ""
echo "1. Run in stdio mode (for MCP integration):"
echo "   ./bin/mcp-server"
echo ""
echo "2. Run in HTTP mode (for remote access):"
echo "   ./bin/mcp-server --mode http --port 8080"
echo ""
echo "3. Analyze a Go file:"
echo "   ./bin/mcp-server analyze --file ./examples/sample.go"
echo ""
echo "4. Run tests:"
echo "   go test ./..."
echo ""
echo "5. Use Docker:"
echo "   docker-compose up -d"
echo ""
echo "📚 For more information, see README.md and docs/"
echo ""
echo "🎉 Setup complete! Happy coding!"
