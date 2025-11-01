# Quick start script for Go Standards MCP Server (Windows)
# Run this script in PowerShell

Write-Host "🚀 Go Standards MCP Server - Quick Start" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check Go installation
$goVersion = go version 2>$null
if (-not $goVersion) {
    Write-Host "❌ Go is not installed. Please install Go 1.21 or higher." -ForegroundColor Red
    exit 1
}

Write-Host "✅ $goVersion" -ForegroundColor Green

# Check golangci-lint (optional)
$lintVersion = golangci-lint --version 2>$null
if ($lintVersion) {
    Write-Host "✅ golangci-lint: $($lintVersion[0])" -ForegroundColor Green
} else {
    Write-Host "⚠️  golangci-lint not found. You may want to install it manually." -ForegroundColor Yellow
    Write-Host "   Visit: https://golangci-lint.run/usage/install/" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📦 Installing dependencies..." -ForegroundColor Cyan
go mod download
go mod tidy

Write-Host ""
Write-Host "🔨 Building server..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path "bin" | Out-Null
go build -o bin\mcp-server.exe .\cmd\server

Write-Host ""
Write-Host "✅ Build complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📖 Quick Start Options:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. Run in stdio mode (for MCP integration):"
Write-Host "   .\bin\mcp-server.exe"
Write-Host ""
Write-Host "2. Run in HTTP mode (for remote access):"
Write-Host "   .\bin\mcp-server.exe --mode http --port 8080"
Write-Host ""
Write-Host "3. Run tests:"
Write-Host "   go test ./..."
Write-Host ""
Write-Host "4. Use Docker:"
Write-Host "   docker-compose up -d"
Write-Host ""
Write-Host "📚 For more information, see README.md and docs/" -ForegroundColor Cyan
Write-Host ""
Write-Host "🎉 Setup complete! Happy coding!" -ForegroundColor Green
