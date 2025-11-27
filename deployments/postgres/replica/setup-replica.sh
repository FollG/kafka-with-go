#!/bin/sh

set -e

echo "Setting up PostgreSQL replica..."

until pg_isready -h postgres-master -p 5432 -U admin; do
  echo "Waiting for master to be ready..."
  sleep 2
done

pg_ctl -D "$PGDATA" -m fast -w stop || true

rm -rf "$PGDATA"/*

echo "Creating base backup from master..."
PGPASSWORD=password pg_basebackup -h postgres-master -U admin -D "$PGDATA" -P -R -X stream

echo "Setting up replication configuration..."
echo "Starting PostgreSQL replica..."
pg_ctl -D "$PGDATA" -w start

echo "PostgreSQL replica setup completed successfully!"