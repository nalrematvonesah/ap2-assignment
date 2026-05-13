````md
# Payment Processing Platform

Production-style event-driven microservices system built with Go, gRPC, RabbitMQ, PostgreSQL, Redis, Docker, and Next.js.

---

# Architecture

```text
Frontend (Next.js)
        ↓ HTTP
Gateway API (Gin)
        ↓ gRPC
Payment Service
   ├── PostgreSQL
   ├── Redis Cache
   └── RabbitMQ Publisher
                ↓
          RabbitMQ Queue
                ↓
      Notification Worker
         ├── Retry Logic
         ├── DLQ
         ├── Redis Idempotency
         └── Email Provider Adapter
````

---

# Tech Stack

| Technology     | Purpose                        |
| -------------- | ------------------------------ |
| Go             | Backend services               |
| gRPC           | Internal service communication |
| RabbitMQ       | Event broker                   |
| PostgreSQL     | Persistent storage             |
| Redis          | Cache + idempotency            |
| Docker Compose | Infrastructure orchestration   |
| Next.js        | Frontend                       |
| Gin            | API Gateway                    |

---

# Features

## Distributed Architecture

* Event-driven communication
* Asynchronous processing
* Background workers
* Service isolation

---

## Reliability

* Retry strategy
* Exponential backoff
* Dead Letter Queue (DLQ)
* Manual ACK/NACK
* Persistent RabbitMQ messages
* Graceful shutdown

---

## Performance

* Redis cache-aside pattern
* Reduced database load
* Fast in-memory idempotency checks

---

## Scalability

* Independent services
* Horizontal worker scaling
* Queue-based workload distribution

---

# Services

## Frontend

Next.js UI for payment processing.

### Port

3000

---

## Gateway

HTTP API gateway using Gin.

### Responsibilities

* Receives HTTP requests
* Validates payloads
* Calls Payment Service via gRPC

### Port

8080

---

## Payment Service

Core payment processing service.

### Responsibilities

* Processes payments
* Stores data in PostgreSQL
* Caches payment data in Redis
* Publishes RabbitMQ events

### Port

50051

---

## Notification Worker

Background worker consuming RabbitMQ events.

### Responsibilities

* Consumes events asynchronously
* Sends notifications
* Retries failed jobs
* Uses Redis idempotency
* Sends failed messages to DLQ

---

# Infrastructure

## RabbitMQ

### Ports

| Port  | Purpose       |
| ----- | ------------- |
| 5672  | AMQP          |
| 15672 | Management UI |

### UI

[http://localhost:15672](http://localhost:15672)

Credentials:

guest
guest

---

## PostgreSQL

### Port

5432

---

## Redis

### Port

6379

---

# Project Structure

```text
.
├── frontend
├── gateway
├── payment-service
├── notification-service
├── docker-compose.yml
└── README.md
```

---

# Queue Topology

## Main Queue

payment.completed

Processes successful payment events.

---

## Dead Letter Queue

payment.completed.dlq

Stores failed events after retries.

---

# Event Flow

## Payment Processing

```text
1. Frontend sends HTTP request
2. Gateway receives request
3. Gateway calls Payment Service via gRPC
4. Payment Service stores payment in PostgreSQL
5. Payment Service caches payment in Redis
6. Payment Service publishes RabbitMQ event
7. Notification Worker consumes event
8. Notification Worker sends notification
9. Worker ACKs message
```

---

# Retry Flow

```text
Message received
      ↓
Processing failed
      ↓
Retry with exponential backoff
      ↓
Still failing
      ↓
NACK(false, false)
      ↓
DLQ
```

---

# Cache-Aside Pattern

```text
Request
   ↓
Redis lookup
   ↓ hit
Return cached data

OR

Redis miss
   ↓
PostgreSQL query
   ↓
Save to Redis
   ↓
Return data
```

---

# Idempotency

Notification Worker prevents duplicate processing using Redis.

### Key format

```text
event:<event_id>
```

---

# Running the Project

## Build and start

```bash
docker compose up --build
```

---

## Stop containers

```bash
docker compose down
```

---

## Remove volumes

```bash
docker compose down -v
```

---

# API

## Process Payment

### Endpoint

POST /payments

### Example Request

```json
{
  "order_id": 1,
  "amount": 250,
  "email": "user@test.com"
}
```

### Example Response

```json
{
  "status": "completed"
}
```

---

# Redis Inspection

Open Redis CLI:

```bash
docker exec -it ap2-assignment-redis-1 redis-cli
```

### List keys

```redis
KEYS *
```

---

# PostgreSQL Inspection

Open PostgreSQL shell:

```bash
docker exec -it ap2-assignment-postgres-1 psql -U postgres
```

### Open database

```sql
\c payments
```

### Show payments

```sql
SELECT * FROM payments;
```

---

# DLQ Testing

Temporary test:

```go
msg.Nack(false, false)
return
```

inside:

notification-service/internal/consumer/consumer.go

Then submit a payment.

Failed messages will appear in:

payment.completed.dlq

---

# Scaling Workers

Run multiple notification workers:

```bash
docker compose up --scale notification-service=3
```

RabbitMQ automatically distributes messages between workers.

---

# Production Concepts Demonstrated

* Event-driven architecture
* Distributed systems
* Message queues
* Background workers
* Reliable delivery
* Dead Letter Queues
* Cache-aside strategy
* Idempotent consumers
* Retry strategies
* Exponential backoff
* Graceful shutdown
* Service isolation
* Horizontal scalability

---

# Author

Tamerlan Khassenov SE-2416

Go • gRPC • RabbitMQ • Redis • PostgreSQL • Docker • Next.js

