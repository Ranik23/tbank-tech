#!/bin/sh
set -e

export DB_DSN="postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=${DATABASE_SSL}"

echo "Using DB_DSN: $DB_DSN"

MIGRATION_DIR=${1:-./migrations}
echo "Running migrations from $MIGRATION_DIR..."
goose -dir "$MIGRATION_DIR" postgres "$DB_DSN" up