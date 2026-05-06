````md
# 🚀 Event-Driven Payment Processing System

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go)
![gRPC](https://img.shields.io/badge/gRPC-Microservices-244c5a?style=for-the-badge&logo=grpc)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-EventDriven-FF6600?style=for-the-badge&logo=rabbitmq)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=for-the-badge&logo=docker)
![Next.js](https://img.shields.io/badge/Next.js-Frontend-black?style=for-the-badge&logo=next.js)

Production-style distributed microservices system built with Go, gRPC, RabbitMQ, PostgreSQL, Docker, and Next.js.

</div>

---

# ✨ Features

- ⚡ gRPC communication
- 📨 RabbitMQ event-driven architecture
- 🗄 PostgreSQL persistence
- 🌐 API Gateway
- 🎨 Modern Next.js frontend
- 🐳 Dockerized infrastructure
- ♻ Retry connection logic
- ✅ Manual ACK handling
- 🛡 Idempotency protection
- 🧹 Graceful shutdown
- 🔥 Production-style architecture

---

# 🏗 Architecture

```text
┌──────────────────────┐
│      Frontend        │
│      Next.js         │
└──────────┬───────────┘
           │ HTTP
           ▼
┌──────────────────────┐
│      API Gateway     │
│      Gin (Go)        │
└──────────┬───────────┘
           │ gRPC
           ▼
┌──────────────────────┐
│    Payment Service   │
│        Go            │
└───────┬───────┬──────┘
        │       │
        │       │
        ▼       ▼
┌────────────┐  ┌──────────────┐
│ PostgreSQL │  │  RabbitMQ    │
└────────────┘  └──────┬───────┘
                       │ Event
                       ▼
              ┌──────────────────┐
              │ Notification     │
              │ Service          │
              └──────────────────┘
````

---

# 🧩 Services

## 🌐 Frontend

Modern UI built with:

* Next.js
* TailwindCSS
* React

Responsibilities:

* Payment form
* User interaction
* API communication

---

## 🚪 API Gateway

Built with:

* Go
* Gin

Responsibilities:

* HTTP API
* Request validation
* gRPC client communication
* CORS handling

---

## 💳 Payment Service

Built with:

* Go
* gRPC
* PostgreSQL
* RabbitMQ

Responsibilities:

* Payment processing
* Database persistence
* Event publishing

---

## 🔔 Notification Service

Built with:

* Go
* RabbitMQ

Responsibilities:

* Event consumption
* Notification handling
* Manual ACK processing
* Idempotency validation

---

# 📨 Event Flow

```text
Frontend
   ↓ HTTP
Gateway
   ↓ gRPC
Payment Service
   ↓
Store payment in PostgreSQL
   ↓
Publish event to RabbitMQ
   ↓
Notification Service consumes event
```

---

# 🛡 Reliability Features

## ✅ Manual ACK

Messages are acknowledged only after successful processing.

```go
msg.Ack(false)
```

Prevents message loss.

---

## ♻ Retry Logic

Services retry connections automatically during startup.

Improves resilience in distributed environments.

---

## 🛡 Idempotency

Duplicate events are ignored using processed event tracking.

---

## 🧹 Graceful Shutdown

Services correctly close:

* gRPC servers
* RabbitMQ channels
* PostgreSQL connections

---

# 🐳 Dockerized Infrastructure

The entire system runs inside Docker containers using Docker Compose.

Services:

* frontend
* gateway
* payment-service
* notification-service
* rabbitmq
* postgres

---

# 📁 Project Structure

```text
ap2-assignment/
│
├── frontend/
├── gateway/
├── payment-service/
├── notification-service/
├── order-service/
├── ap2-proto-contracts/
│
└── docker-compose.yml
```

---

# ⚙️ Technologies Used

| Technology  | Purpose                     |
| ----------- | --------------------------- |
| Go          | Backend services            |
| gRPC        | Inter-service communication |
| RabbitMQ    | Event broker                |
| PostgreSQL  | Database                    |
| Docker      | Containerization            |
| Next.js     | Frontend                    |
| TailwindCSS | UI styling                  |

---

# 🚀 Getting Started

## 1️⃣ Clone repository

```bash
git clone <your-repository>
```

---

## 2️⃣ Start system

```bash
docker compose up --build
```

---

# 🌐 Available Services

| Service     | URL                                              |
| ----------- | ------------------------------------------------ |
| Frontend    | [http://localhost:3000](http://localhost:3000)   |
| Gateway API | [http://localhost:8080](http://localhost:8080)   |
| RabbitMQ UI | [http://localhost:15672](http://localhost:15672) |

---

# 🔑 RabbitMQ Credentials

```text
Username: guest
Password: guest
```

---

# 🧪 Example Workflow

1. Open frontend
2. Submit payment
3. Gateway sends gRPC request
4. Payment stored in PostgreSQL
5. Event published to RabbitMQ
6. Notification Service processes event

---

# 📸 UI Preview

```text
Modern glassmorphism payment dashboard
with animated gradients and live processing state
```

---

# 🔥 Production Concepts Demonstrated

* Distributed Systems
* Event-Driven Architecture
* API Gateway Pattern
* Producer / Consumer Pattern
* Retry Pattern
* Idempotency Pattern
* Graceful Shutdown Pattern

---

# 🚀 Future Improvements

* JWT Authentication
* Redis caching
* Kubernetes deployment
* CI/CD pipelines
* Prometheus & Grafana monitoring
* WebSocket live events
* Email notifications
* Dead Letter Queues (DLQ)

---

# 👨‍💻 Author

Tamerlan Khassenov SE-2416

---

# ⭐ Project Status

```text
Production-style educational microservices system
```
