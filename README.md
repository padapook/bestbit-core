# BestBit-Core (Backend)

### Project Description
BestBit-Core คือระบบฝั่ง Backend ที่เขียนขึ้นมาเพื่อศึกษาการทำงานของระบบ Digital Asset Trading Platform

1. Financial Precision: ใช้การคำนวณทศนิยม 16 หลัก เพื่อรองรับทุก currency และป้องกัน Rounding errors
2. High Concurrency: ออกแบบ Architech มาเพื่อรองรับคำสั่งซื้อขาย (Orders) จำนวนมากพร้อมกัน โดยใช้ระบบ Hybrid ระหว่าง GORM และ pgxpool เพื่อ performance สูงสุดจาก Database
3. System Integrity: ระบบ Lock เงิน (Amount Locked) และ Transactional Integrity เพื่อป้องกันปัญหา Double-spending และรักษาความถูกต้องของยอดเงิน

### Key Objectives
- Wallet System: จัดการยอดเงิน (Available/Locked) ของ User
- Order Management: ระบบวางคำสั่งซื้อขายแบบ Real-time (Limit/Market Order)
- Matching Engine Integration: เตรียมความพร้อมสำหรับการเชื่อมต่อกับระบบจับคู่คำสั่งซื้อขาย
- Auditability: คำสั่งซื้อขายต้องสามารถตรวจสอบย้อนกลับได้ทั้งหมด

---

## Project Structure (Modular Monolith)
- cmd/server/main.go → Entry point ของระบบและการตั้งค่า Middleware
- internal/database/ → จัดการ GormConnectDB และ AutoMigrate
- internal/account/.../ → ข้อมูล User และ Profile (Singular naming)
- internal/wallet/.../ → ข้อมูล Balance, Locked และ Transaction
- internal/order/.../ → ข้อมูล Limit/Market Orders
- internal/trade/.../ → ข้อมูลการจับคู่ซื้อขาย (Match results)
- internal/routes/ → จัดการ Route Grouping (v1/api/...)

## Tech Specification
- Language: Golang (Gin)
- DB: PostgreSQL
- ORM & Driver: GORM (Schema/Migration) and pgxpool (High-speed Raw SQL).

### Performance & Scalability Layers
- In-Memory Cache: Redis
- Message Broker: RabbitMQ
- WebSocket
- Distributed Locking

### Financial Standards
- Precision: decimal(32,16) for all balance and transaction fields to ensure zero rounding errors.
- Naming Convention: Singular file and struct names (Go Best Practices).
- Data Integrity: Database-level constraints combined with ACID-compliant transactions.

### Testing Stack
- testify

---

# Development Guidelines (BestBit-Core)

### 1. Commit Message Pattern (Conventional Commits)

Follow this conventional commit format:

```
<type>(<scope>): <description>
```

**Types:**
- `feat` - A new feature
- `fix` - A bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, missing semicolons, etc.)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks
- `add` - Adding new files or dependencies

**Scope:** The scope should indicate what part of the codebase is affected (e.g., api, docs, readme, auth, database, etc.)

**Examples:**

❌ Bad:
```
feat: auth
done
done: Added JWT token-based authentication for secure login.
```

✅ Good:
```
feat(auth): Added JWT token-based authentication for secure login
fix(database): Fixed connection pool timeout issue
docs(readme): Updated installation instructions
refactor(controller): Simplified user validation logic
test(service): add thoroughness tests for withdrawal
```

---

### 2. Naming Conventions (Go Idioms)

#### Visibility (Public/Private)
- Capitalized (PascalCase) = Exported (Public)
- Lowercase (camelCase) = Unexported (Private)

---

### 3. Project Focus
- All monetary and balance-related fields must use the decimal(32,16) database type.
- Every core model must include standard Audit fields: CreatedAt, UpdatedAt, and DeletedAt.
- Sensitive fields, such as Password or internal keys, must use the json:"-" tag.