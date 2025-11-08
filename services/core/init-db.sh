#!/bin/bash
set -e

# Trap SIGTERM and SIGINT for graceful shutdown
trap 'echo "🛑 Received shutdown signal, stopping PostgreSQL gracefully..."; pg_ctl stop -D "$PGDATA" -m fast; exit 0' SIGTERM SIGINT

echo '✅ Starting PostgreSQL...'
docker-entrypoint.sh postgres &
PG_PID=$!

echo '⏳ Waiting for PostgreSQL to start...'
until pg_isready -h localhost -p 5432 -U "${POSTGRES_USER}"; do
  echo '⏳ Still waiting for PostgreSQL...'
  sleep 2
done

echo '✅ PostgreSQL is ready! Proceeding with initialization...'

DB_NAME=""
case ${APP_ENV} in
  prod)
    DB_NAME=${POSTGRES_DB_PROD}
    ;;
  test)
    DB_NAME=${POSTGRES_DB_TEST}
    ;;
  dev)
    DB_NAME=${POSTGRES_DB_DEV}
    ;;
  *)
    echo '❌ APP_ENV must be one of prod, test, or dev'
    exit 1
    ;;
esac

# Ensure DB_NAME is not empty
if [ -z "${DB_NAME}" ]; then
  echo '❌ ERROR: DB_NAME is not set. Please check your .env file.'
  echo "❌ Required: POSTGRES_DB_PROD, POSTGRES_DB_TEST, or POSTGRES_DB_DEV based on APP_ENV=${APP_ENV}"
  exit 1
fi

# Ensure the database exists
echo "Checking if database '${DB_NAME}' exists..."
DB_EXISTS=$(psql -U "${POSTGRES_USER}" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" 2>/dev/null || echo "0")
if echo "${DB_EXISTS}" | grep -q "1"; then
  echo "✅ Database '${DB_NAME}' already exists."
else
  echo "🚀 Creating database: '${DB_NAME}'"
  psql -U "${POSTGRES_USER}" -d postgres -c "CREATE DATABASE ${DB_NAME};"
fi

echo '🎉 Database initialization complete!'

# Keep PostgreSQL running in the foreground and wait for it
wait $PG_PID
