#!/bin/sh

echo "Waiting for Postgres at ${PG_DB_HOST:-auth-postgres-db}:${PG_DB_PORT:-5432}..."
until nc -z "${PG_DB_HOST:-auth-postgres-db}" "${PG_DB_PORT:-5432}"; do
    sleep 1
done
echo "Postgres is ready."

echo "Waiting for RabbitMQ at ${RABBIT_HOST:-rabbitmq}:${RABBIT_PORT:-5672}..."
until nc -z "${RABBIT_HOST:-rabbitmq}" "${RABBIT_PORT:-5672}"; do
    sleep 1
done
echo "RabbitMQ is ready."

echo "Starting Go application..."
exec go run ./cmd/api/ --host 0.0.0.0 --port 8082