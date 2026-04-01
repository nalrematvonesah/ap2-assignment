# AP2 Assignment 1 — Clean Architecture based Microservices (Order & Payment)

## 📌 Overview
This project implements a two-service microservice platform in **Go** using **Gin**, **PostgreSQL**, and **Clean Architecture**.

The system consists of:

- **Order Service** — manages customer orders and their lifecycle
- **Payment Service** — authorizes payments and validates transaction limits

Communication between services is implemented **strictly via REST**, as required in the assignment. :contentReference[oaicite:1]{index=1}

---

## 🏗 Architecture Decisions
Each service follows **Clean Architecture** and is split into the following layers:

- `domain/` → business entities
- `usecase/` → business logic and invariants
- `repository/` → database persistence logic
- `transport/http/` → thin HTTP handlers
- `app/` → route registration
- `cmd/` → composition root and manual dependency injection

### Why this architecture?
This design provides:

- clear separation of concerns
- dependency inversion
- testability
- maintainability
- isolated business rules
- easier service evolution

HTTP handlers are intentionally kept thin.  
All business rules and state transitions are implemented inside the **use case layer**, while SQL logic is isolated in repositories.

---

## 🧩 Bounded Contexts
The system is decomposed into **two bounded contexts**.

### 🟦 Order Context
Responsible for:

- creating orders
- retrieving orders
- cancelling orders
- updating order statuses

Owns:
- `orders` table
- `order_db`

### 🟩 Payment Context
Responsible for:

- payment authorization
- transaction creation
- declined payment rules
- payment lookup

Owns:
- `payments` table
- `payment_db`

### Important decomposition rules
This project intentionally avoids:

- ❌ shared database
- ❌ shared `common` package
- ❌ shared entities
- ❌ direct SQL access between services

This prevents the **distributed monolith anti-pattern**. :contentReference[oaicite:2]{index=2}

---

## 🔄 Service Communication Flow
The **Order Service** communicates with the **Payment Service** via REST:

```text
POST /orders
   ↓
Save order as Pending
   ↓
POST /payments
   ↓
Receive Authorized / Declined
   ↓
Update order status → Paid / Failed
```

The outbound communication uses a **custom `http.Client` with a 2-second timeout**, which satisfies the resiliency requirement. :contentReference[oaicite:3]{index=3}

### Failure Scenario
If the Payment Service is unavailable:

- the timeout is triggered
- the request does not hang
- the order status becomes **Failed**
- API returns **503 Service Unavailable**

This explicitly communicates payment failure to the client.

---

## 🗄 Database Design
Each microservice owns its own dedicated PostgreSQL database container.

### 🟦 Order Database
- Host: `127.0.0.1`
- Port: `55432`
- DB: `order_db`

### 🟩 Payment Database
- Host: `127.0.0.1`
- Port: `55433`
- DB: `payment_db`

This preserves **data ownership boundaries** and clean bounded contexts.

---

## 📊 Architecture Diagram
```text
┌───────────────┐
│  Frontend UI  │
└───────┬───────┘
        │ HTTP
        ▼
┌────────────────────┐
│  Order Service     │  :8080
│  - handlers        │
│  - use cases       │
│  - repository      │
└─────────┬──────────┘
          │ REST POST /payments
          ▼
┌────────────────────┐
│ Payment Service    │  :8081
│ - handlers         │
│ - use cases        │
│ - repository       │
└────────────────────┘

┌────────────────────┐   ┌────────────────────┐
│ order_db           │   │ payment_db         │
│ :55432             │   │ :55433             │
└────────────────────┘   └────────────────────┘
```

---

## 🚀 How to Run

### 1) Start PostgreSQL containers
```bash
docker compose up -d
```

### 2) Run migrations
#### PowerShell
```powershell
Get-Content .\order-service\migrations\001_create_orders.sql | docker exec -i order-db psql -U postgres -d order_db
Get-Content .\payment-service\migrations\001_create_payments.sql | docker exec -i payment-db psql -U postgres -d payment_db
```

### 3) Start Payment Service
```bash
cd payment-service
go run cmd/payment-service/main.go
```

### 4) Start Order Service
```bash
cd order-service
go run cmd/order-service/main.go
```

### 5) Start frontend
Open:

```text
frontend/index.html
```

using **Live Server** in VS Code.

---

## 📡 API Examples

### Create Order
```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{
  "customer_id":"cust-1",
  "item_name":"MacBook",
  "amount":50000
}'
```

### Get Order
```bash
curl http://localhost:8080/orders/{id}
```

### Cancel Order
```bash
curl -X PATCH http://localhost:8080/orders/{id}/cancel
```

### Payment Decline Example
```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{
  "customer_id":"cust-1",
  "item_name":"iPhone",
  "amount":150000
}'
```

Expected result: **Failed**, because amount > 100000. :contentReference[oaicite:4]{index=4}

---

## 💰 Business Rules
The following rules are enforced:

- money uses **int64**
- amount must be **greater than 0**
- amount **> 100000 → Declined**
- paid orders **cannot be cancelled**
- payment call timeout **≤ 2 seconds**

These invariants are implemented in the **use case layer**, not in handlers.

---

## ⚠️ Failure Handling Decision
If the Payment Service is unavailable, the Order Service updates the order status to **Failed**.

### Why `Failed` instead of `Pending`?
I intentionally chose `Failed` because it provides:

- explicit client feedback
- stronger consistency
- easier debugging
- deterministic order state

This decision simplifies client-side retry logic and makes the failure visible.

---

## 🎯 Assignment Compliance Checklist
This solution fully satisfies the assignment requirements:

- ✅ Clean Architecture inside each service
- ✅ Thin handlers
- ✅ business rules in use cases
- ✅ repository abstraction
- ✅ real PostgreSQL databases
- ✅ separate database per service
- ✅ no shared code/models
- ✅ REST-only communication
- ✅ custom timeout HTTP client
- ✅ failure scenario + 503
- ✅ all required endpoints
- ✅ business invariants
- ✅ SQL migration scripts
- ✅ frontend testing UI
- ✅ README and architecture diagram