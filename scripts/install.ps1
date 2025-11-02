# Go Standards MCP Server - Installation Script
# For Windows PowerShell

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Go Standards MCP Server - Installation" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check Go installation
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go not found" -ForegroundColor Red
    Write-Host "Please install Go 1.21+ from https://go.dev/dl/"
    exit 1
}

$goVersion = go version
Write-Host "Detected Go: $goVersion" -ForegroundColor Green
Write-Host ""

# Set project directory
$projectDir = "$HOME\go-standards-mcp-server"

# Check if already in project directory
if (Test-Path "go.mod") {
    $content = Get-Content "go.mod"
    if ($content -match "go-standards-mcp-server") {
        Write-Host "Already in project directory" -ForegroundColor Green
        $projectDir = Get-Location
    }
} else {
    # Clone or update project
    if (Test-Path $projectDir) {
        Write-Host "Directory exists, updating..." -ForegroundColor Yellow
        Set-Location $projectDir
        git pull
    } else {
        Write-Host "Cloning repository..." -ForegroundColor Cyan
        Set-Location $HOME
        git clone https://github.com/MOONL0323/go-standards-mcp-server.git
        Set-Location go-standards-mcp-server
    }
    $projectDir = Get-Location
}

Write-Host ""
Write-Host "Installing dependencies..." -ForegroundColor Cyan
go mod download

Write-Host ""
Write-Host "Building server..." -ForegroundColor Cyan
if (Test-Path "Makefile") {
    make build
} else {
    go build -o bin\mcp-server.exe .\cmd\server
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Installation Complete" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next Step: Configure Cursor IDE"
Write-Host ""
Write-Host "1. Open Cursor Settings (Gear Icon > Settings)"
Write-Host "2. Search for 'mcp'"
Write-Host "3. Add the following configuration:"
Write-Host ""

$configPath = "$projectDir\bin\mcp-server.exe" -replace '\\', '\\'
$config = @"
{
  "mcpServers": {
    "go-standards": {
      "command": "$configPath"
    }
  }
}
"@

Write-Host $config -ForegroundColor Yellow
Write-Host ""
Write-Host "4. Save and restart Cursor"
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Server Path: $projectDir\bin\mcp-server.exe" -ForegroundColor Green
Write-Host "Test Command: @go-standards health_check"
Write-Host ""
