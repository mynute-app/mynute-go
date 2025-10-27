# Production Migration Runner (PowerShell)
# Usage: .\scripts\migrate.ps1 [up|down|version|force] [-Steps 1] [-Version 123]

param(
    [Parameter(Position=0)]
    [ValidateSet("up", "down", "version", "force")]
    [string]$Action = "up",
    
    [Parameter()]
    [int]$Steps = 1,
    
    [Parameter()]
    [int]$Version = -1
)

# Load .env file if it exists
$envFile = ".env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

# Check if APP_ENV is set
$appEnv = $env:APP_ENV
if (-not $appEnv) {
    Write-Error "Error: APP_ENV is not set"
    exit 1
}

# Validate environment
if ($appEnv -notin @("prod", "dev", "test")) {
    Write-Error "Error: APP_ENV must be one of: prod, dev, test"
    exit 1
}

Write-Host "================================" -ForegroundColor Cyan
Write-Host "Database Migration Runner" -ForegroundColor Cyan
Write-Host "Environment: $appEnv" -ForegroundColor Cyan
Write-Host "Action: $Action" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Confirmation for production
if ($appEnv -eq "prod") {
    Write-Host "⚠️  WARNING: Running migrations in PRODUCTION environment!" -ForegroundColor Yellow
    Write-Host "Database: $env:POSTGRES_DB" -ForegroundColor Yellow
    Write-Host ""
    $confirmation = Read-Host "Are you sure you want to continue? (yes/no)"
    if ($confirmation -ne "yes") {
        Write-Host "Migration cancelled." -ForegroundColor Red
        exit 0
    }
    Write-Host ""
}

# Run migration based on action
try {
    switch ($Action) {
        "up" {
            Write-Host "Running migrations..." -ForegroundColor Green
            go run migrate/main.go -action=up -path=./migrations
        }
        "down" {
            Write-Host "Rolling back $Steps migration(s)..." -ForegroundColor Yellow
            go run migrate/main.go -action=down -steps=$Steps -path=./migrations
        }
        "version" {
            Write-Host "Checking migration version..." -ForegroundColor Cyan
            go run migrate/main.go -action=version -path=./migrations
        }
        "force" {
            if ($Version -lt 0) {
                Write-Error "Error: -Version is required for force action"
                exit 1
            }
            Write-Host "Forcing migration to version $Version..." -ForegroundColor Magenta
            go run migrate/main.go -action=force -version=$Version -path=./migrations
        }
    }
    
    Write-Host ""
    Write-Host "✅ Migration completed successfully!" -ForegroundColor Green
}
catch {
    Write-Error "Migration failed: $_"
    exit 1
}
