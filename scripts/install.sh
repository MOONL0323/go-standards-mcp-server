#!/bin/bash

# Go Standards MCP Server - Installation Script
# For macOS and Linux

set -e

echo "=========================================="
echo "Go Standards MCP Server - Installation"
echo "=========================================="
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "ERROR: Go not found"
    echo "Please install Go 1.21+ from https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "Detected Go: $GO_VERSION"
echo ""

# Check if already in project directory
if [ -f "go.mod" ] && grep -q "go-standards-mcp-server" go.mod; then
    echo "Already in project directory"
    PROJECT_DIR=$(pwd)
else
    # Clone or update project
    if [ -d "$HOME/go-standards-mcp-server" ]; then
        echo "Directory exists, updating..."
        cd "$HOME/go-standards-mcp-server"
        git pull
    else
        echo "Cloning repository..."
        cd "$HOME"
        git clone https://github.com/MOONL0323/go-standards-mcp-server.git
        cd go-standards-mcp-server
    fi
    PROJECT_DIR=$(pwd)
fi

echo ""
echo "Installing dependencies..."
go mod download

echo ""
echo "Building server..."
make build || go build -o bin/mcp-server ./cmd/server

echo ""
echo "=========================================="
echo "Installation Complete"
echo "=========================================="
echo ""
echo "Next Step: Configure Cursor IDE"
echo ""
echo "1. Open Cursor Settings (Gear Icon > Settings)"
echo "2. Search for 'mcp'"
echo "3. Add the following configuration:"
echo ""
echo '{
  "mcpServers": {
    "go-standards": {
      "command": "'$PROJECT_DIR'/bin/mcp-server"
    }
  }
}'
echo ""
echo "4. Save and restart Cursor"
echo ""
echo "=========================================="
echo ""
echo "Server Path: $PROJECT_DIR/bin/mcp-server"
echo "Test Command: @go-standards health_check"
echo ""
