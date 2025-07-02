#!/bin/sh

echo "Waiting for PostgreSQL main..."
until pg_isready -h main -p 5432 -U user; do
  sleep 2
done

# echo "Waiting for Kafka..."
# until kafkacat -b kafka:9092 -L >/dev/null 2>&1; do
#   echo "Kafka is not ready yet..."
#   sleep 2
# done

echo "All services are up. Starting the app..."
exec "$@"
