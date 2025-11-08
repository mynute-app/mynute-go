#!/bin/bash
# PostgreSQL healthcheck script for core service

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

pg_isready -U "${POSTGRES_USER}" -d "${DB_NAME}"
