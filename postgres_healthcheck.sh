#!/bin/sh

POSTGRES_HOST=$1

POSTGRES_PORT=${2:-$DATABASE_PORT}

if [ -z "$POSTGRES_HOST" ]; then
  echo "PostgreSQL host not provided"
  exit 1
fi

if [ -z "$POSTGRES_PORT" ]; then
  echo "PostgreSQL port not provided and not set in environment"
  exit 1
fi

# Ожидаем, пока PostgreSQL станет доступной
until nc -z -v -w30 "$POSTGRES_HOST" "$POSTGRES_PORT"
  do
    echo "Waiting for PostgreSQL to be available at $POSTGRES_HOST:$POSTGRES_PORT..."
    sleep 5
done

echo "PostgreSQL is up and running at $POSTGRES_HOST:$POSTGRES_PORT!"
