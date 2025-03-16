#!/bin/sh
set -e

MIGRATION_DIR=${1:-./migrations}

echo "Running migrations from $MIGRATION_DIR..."
goose -dir "$MIGRATION_DIR" postgres "$DB_DSN" up

echo "Starting application..."
exec ./merch-store
