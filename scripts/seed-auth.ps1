# Seed Auth Service
# This script seeds the auth service with endpoints and resources from the business service

param(
    [string]$AuthServiceUrl = "http://localhost:4001",
    [switch]$Verbose
)

Write-Host "=== Auth Service Seeder ===" -ForegroundColor Cyan
Write-Host "Auth Service URL: $AuthServiceUrl" -ForegroundColor Gray
Write-Host ""

# Check if auth service is running
Write-Host "Checking auth service health..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$AuthServiceUrl/health" -Method Get -TimeoutSec 5
    if ($health.status -eq "healthy") {
        Write-Host "✓ Auth service is healthy" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Auth service returned unexpected health status: $($health.status)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ Auth service is not responding at $AuthServiceUrl" -ForegroundColor Red
    Write-Host "  Please start the auth service first:" -ForegroundColor Yellow
    Write-Host "  go run cmd/auth-service/main.go" -ForegroundColor Gray
    exit 1
}

Write-Host ""
Write-Host "Starting seeding process..." -ForegroundColor Yellow

# Run the Go seeder
$env:AUTH_SERVICE_URL = $AuthServiceUrl

if ($Verbose) {
    go run cmd/seed-auth/main.go
} else {
    go run cmd/seed-auth/main.go 2>&1 | Where-Object { 
        $_ -match "^(✓|⚠️|===|\d+/)" -or $_ -match "Seeding completed" 
    }
}

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "=== Seeding Complete ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host "1. Review policy definitions in core/src/config/db/seed/policy/" -ForegroundColor Gray
    Write-Host "2. Create policies via auth service admin panel" -ForegroundColor Gray
    Write-Host "3. Test authorization with: POST /authorize/by-method-and-path" -ForegroundColor Gray
} else {
    Write-Host ""
    Write-Host "=== Seeding Failed ===" -ForegroundColor Red
    Write-Host "Check the error messages above for details" -ForegroundColor Yellow
    exit 1
}
