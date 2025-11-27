<div align="center">

# 🚀 Kafka with Go

**A practical guide and code examples for building robust event-driven applications with Apache Kafka and the Go programming language.**

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go)](https://golang.org)
[![Kafka](https://img.shields.io/badge/Apache%20Kafka-4%2B-231F20?logo=apachekafka)](https://kafka.apache.org/40/documentation.html)
[![Postgres](https://img.shields.io/badge/PSQL-16-blue)](https://www.postgresql.org/about/news/postgresql-16-released-2715/)
[![Prometheus](https://img.shields.io/badge/Promethus-last-red)](https://prometheus.io/docs/introduction/overview/)
[![Grafana](https://img.shields.io/badge/Grafana-last-green)](https://grafana.com/docs/)
[![Redis](https://img.shields.io/badge/Redis-7.2-violet)](https://redis.io/docs/latest/)
[![License](https://img.shields.io/badge/License-Apache2.0-gray.svg)](https://opensource.org/license/apache-2-0)


</div>

---

## 📖 Table of Contents

- [Overview](#overview)
- [What's Inside?](#whats-inside)
- [Quick Start](#quick-start)

---

## Overview

The project can be used as a platform for training in performing work tasks. The project has a deliberately poorly designed database subject area. The request path starts with an API written in Golang using SOLID and CLEAN ARCHITECTURE. Upon a successful POST request, the API returns a 202 HTTP status code, and the request goes to a Kafka topic, from where psql-master pulls it and replicates it to psql-replica. The application also uses a Redis cache.

![Go Gopher](https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png)

---

## What's Inside?

- **`producer/`**: Go application to produce messages to Kafka
- **`consumer/`**: Go application to consume messages from Kafka
- **`psql/`**: Main database
- **`psql replica/`**: Replication of main database
- **`redis/`**: Cache 
- **`promethus/`**: Metrics
- **`grafana/`**: Monitoring
- **`docker-compose.yml`**: Docker setup for local Kafka cluster

---

## Quick Start

### Prerequisites

- **Docker & Docker Compose** - [Get Docker](https://docs.docker.com/get-docker/)

### Running project with Docker

```bash
git clone https://github.com/FollG/kafka-with-go.git
cd kafka-with-go
make build
make docker-up
```
