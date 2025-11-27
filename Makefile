.PHONY: build build-api build-processor test migrate-up migrate-down docker-up docker-down logs backup healthcheck kafka-topics

build: build-api build-processor

build-api:
	@echo "Building API service..."
	docker build -f deployments/Dockerfile.api -t product-api:latest .

build-processor:
	@echo "Building Processor service..."
	docker build -f deployments/Dockerfile.processor -t product-processor:latest .

test:
	@echo "Running tests..."
	go test -v ./...

migrate-up:
	@echo "Running database migrations..."
	docker exec -i postgres-master psql -U admin -d products -c "CREATE TYPE product_type AS ENUM ('clothing_headwear', 'clothing_body', 'clothing_pants', 'clothing_shoes', 'food', 'furniture', 'electronics', 'adult', 'home_goods');" || true
	docker exec -i postgres-master psql -U admin -d products -f /docker-entrypoint-initdb.d/init.sql

migrate-down:
	@echo "Reverting database migrations..."
	docker exec -i postgres-master psql -U admin -d products -c "DROP TABLE IF EXISTS products; DROP TYPE IF EXISTS product_type;"

docker-up:
	@echo "Starting all services..."
	cd deployments && docker-compose up -d

docker-down:
	@echo "Stopping all services..."
	cd deployments && docker-compose down

docker-clean:
	@echo "Stopping and removing all services with volumes..."
	cd deployments && docker-compose down -v

logs:
	@echo "Showing logs..."
	cd deployments && docker-compose logs -f

logs-api:
	cd deployments && docker-compose logs -f api

logs-processor:
	cd deployments && docker-compose logs -f processor

logs-kafka:
	cd deployments && docker-compose logs -f kafka1 kafka2 kafka3

backup:
	@echo "Creating database backup..."
	./scripts/backup.sh

healthcheck:
	@echo "Running health checks..."
	curl -f http://localhost:8080/health || exit 1

kafka-topics:
	@echo "Listing Kafka topics..."
	docker exec kafka1 kafka-topics.sh --list --bootstrap-server kafka1:9092

kafka-describe:
	@echo "Describing products topic..."
	docker exec kafka1 kafka-topics.sh --describe --topic products --bootstrap-server kafka1:9092

psql-master:
	docker exec -it postgres-master psql -U admin -d products

psql-replica:
	docker exec -it postgres-replica psql -U admin -d products

redis-cli:
	docker exec -it redis redis-cli

metrics:
	@echo "API Metrics: curl http://localhost:9091/metrics"
	@echo "Processor Metrics: curl http://localhost:9092/metrics"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin)"

prometheus:
	@echo "Opening Prometheus..."
	open http://localhost:9090

grafana:
	@echo "Opening Grafana..."
	open http://localhost:3000

monitor-metrics: metrics
	@echo "Use 'make prometheus' to open Prometheus"
	@echo "Use 'make grafana' to open Grafana"

monitor:
	@echo "Service URLs:"
	@echo "API: http://localhost:8080"
	@echo "API Health: http://localhost:8080/health"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "API Metrics: http://localhost:9091/metrics"
	@echo "Processor Metrics: http://localhost:9092/metrics"

init-dirs:
	@echo "Creating necessary directories..."
	mkdir -p deployments/postgres/master
	mkdir -p deployments/postgres/replica
	mkdir -p deployments/kafka
	mkdir -p deployments/monitoring/grafana/provisioning/datasources
	mkdir -p deployments/monitoring/grafana/provisioning/dashboards

setup-permissions:
	@echo "Setting up script permissions..."
	chmod +x deployments/postgres/replica/setup-replica.sh


.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build Docker images"
	@echo "  build-api     - Build only API service"
	@echo "  build-processor - Build only Processor service"
	@echo "  test          - Run tests"
	@echo "  docker-up     - Start all services"
	@echo "  docker-down   - Stop all services"
	@echo "  docker-clean  - Stop and remove all services with volumes"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Revert database migrations"
	@echo "  logs          - Show all logs"
	@echo "  logs-api      - Show API logs"
	@echo "  logs-processor - Show Processor logs"
	@echo "  backup        - Create database backup"
	@echo "  healthcheck   - Run health checks"
	@echo "  kafka-topics  - List Kafka topics"
	@echo "  kafka-describe - Describe products topic"
	@echo "  psql-master   - Connect to master PostgreSQL"
	@echo "  psql-replica  - Connect to replica PostgreSQL"
	@echo "  redis-cli     - Connect to Redis"
	@echo "  monitor       - Show service URLs"
	@echo "  init-dirs     - Create necessary directories"
	@echo "  setup-permissions - Set up script permissions"