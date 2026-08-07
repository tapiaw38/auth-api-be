#!/bin/sh
set -e


echo "Starting Go application..."
exec go run ./cmd/api/ --host 0.0.0.0 --port 8082