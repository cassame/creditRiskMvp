# 🏦 Credit Risk Assessment API

A REST API service for automated credit application processing with risk assessment capabilities. The system evaluates loan applications through configurable validation strategies and makes automated decisions: approve, reject, or route to manual review.

### Key Features

- ✅ **Multi-strategy validation** - Different rule sets for residents/non-residents and first-time/repeat customers
- 🔄 **Parallel check execution** - Concurrent validation for improved performance
- 💾 **Transactional data persistence** - Atomic storage of applications and check results
- 🚀 **External service integration** - Credit history, terrorist watchlist, and bankruptcy checks
- 📊 **Event streaming** - Kafka integration for application status notifications
- 🎯 **Intelligent caching** - In-memory cache with TTL for watchlist data

## 🏗️ Architecture

```mermaid
graph TB
    Client[Client Application]
    
    subgraph "HTTP Layer"
        API[REST API Handlers<br/>POST /applications<br/>GET /applications/:id]
    end
    
    subgraph "Business Logic Layer"
        Strategy[Strategy Selection<br/>chooseStrategy]
        Runner[Parallel Check Execution<br/>runStrategy]
        Decision[Decision Engine<br/>decideStatus]
    end
    
    subgraph "Validation Checks"
        Local[Local Validators<br/>• Age 18+<br/>• Phone Format<br/>• Passport<br/>• Patronymic<br/>• Amount Limit]
        External[External Services<br/>• Credit History API<br/>• Terrorist Watchlist<br/>• Bankruptcy Check]
    end
    
    subgraph "Data Layer"
        DB[(PostgreSQL<br/>applications<br/>check_results)]
        Cache[In-Memory Cache<br/>Terrorist List TTL]
    end
    
    subgraph "Event Streaming"
        Kafka[Apache Kafka<br/>Application Events]
    end
    
    Client -->|HTTP Request| API
    API --> Strategy
    Strategy --> Runner
    Runner --> Local
    Runner --> External
    External -.->|Cache Hit| Cache
    Local --> Decision
    External --> Decision
    Decision --> DB
    Decision --> Kafka
    DB -->|Response| API
    API -->|HTTP Response| Client
    
    style API fill:#e1f5ff
    style Strategy fill:#fff3e0
    style Runner fill:#fff3e0
    style Decision fill:#fff3e0
    style DB fill:#f3e5f5
    style Kafka fill:#e8f5e9
    style Cache fill:#fce4ec
```

## 🔄 Request Flow

1. **POST /applications** - Client submits loan application
2. **parseApplication** - JSON validation and data extraction
3. **chooseStrategy** - Select validation strategy based on customer profile
4. **runStrategy** - Execute validation checks in parallel
5. **decideStatus** - Determine final decision (approved/rejected/manual_review)
6. **saveApplication** - Persist to database within transaction
7. **notifier.Notify** - Publish event to Kafka
8. **Response** - Return decision to client

## 🎯 Validation Strategies

| Customer Type | First Time | Checks Applied |
|--------------|------------|----------------|
| **Resident** | ✅ Yes | age, phone, passport, patronymic, amount, terrorist, credit_history |
| **Resident** | ❌ No | age, phone, passport, amount |
| **Non-Resident** | ✅ Yes | age, phone, passport, amount, terrorist, credit_history |
| **Non-Resident** | ❌ No | age, phone, amount |

### Decision Logic

- **Approved** ✅ - All checks passed
- **Rejected** ❌ - Any check failed
- **Manual Review** ⚠️ - Critical external service unavailable (terrorist/bankruptcy check)


## 📦 Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/credit-risk-api.git
cd credit-risk-api
```

2. **Start infrastructure services**
```bash
docker-compose up -d
```

3. **Initialize database schema**
```bash
psql postgres://app:app@localhost:5432/creditrisk < schema.sql
```

4. **Run the application**
```bash
go run .
```

The API will be available at `http://localhost:8080`

## 🔧 Configuration

Environment variables:

```bash
DATABASE_URL=postgres://app:app@localhost:5432/creditrisk?sslmode=disable
HTTP_TIMEOUT_MS=1000
CREDIT_HISTORY_URL=http://localhost:8080/mock/credit-history
TERRORIST_LIST_URL=http://localhost:8080/mock/terrorist/list
BANKRUPTCY_URL=http://localhost:8080/mock/bankruptcy
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=credit-applications
```

## 📡 API Endpoints

### Create Application

```http
POST /applications
Content-Type: application/json

{
  "name": "Ivanov Ivan Ivanovich",
  "birthdate": "1990-05-15",
  "phone": "+79161234567",
  "passport": "1234 567890",
  "residency": "resident",
  "first_time": true,
  "requested_amount": 50000
}
```

**Response:**
```json
{
  "application_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "approved",
  "strategy": "resident_first_time",
  "checks": [
    {
      "check": "age>=18",
      "status": "passed",
      "reason": ""
    },
    {
      "check": "valid_phone",
      "status": "passed",
      "reason": ""
    }
  ]
}
```

### Get Application

```http
GET /applications/{id}
```

**Response:**
```json
{
  "application_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-02-04T10:30:00Z",
  "strategy": "resident_first_time",
  "status": "approved",
  "payload": { },
  "checks": [  ]
}
```

## 🧪 Testing

### Test application submission
```bash
curl -X POST http://localhost:8080/applications \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User Testovich",
    "birthdate": "1995-01-01",
    "phone": "+79161234567",
    "passport": "1234 567890",
    "residency": "resident",
    "first_time": true,
    "requested_amount": 75000
  }'
```

## 👨‍💻 Author
**cassame**
- GitHub: [@cassame](https://github.com/cassame)
- LinkedIn: [cassame](https://linkedin.com/in/cassame)

