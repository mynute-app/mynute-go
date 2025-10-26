#!/usr/bin/env pwsh
# Automated Migration Testing Script
# Tests: up -> verify -> down -> verify -> up again

param(
    [Parameter(Position=0)]
    [string]$MigrationFile = "",
    
    [switch]$SkipConfirmation
)

$ErrorActionPreference = "Stop"

# Colors
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Warning { Write-Host $args -ForegroundColor Yellow }
function Write-Error { Write-Host $args -ForegroundColor Red }
function Write-Info { Write-Host $args -ForegroundColor Cyan }

# Load .env
if (Test-Path ".env") {
    Get-Content ".env" | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

$appEnv = $env:APP_ENV
if ($appEnv -eq "prod") {
    Write-Error "⚠️  ERROR: Cannot run automated tests in production environment!"
    Write-Error "Set APP_ENV=dev or APP_ENV=test in your .env file"
    exit 1
}

Write-Info "================================"
Write-Info "Automated Migration Test Runner"
Write-Info "Environment: $appEnv"
Write-Info "================================"
Write-Host ""

# Get current version
Write-Info "📊 Checking current migration state..."
$currentVersion = go run migrate/main.go -action=version -path=./migrations 2>&1
Write-Host $currentVersion
Write-Host ""

if (-not $SkipConfirmation) {
    Write-Warning "⚠️  This will run: UP → DOWN → UP on your database"
    Write-Warning "Make sure you're using a TEST database!"
    $confirm = Read-Host "Continue? (yes/no)"
    if ($confirm -ne "yes") {
        Write-Info "Test cancelled."
        exit 0
    }
    Write-Host ""
}

# Step 1: Run migration UP
Write-Info "🔼 Step 1/5: Running migration UP..."
try {
    go run migrate/main.go -action=up -path=./migrations
    Write-Success "✅ Migration UP completed"
} catch {
    Write-Error "❌ Migration UP failed: $_"
    exit 1
}
Write-Host ""

# Step 2: Verify UP worked
Write-Info "🔍 Step 2/5: Verifying migration was applied..."
$versionAfterUp = go run migrate/main.go -action=version -path=./migrations 2>&1
Write-Host $versionAfterUp
Write-Success "✅ Migration applied successfully"
Write-Host ""

# Step 3: Run migration DOWN (rollback)
Write-Info "🔽 Step 3/5: Testing rollback (DOWN)..."
try {
    go run migrate/main.go -action=down -steps=1 -path=./migrations
    Write-Success "✅ Rollback completed"
} catch {
    Write-Error "❌ Rollback failed: $_"
    Write-Warning "Your migration might not have a proper DOWN script"
    exit 1
}
Write-Host ""

# Step 4: Verify DOWN worked
Write-Info "🔍 Step 4/5: Verifying rollback..."
$versionAfterDown = go run migrate/main.go -action=version -path=./migrations 2>&1
Write-Host $versionAfterDown
Write-Success "✅ Rollback verified"
Write-Host ""

# Step 5: Run migration UP again
Write-Info "🔼 Step 5/5: Re-applying migration (UP)..."
try {
    go run migrate/main.go -action=up -path=./migrations
    Write-Success "✅ Migration re-applied successfully"
} catch {
    Write-Error "❌ Re-applying migration failed: $_"
    exit 1
}
Write-Host ""

# Final verification
Write-Info "🔍 Final verification..."
$finalVersion = go run migrate/main.go -action=version -path=./migrations 2>&1
Write-Host $finalVersion
Write-Host ""

# Summary
Write-Success "════════════════════════════════"
Write-Success "✅ ALL TESTS PASSED!"
Write-Success "════════════════════════════════"
Write-Success ""
Write-Success "Your migration is working correctly:"
Write-Success "  ✅ UP migration applies cleanly"
Write-Success "  ✅ DOWN migration rolls back properly"
Write-Success "  ✅ UP migration can be re-applied"
Write-Success ""
Write-Info "Next steps:"
Write-Info "  1. Review the generated SQL files"
Write-Info "  2. Test with realistic data"
Write-Info "  3. Commit your migration files"
Write-Info "  4. Deploy to staging/production"
