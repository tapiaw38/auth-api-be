#!/bin/sh

until pg_isready --host=cardon-postgres-db --port=5432 --username=postgres --dbname=auth-api-db
do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "========== Starting Go application =========="
exec go run ./cmd/api/ --host 0.0.0.0