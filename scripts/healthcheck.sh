#!/bin/bash

# Health check для всех сервисов

# Проверяем API
curl -f http://localhost:8080/health || exit 1

# Проверяем PostgreSQL
pg_isready -h localhost -p 5432 || exit 1

# Проверяем Redis
redis-cli ping || exit 1

# Проверяем Kafka
kafka-topics.sh --list --bootstrap-server localhost:9092 > /dev/null 2>&1 || exit 1

echo "All services are healthy"