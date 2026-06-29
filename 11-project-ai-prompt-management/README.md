# AI Prompt Management API

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Cache-Redis-red.svg)](https://redis.io)

AI Prompt Management API adalah repositori terpusat untuk mengelola, mengompilasi, dan melacak versi prompt template dinamis untuk LLM. Layanan ini mendukung arsitektur multi-tenancy berbasis Workspace, integrasi autentikasi API Key terdistribusi dengan Redis cache, verifikasi token JWT asimetris secara offline (mengonsumsi kunci publik milik Auth Service Project 7), dan tracking analitik pemakaian compiler secara asinkron.

---

## 1. Tech stack

- **Core:** Go (Golang) v1.20+
- **REST Framework:** Gin Web Framework
- **ORM:** GORM (v2)
- **Database:** PostgreSQL 15 (SQLite in-memory for testing)
- **Token & Key Caching:** Redis 7
- **Cryptography:** Asymmetric RS256 Validation, SHA-256 Key Hashing, Crypto Secure Random Generator
- **Daemon Logging:** Async Buffered Channels Analytics Workers

## 2. Key features

- **Asymmetric Offline JWT Verification:** Dashboard admin memvalidasi login JWT secara mandiri dengan membaca `certs/public.key` tanpa call RPC/Database sinkron ke Auth Service.
- **Secure API Key System:** Downstream client mengakses compiler prompt menggunakan header API Key (`prompt_live_...`). Hanya hash SHA-256 key yang disimpan di database PostgreSQL.
- **Redis Cache-Aside API Key Validation:** Hasil otentikasi API Key yang lolos validasi disimpan di Redis cache selama 1 jam (`apikey:<hash> -> workspace_id`), memotong latency validasi di bawah 2ms.
- **Prompt Version Snapshots:** Menyimpan teks instruksi AI secara utuh (*full snapshot*) sebagai record immutable per versi baru (v1, v2, v3, dst.) untuk mempercepat retrieval $O(1)$.
- **Template Compiler Engine:** Parser dinamis dengan regular expression (regex) Go untuk me-replace variabel template `{{var}}` dan mengestimasi panjang token ($\text{word count} \times 1.33$).
- **Async Analytics Worker Daemon:** Penulisan log data analitik (latensi, hit, token estimate) diproses non-blocking menggunakan buffered Go channels dan background daemon worker.

## 3. Quick start

1. **Copy Config:**
   ```bash
   cp .env.example .env
   ```
2. **Copy Public Key dari Project 7:**
   ```bash
   mkdir -p certs
   cp ../10-project-auth-service/certs/public.key certs/public.key
   ```
3. **Start Containers:**
   ```bash
   docker-compose up -d
   ```
4. **Run Server Local:**
   ```bash
   go run cmd/server/main.go
   ```
5. **Jalankan Unit Test:**
   ```bash
   go test -v ./...
   ```

## 4. Documentation index

* [PRD.md](./PRD.md) - Definisi kebutuhan produk, delivery guarantee, dan sukses kriteria.
* [ROADMAP.md](./ROADMAP.md) - Peta jalan pembangunan fitur platform.
* [ARCHITECTURE.md](./ARCHITECTURE.md) - Diagram offline JWT verification dan aliran kompilasi API Key.
* [DATABASE.md](./DATABASE.md) - Skema ERD relasi join table dan analitik.
* [API.md](./API.md) - Payload JSON HTTP dashboard admin dan prompt compiler.
* [SETUP.md](./SETUP.md) - Panduan cURL setup lokal dan simulasi kompilasi template.
* [DEPLOYMENT.md](./DEPLOYMENT.md) - Konfigurasi Docker multi-stage.
* [TESTING.md](./TESTING.md) - Strategi pengujian SQLite in-memory, API Key caching, dan regex compiler.
* [CHANGELOG.md](./CHANGELOG.md) - Riwayat perubahan versi.
* [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) - Rencana integrasi SDK client dan backup cache.
* [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) - Retrospektif penanganan buffered channels, SHA-256 caching, dan modular multi-tenancy.
* [ADR-001 (Versioning Strategy)](./adr/001-versioning-strategy.md) - Pilihan strategi full snapshot dibanding line deltas/diffs.
* [ADR-002 (API Key Caching)](./adr/002-apikey-caching.md) - Justifikasi caching Redis dibanding query DB relasional langsung.
