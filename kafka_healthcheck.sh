#!/bin/sh

# Параметры (передаем хост и порт через аргументы)
KAFKA_HOST=$1
KAFKA_PORT=$2

if [ -z "$KAFKA_HOST" ]; then
  echo "Kafka host not provided"
  exit 1
fi

if [ -z "$KAFKA_PORT" ]; then
  echo "Kafka port not provided"
  exit 1
fi

# Ожидаем, пока Kafka станет доступной
until nc -z -v -w30 $KAFKA_HOST $KAFKA_PORT
do
  echo "Waiting for Kafka to be available at $KAFKA_HOST:$KAFKA_PORT..."
  sleep 5
done

echo "Kafka is up and running at $KAFKA_HOST:$KAFKA_PORT!"

# Теперь можно запускать приложение
# exec "$@"
