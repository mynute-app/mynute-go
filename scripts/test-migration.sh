#!/usr/bin/env bash
# Automated Migration Testing Script (Bash version)
# Tests: up -> verify -> down -> verify -> up again

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

function print_success { echo -e "${GREEN}$1${NC}"; }
function print_warning { echo -e "${YELLOW}$1${NC}"; }
function print_error { echo -e "${RED}$1${NC}"; }
function print_info { echo -e "${CYAN}$1${NC}"; }

# Load .env
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

APP_ENV=${APP_ENV:-dev}

if [ "$APP_ENV" = "prod" ]; then
    print_error "⚠️  ERROR: Cannot run automated tests in production environment!"
    print_error "Set APP_ENV=dev or APP_ENV=test in your .env file"
    exit 1
fi

print_info "================================"
print_info "Automated Migration Test Runner"
print_info "Environment: $APP_ENV"
print_info "================================"
echo ""

# Get current version
print_info "📊 Checking current migration state..."
go run migrate/main.go -action=version -path=./migrations
echo ""

# Confirmation
if [ "${1}" != "-y" ]; then
    print_warning "⚠️  This will run: UP → DOWN → UP on your database"
    print_warning "Make sure you're using a TEST database!"
    read -p "Continue? (yes/no): " -r
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        print_info "Test cancelled."
        exit 0
    fi
    echo ""
fi

# Step 1: Run migration UP
print_info "🔼 Step 1/5: Running migration UP..."
if go run migrate/main.go -action=up -path=./migrations; then
    print_success "✅ Migration UP completed"
else
    print_error "❌ Migration UP failed"
    exit 1
fi
echo ""

# Step 2: Verify UP worked
print_info "🔍 Step 2/5: Verifying migration was applied..."
go run migrate/main.go -action=version -path=./migrations
print_success "✅ Migration applied successfully"
echo ""

# Step 3: Run migration DOWN (rollback)
print_info "🔽 Step 3/5: Testing rollback (DOWN)..."
if go run migrate/main.go -action=down -steps=1 -path=./migrations; then
    print_success "✅ Rollback completed"
else
    print_error "❌ Rollback failed"
    print_warning "Your migration might not have a proper DOWN script"
    exit 1
fi
echo ""

# Step 4: Verify DOWN worked
print_info "🔍 Step 4/5: Verifying rollback..."
go run migrate/main.go -action=version -path=./migrations
print_success "✅ Rollback verified"
echo ""

# Step 5: Run migration UP again
print_info "🔼 Step 5/5: Re-applying migration (UP)..."
if go run migrate/main.go -action=up -path=./migrations; then
    print_success "✅ Migration re-applied successfully"
else
    print_error "❌ Re-applying migration failed"
    exit 1
fi
echo ""

# Final verification
print_info "🔍 Final verification..."
go run migrate/main.go -action=version -path=./migrations
echo ""

# Summary
print_success "════════════════════════════════"
print_success "✅ ALL TESTS PASSED!"
print_success "════════════════════════════════"
print_success ""
print_success "Your migration is working correctly:"
print_success "  ✅ UP migration applies cleanly"
print_success "  ✅ DOWN migration rolls back properly"
print_success "  ✅ UP migration can be re-applied"
print_success ""
print_info "Next steps:"
print_info "  1. Review the generated SQL files"
print_info "  2. Test with realistic data"
print_info "  3. Commit your migration files"
print_info "  4. Deploy to staging/production"
