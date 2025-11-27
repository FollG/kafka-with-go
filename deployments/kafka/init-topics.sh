
#!/bin/bash

echo "Starting Kafka topics initialization..."
echo "Waiting for Kafka to be ready..."
until /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka1:9092 --list > /dev/null 2>&1; do
  echo "Waiting for Kafka to be ready..."
  sleep 5
done

echo "Kafka is ready. Creating topics..."

/opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka1:9092 \
  --create \
  --if-not-exists \
  --topic products \
  --partitions 3 \
  --replication-factor 3 \
  --config retention.ms=86400000 \
  --config cleanup.policy=delete

echo "Topic 'products' created successfully"

echo "Current topics:"
/opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka1:9092 --list

echo "Kafka topics initialization completed!"