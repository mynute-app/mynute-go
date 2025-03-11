#!/bin/bash
set -e

echo "Initializing PostgreSQL database script..."

# Ensure the main database exists before running any commands
echo "Checking if main database '$POSTGRES_DB' exists..."
MAIN_DB_EXISTS=$(psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$POSTGRES_DB'")

if [ "$MAIN_DB_EXISTS" == "1" ]; then
  echo "✅   Main database '$POSTGRES_DB' already exists."
else
  echo "🚀   Creating main database: $POSTGRES_DB"
  psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE \"$POSTGRES_DB\";"
fi

# Check if APP_ENV is set to "test" and create additional test database
if [ "$APP_ENV" == "test" ]; then
  TEST_DB="${POSTGRES_DB}-${APP_ENV}"
  echo "Checking if test database '$TEST_DB' exists..."

  TEST_DB_EXISTS=$(psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$TEST_DB'")

  if [ "$TEST_DB_EXISTS" == "1" ]; then
    echo "✅   Test database '$TEST_DB' already exists. Skipping creation."
  else
    echo "🚀   Creating test database: $TEST_DB"
    psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE \"$TEST_DB\";"
  fi
else
  echo "APP_ENV is not 'test'. Skipping test database creation."
fi

echo "🎉   Database initialization script complete."