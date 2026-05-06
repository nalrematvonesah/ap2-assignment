````md
# Event-Driven Payment Processing System

A distributed microservices system built with Go, gRPC, RabbitMQ, and Docker.

The project demonstrates synchronous and asynchronous communication between services using an event-driven architecture.

---

# Architecture

```text
Order Service
      │
      │ gRPC
      ▼
Payment Service
      │
      │ RabbitMQ Event
      ▼
Notification Service
````

---

# Technologies Used

* Go
* gRPC
* Protocol Buffers
* RabbitMQ
* Docker
* Docker Compose

---

# Services

## 1. Order Service

Acts as a gRPC client.

Responsibilities:

* Sends payment requests to Payment Service
* Demonstrates synchronous communication

---

## 2. Payment Service

Acts as a gRPC server and RabbitMQ producer.

Responsibilities:

* Processes payment requests
* Publishes payment events to RabbitMQ
* Uses durable queues
* Implements retry connection logic
* Supports graceful shutdown

---

## 3. Notification Service

Acts as RabbitMQ consumer.

Responsibilities:

* Consumes payment events
* Sends notifications
* Uses manual ACK
* Implements idempotency
* Supports graceful shutdown

---

# RabbitMQ Event Flow

1. Order Service sends gRPC request
2. Payment Service processes payment
3. Payment Service publishes event to RabbitMQ
4. Notification Service consumes event
5. Notification Service acknowledges message

---

# Reliability Features

## Manual ACK

The Notification Service acknowledges messages only after successful processing.

```go
msg.Ack(false)
```

This prevents message loss.

---

## Durable Queue

RabbitMQ queue is configured as durable.

```go
QueueDeclare(
    "payment.completed",
    true,
)
```

Messages survive broker restarts.

---

## Retry Logic

Services retry RabbitMQ connection during startup.

This improves resilience when containers start simultaneously.

---

## Idempotency

Duplicate events are prevented using unique event IDs.

```go
seen map[string]bool
```

The Notification Service ignores already processed events.

---

# Graceful Shutdown

Services handle SIGINT/SIGTERM signals and close resources properly.

Features:

* gRPC graceful stop
* RabbitMQ connection cleanup
* channel cleanup

---

# Project Structure

```text
ap2-assignment/
│
├── docker-compose.yml
│
├── order-service/
│
├── payment-service/
│
├── notification-service/
│
└── ap2-proto-contracts/
```

---

# How to Run

## 1. Clone repositories

```bash
git clone <repo>
```

---

## 2. Start system

```bash
docker compose up --build
```

---

# RabbitMQ Management UI

Available at:

```text
http://localhost:15672
```

Credentials:

```text
guest / guest
```

---

# Expected Workflow

Successful execution:

```text
Order Service
    ↓
Payment processed
    ↓
Event published
    ↓
Notification received
```

---

# Example Logs

## Payment Service

```text
Processing payment...
Published event...
```

## Notification Service

```text
Received event...
Sending notification...
```

---

# Design Patterns Used

* Event-Driven Architecture
* Producer / Consumer Pattern
* Retry Pattern
* Idempotency Pattern
* Graceful Shutdown Pattern

---

# Future Improvements

* Dead Letter Queue (DLQ)
* Persistent idempotency storage
* Health checks
* Structured logging
* CI/CD pipeline
* Kubernetes deployment

---

# Author

Tamerlan Khassenov SE-2416