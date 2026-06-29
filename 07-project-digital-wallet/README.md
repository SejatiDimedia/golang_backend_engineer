# Digital Wallet API

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Cache-Redis-red.svg)](https://redis.io)

Digital Wallet API adalah layanan backend dompet digital transaksional yang mengedepankan konsistensi finansial data saldo. Proyek ini mengimplementasikan pencatatan akuntansi Double-Entry Ledger, distributed locking Redis hand-rolled untuk memblokir double-spending konkuren, proteksi idempotensi request, dan caching saldo dompet digital.

---

## 1. Tech stack

- **Core Engine:** Go (Golang) v1.20+
- **REST Framework:** Gin Web Framework
- **ORM:** GORM (v2)
- **Database:** PostgreSQL 15
- **Cache & Distributed Mutex:** Redis 7 (Alpine)
- **Security:** JWT (`golang-jwt/jwt/v5`), Password Hashing (`bcrypt`)
- **Containerization:** Docker & Docker Compose

## 2. Key features

- **Otentikasi JWT Ad-Hoc:** Registrasi akun otomatis menginisialisasi satu `Wallet` dengan nomor rekening unik.
- **Double-Entry Ledger Bookkeeping:** Setiap mutasi saldo (top-up, withdraw, transfer) wajib dicatat sebagai entri debit/kredit yang seimbang di database relasional.
- **Redis Hand-Rolled Distributed Lock:** Mengunci ID wallet konkuren menggunakan command primitive `SET key token NX PX` dan pelepasan kunci Lua script atomik untuk mencegah double-spending.
- **Idempotency Key Middleware:** Header `X-Idempotency-Key` menyaring request retries ganda dengan menyimpan salinan respons sukses di Redis cache (TTL 1 jam).
- **Balance Caching:** Kueri saldo (`GET /wallet/balance`) dilayani oleh Redis cache, dengan fitur otomatis invalidasi (penghapusan cache) saat mutasi ledger baru sukses ditulis.

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
* [TESTING.md](./TESTING.md) - Strategi pengujian unit, concurrency test, dan idempotency cache.
* [CHANGELOG.md](./CHANGELOG.md) - Riwayat rilis perubahan.
* [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) - Rencana retrofit Auth terpusat dan multi-node cluster locking.
* [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) - Retrospektif deadlock dan Lua scripts.
* [ADR-001 (Ledger Model)](./adr/001-ledger-model-strategy.md) - Justifikasi running balance dibanding SUM query ledger.
* [ADR-002 (Redis Lock)](./adr/002-redis-lock-strategy.md) - Justifikasi penulisan lock manual dibanding library Redlock.
