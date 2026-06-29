# Notification Service

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![Redis](https://img.shields.io/badge/Queue-Redis-red.svg)](https://redis.io)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://postgresql.org)

Notification Service adalah layanan backend pengirim notifikasi (email, webhook, push) asinkron yang tangguh. Mengimplementasikan antrean pesan (message queue) manual Redis List & Sorted Set, daemon worker pool paralel, poller scheduler, log audit pengiriman, dan retry exponential backoff otomatis untuk menangani kegagalan eksternal.

---

## 1. Tech stack

- **Core Engine:** Go (Golang) v1.20+
- **REST Framework:** Gin Web Framework
- **ORM:** GORM (v2)
- **Database:** PostgreSQL 15
- **Message Broker:** Redis 7 (Alpine)
- **Security:** JWT (`golang-jwt/jwt/v5`), Password Hashing (`bcrypt`)
- **Containerization:** Docker & Docker Compose

## 2. Key features

- **Hand-Rolled Redis Queue:**
  - Instant queue (`LPUSH`/`BRPOP` blocking dequeue) untuk pesan instan.
  - Scheduled & retry queue (`ZADD`/`ZRANGEBYSCORE` Sorted Set) untuk menunda tugas.
- **Atomic Lua Scheduler:** Pemindahan tugas jatuh tempo dari Sorted Set ke instant List dijalankan menggunakan Lua script atomik untuk menjamin keselamatan konkurensi di lingkungan multi-node.
- **Worker Pool Daemon:** Concurrency pool paralel asinkron yang membongkar antrean di background thread Go tanpa menghambat API utama.
- **Exponential Backoff Retry:** Otomatis menghitung penundaan retry dinamis ($2^{\text{attempt}} \times 2$ detik) saat terjadi kegagalan eksternal, memindahkan notifikasi ke Dead-Letter Queue (status `FAILED`) jika melampaui 5 kali percobaan.
- **Audit Logs:** Log lengkap dari error message dan timestamps setiap percobaan pengiriman dicatat permanen ke PostgreSQL.

## 3. Quick start

Pastikan Docker Desktop aktif di komputer Anda.

1. **Setup Environment:**
   ```bash
   cp .env.example .env
   ```
2. **Start Containers:**
   ```bash
   docker-compose up -d
   ```
3. **Run Server Local:**
   ```bash
   go run cmd/server/main.go
   ```
4. **Jalankan Unit Test:**
   ```bash
   go test -v ./...
   ```

## 4. Documentation index

* [PRD.md](./PRD.md) - Kebutuhan produk dan sukses kriteria.
* [ROADMAP.md](./ROADMAP.md) - Peta jalan pembangunan fitur.
* [ARCHITECTURE.md](./ARCHITECTURE.md) - Diagram data-flow dan struktur Clean Architecture.
* [DATABASE.md](./DATABASE.md) - Skema ERD dan indeks tabel.
* [API.md](./API.md) - Spesifikasi lengkap payload HTTP Endpoint.
* [SETUP.md](./SETUP.md) - Panduan instalasi dan pengujian cURL lokal.
* [DEPLOYMENT.md](./DEPLOYMENT.md) - Konfigurasi Docker multi-stage.
* [TESTING.md](./TESTING.md) - Strategi pengujian unit, miniredis queue, dan backoff math.
* [CHANGELOG.md](./CHANGELOG.md) - Riwayat rilis perubahan.
* [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) - Rencana retrofit Auth terpusat dan webhook validation.
* [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) - Retrospektif atomic Lua scheduler dan blocking BRPOP.
* [ADR-001 (Queue Strategy)](./adr/001-queue-strategy.md) - Justifikasi Hand-rolled Redis List/ZSET dibanding library Asynq.
* [ADR-002 (Delivery Guarantee)](./adr/002-delivery-guarantee.md) - Justifikasi At-Least-Once delivery.
