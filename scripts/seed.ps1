#!/usr/bin/env pwsh
# seed.ps1 - Database seeding script for Windows
# Usage: .\scripts\seed.ps1

param(
    [switch]$Build,
    [switch]$Help
)

$ErrorActionPreference = "Stop"

function Show-Help {
    Write-Host ""
    Write-Host "Database Seeding Script" -ForegroundColor Cyan
    Write-Host "======================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage:" -ForegroundColor Yellow
    Write-Host "  .\scripts\seed.ps1           - Run seeding directly"
    Write-Host "  .\scripts\seed.ps1 -Build    - Build seed binary"
    Write-Host "  .\scripts\seed.ps1 -Help     - Show this help"
    Write-Host ""
    Write-Host "What gets seeded:" -ForegroundColor Yellow
    Write-Host "  - System Resources"
    Write-Host "  - System Roles"
    Write-Host "  - API Endpoints"
    Write-Host "  - Access Policies"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  Development:  .\scripts\seed.ps1"
    Write-Host "  Production:   .\scripts\seed.ps1 -Build"
    Write-Host "               (then copy bin\seed.exe to server)"
    Write-Host ""
}

if ($Help) {
    Show-Help
    exit 0
}

if ($Build) {
    Write-Host "Building seed binary..." -ForegroundColor Cyan
    go build -o bin/seed.exe cmd/seed/main.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Binary created at: bin\seed.exe" -ForegroundColor Green
        Write-Host ""
        Write-Host "To run in production:" -ForegroundColor Yellow
        Write-Host "  1. Copy bin\seed.exe to your production server"
        Write-Host "  2. Set APP_ENV and database environment variables"
        Write-Host "  3. Run: .\seed.exe"
    } else {
        Write-Host "✗ Build failed" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Running database seeding..." -ForegroundColor Cyan
    go run cmd/seed/main.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Seeding completed successfully" -ForegroundColor Green
    } else {
        Write-Host "✗ Seeding failed" -ForegroundColor Red
        exit 1
    }
}
