#!/bin/sh

echo "========== Starting Go application =========="
exec go run ./cmd/api/ --host 0.0.0.0