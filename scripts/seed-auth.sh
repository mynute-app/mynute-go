#!/bin/bash

# Seed Auth Service
# This script seeds the auth service with endpoints and resources from the business service

AUTH_SERVICE_URL="${AUTH_SERVICE_URL:-http://localhost:4001}"
VERBOSE=false

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -u|--url) AUTH_SERVICE_URL="$2"; shift ;;
        -v|--verbose) VERBOSE=true ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

echo -e "\033[0;36m=== Auth Service Seeder ===\033[0m"
echo -e "\033[0;90mAuth Service URL: $AUTH_SERVICE_URL\033[0m"
echo ""

# Check if auth service is running
echo -e "\033[0;33mChecking auth service health...\033[0m"
if curl -s -f "$AUTH_SERVICE_URL/health" > /dev/null 2>&1; then
    HEALTH=$(curl -s "$AUTH_SERVICE_URL/health" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
    if [ "$HEALTH" = "healthy" ]; then
        echo -e "\033[0;32m✓ Auth service is healthy\033[0m"
    else
        echo -e "\033[0;33m⚠️  Auth service returned unexpected health status: $HEALTH\033[0m"
    fi
else
    echo -e "\033[0;31m✗ Auth service is not responding at $AUTH_SERVICE_URL\033[0m"
    echo -e "\033[0;33m  Please start the auth service first:\033[0m"
    echo -e "\033[0;90m  go run cmd/auth-service/main.go\033[0m"
    exit 1
fi

echo ""
echo -e "\033[0;33mStarting seeding process...\033[0m"

# Run the Go seeder
export AUTH_SERVICE_URL

if [ "$VERBOSE" = true ]; then
    go run cmd/seed-auth/main.go
else
    go run cmd/seed-auth/main.go 2>&1 | grep -E "^(✓|⚠️|===|[0-9]+/)|Seeding completed"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo -e "\033[0;32m=== Seeding Complete ===\033[0m"
    echo ""
    echo -e "\033[0;36mNext steps:\033[0m"
    echo -e "\033[0;90m1. Review policy definitions in core/src/config/db/seed/policy/\033[0m"
    echo -e "\033[0;90m2. Create policies via auth service admin panel\033[0m"
    echo -e "\033[0;90m3. Test authorization with: POST /authorize/by-method-and-path\033[0m"
else
    echo ""
    echo -e "\033[0;31m=== Seeding Failed ===\033[0m"
    echo -e "\033[0;33mCheck the error messages above for details\033[0m"
    exit 1
fi
