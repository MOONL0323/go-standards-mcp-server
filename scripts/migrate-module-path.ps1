#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Migrate Go module path safely with UTF-8 encoding preservation

.DESCRIPTION
    This script migrates the Go module path from GitHub-specific path to a generic one,
    while preserving UTF-8 encoding in all Go files.

.PARAMETER NewModulePath
    The new module path (e.g., "go-standards-mcp-server")

.PARAMETER OldModulePath
    The old module path (default: "github.com/MOONL0323/go-standards-mcp-server")

.PARAMETER Backup
    Create a backup branch before migration (default: true)

.EXAMPLE
    .\migrate-module-path.ps1 -NewModulePath "go-standards-mcp-server"

.EXAMPLE
    .\migrate-module-path.ps1 -NewModulePath "gitlab.com/myorg/go-standards" -OldModulePath "github.com/MOONL0323/go-standards-mcp-server"
#>

param(
    [Parameter(Mandatory=$true)]
    [string]$NewModulePath,
    
    [Parameter(Mandatory=$false)]
    [string]$OldModulePath = "github.com/MOONL0323/go-standards-mcp-server",
    
    [Parameter(Mandatory=$false)]
    [bool]$Backup = $true
)

$ErrorActionPreference = "Stop"

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Go Module Path Migration Script" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Old module path: $OldModulePath" -ForegroundColor Yellow
Write-Host "New module path: $NewModulePath" -ForegroundColor Green
Write-Host ""

# Check if git is available
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Git is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

# Check if go is available
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

# Check if we're in a git repository
$gitStatus = git rev-parse --is-inside-work-tree 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Not in a git repository" -ForegroundColor Red
    exit 1
}

# Check for uncommitted changes
$gitDiff = git status --porcelain
if ($gitDiff) {
    Write-Host "WARNING: You have uncommitted changes" -ForegroundColor Yellow
    Write-Host "It's recommended to commit or stash them before migration" -ForegroundColor Yellow
    $response = Read-Host "Continue anyway? (y/N)"
    if ($response -ne 'y' -and $response -ne 'Y') {
        Write-Host "Migration cancelled" -ForegroundColor Yellow
        exit 0
    }
}

# Create backup branch if requested
if ($Backup) {
    $backupBranch = "backup-before-migration-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
    Write-Host "Creating backup branch: $backupBranch" -ForegroundColor Cyan
    git checkout -b $backupBranch
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Failed to create backup branch" -ForegroundColor Red
        exit 1
    }
    
    # Return to original branch
    git checkout -
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Failed to return to original branch" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Backup branch created successfully" -ForegroundColor Green
    Write-Host ""
}

# Step 1: Update go.mod
Write-Host "Step 1: Updating go.mod..." -ForegroundColor Cyan
$goModPath = "go.mod"
if (Test-Path $goModPath) {
    $content = [System.IO.File]::ReadAllText($goModPath, [System.Text.Encoding]::UTF8)
    $content = $content -replace [regex]::Escape($OldModulePath), $NewModulePath
    [System.IO.File]::WriteAllText($goModPath, $content, [System.Text.Encoding]::UTF8)
    Write-Host "✓ go.mod updated" -ForegroundColor Green
} else {
    Write-Host "ERROR: go.mod not found" -ForegroundColor Red
    exit 1
}

# Step 2: Update all .go files
Write-Host "Step 2: Updating all .go files..." -ForegroundColor Cyan
$goFiles = Get-ChildItem -Recurse -Include *.go -File

$updatedCount = 0
$totalCount = $goFiles.Count

foreach ($file in $goFiles) {
    try {
        $content = [System.IO.File]::ReadAllText($file.FullName, [System.Text.Encoding]::UTF8)
        $originalContent = $content
        $content = $content -replace [regex]::Escape($OldModulePath), $NewModulePath
        
        if ($content -ne $originalContent) {
            [System.IO.File]::WriteAllText($file.FullName, $content, [System.Text.Encoding]::UTF8)
            $updatedCount++
        }
    } catch {
        Write-Host "ERROR: Failed to process $($file.FullName): $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
}

Write-Host "✓ Updated $updatedCount out of $totalCount Go files" -ForegroundColor Green
Write-Host ""

# Step 3: Clean and rebuild
Write-Host "Step 3: Cleaning and rebuilding..." -ForegroundColor Cyan
Write-Host "Cleaning module cache..." -ForegroundColor Yellow
go clean -modcache
if ($LASTEXITCODE -ne 0) {
    Write-Host "WARNING: Failed to clean module cache (non-critical)" -ForegroundColor Yellow
}

Write-Host "Tidying dependencies..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: go mod tidy failed" -ForegroundColor Red
    exit 1
}

Write-Host "Building server..." -ForegroundColor Yellow
go build -o bin/mcp-server.exe ./cmd/server
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to build server" -ForegroundColor Red
    exit 1
}

Write-Host "Building CLI..." -ForegroundColor Yellow
go build -o bin/go-standards-cli.exe ./cmd/cli
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to build CLI" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Build successful" -ForegroundColor Green
Write-Host ""

# Step 4: Run tests
Write-Host "Step 4: Running tests..." -ForegroundColor Cyan
go test ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "WARNING: Some tests failed (review manually)" -ForegroundColor Yellow
} else {
    Write-Host "✓ All tests passed" -ForegroundColor Green
}
Write-Host ""

# Step 5: Verify changes
Write-Host "Step 5: Verifying changes..." -ForegroundColor Cyan
$remainingOldPaths = Select-String -Path *.go -Pattern ([regex]::Escape($OldModulePath)) -Recurse 2>$null
if ($remainingOldPaths) {
    Write-Host "WARNING: Some files still contain old module path:" -ForegroundColor Yellow
    $remainingOldPaths | ForEach-Object { Write-Host "  $($_.Path):$($_.LineNumber)" -ForegroundColor Yellow }
    Write-Host ""
} else {
    Write-Host "✓ No old module paths found" -ForegroundColor Green
}

# Summary
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Migration Summary" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Old module path: $OldModulePath" -ForegroundColor Yellow
Write-Host "New module path: $NewModulePath" -ForegroundColor Green
Write-Host "Files updated:   $updatedCount Go files" -ForegroundColor Green
if ($Backup) {
    Write-Host "Backup branch:   $backupBranch" -ForegroundColor Cyan
}
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Review changes: git diff" -ForegroundColor White
Write-Host "  2. Test thoroughly: go test ./..." -ForegroundColor White
Write-Host "  3. Commit changes: git add . && git commit -m 'chore: migrate module path'" -ForegroundColor White
Write-Host "  4. Push to repository: git push origin main" -ForegroundColor White
Write-Host ""
Write-Host "Migration completed successfully!" -ForegroundColor Green
