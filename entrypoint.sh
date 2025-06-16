#!/bin/sh

set -e

sleep 5

echo "Running PostgreSQL migrations..."
./migrator -mode up \
  -storage-path "postgres://evans:evans@postgres:8081/postgres?sslmode=disable" \
  -migrations-path "./migrations/postgres"

echo "Running ClickHouse migrations..."
./migrator -mode up \
  -storage-path "clickhouse://evans:evans@clickhouse:9000/logs" \
  -migrations-path "./migrations/clickhouse"

echo "All migrations completed successfully."

echo "Starting main application..."
./main -config ./example.env