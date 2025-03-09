#!/bin/bash
set -e

echo "Initializing PostgreSQL database script..."

# Check if APP_ENV is set to "test" and create additional database
if [ "$APP_ENV" == "test" ]; then
  TEST_DB="${POSTGRES_DB}-${APP_ENV}"
  echo "Creating test database: $TEST_DB"
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE DATABASE \"$TEST_DB\";"
else
  echo "APP_ENV is not 'test'. Skipping test database creation."
fi

echo "Database initialization script complete."