#!/usr/bin/env bash
# Production Migration Runner
# Usage: ./scripts/migrate.sh [up|down|version|force VERSION]

set -e

# Load environment variables from .env if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if APP_ENV is set
if [ -z "$APP_ENV" ]; then
    echo "Error: APP_ENV is not set"
    exit 1
fi

# Validate environment
if [ "$APP_ENV" != "prod" ] && [ "$APP_ENV" != "dev" ] && [ "$APP_ENV" != "test" ]; then
    echo "Error: APP_ENV must be one of: prod, dev, test"
    exit 1
fi

# Check for action
ACTION=${1:-up}

echo "================================"
echo "Database Migration Runner"
echo "Environment: $APP_ENV"
echo "Action: $ACTION"
echo "================================"
echo ""

# Confirmation for production
if [ "$APP_ENV" = "prod" ]; then
    echo "⚠️  WARNING: Running migrations in PRODUCTION environment!"
    echo "Database: $POSTGRES_DB"
    echo ""
    read -p "Are you sure you want to continue? (yes/no): " -r
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        echo "Migration cancelled."
        exit 0
    fi
    echo ""
fi

# Run migration based on action
case $ACTION in
    up)
        echo "Running migrations..."
        go run migrate/main.go -action=up -path=./migrations
        ;;
    down)
        STEPS=${2:-1}
        echo "Rolling back $STEPS migration(s)..."
        go run migrate/main.go -action=down -steps=$STEPS -path=./migrations
        ;;
    version)
        echo "Checking migration version..."
        go run migrate/main.go -action=version -path=./migrations
        ;;
    force)
        VERSION=$2
        if [ -z "$VERSION" ]; then
            echo "Error: VERSION is required for force action"
            echo "Usage: $0 force VERSION"
            exit 1
        fi
        echo "Forcing migration to version $VERSION..."
        go run migrate/main.go -action=force -version=$VERSION -path=./migrations
        ;;
    *)
        echo "Error: Unknown action '$ACTION'"
        echo "Usage: $0 [up|down|version|force VERSION]"
        exit 1
        ;;
esac

echo ""
echo "✅ Migration completed successfully!"
