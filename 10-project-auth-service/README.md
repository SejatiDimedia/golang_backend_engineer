# Authentication Service (Auth Service)

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://postgresql.org)
[![RS256](https://img.shields.io/badge/Security-Asymmetric_RS256-green.svg)](https://jwt.io)

Auth Service adalah layanan Identity Provider (IdP) terpusat untuk platform microservices. Menyediakan pendaftaran, login aman, rotasi token refresh otomatis (RTR), pemulihan password, aktivasi verifikasi email, serta otorisasi berbasis Role-Based Access Control (RBAC). Token JWT ditandatangani menggunakan algoritma asimetris RS256 (Private/Public RSA Keypair) untuk mendukung offline verification di downstream services.

---

## 1. Tech stack

- **Core:** Go (Golang) v1.20+
- **REST Framework:** Gin Web Framework
- **ORM:** GORM (v2)
- **Database:** PostgreSQL 15 (SQLite in-memory for testing)
- **Token Caching:** Redis 7
- **Cryptography:** RS256 JWT, BCrypt Hashing, Crypto Secure Random Tokens
- **Autogenerator:** RSA 2048-bit Key Pair Auto-generator

## 2. Key features

- **Asymmetric RS256 Signing:** Auth Service memegang private key untuk menandatangani JWT. Downstream services (Booking, Wallet, Notification) dapat memverifikasi token offline hanya dengan public key secara instan ($O(1)$) tanpa interupsi network ke Auth Service.
- **Refresh Token Rotation (RTR):** Setiap kali refresh token digunakan, pasangan token access + refresh baru dikembalikan dan token lama dicabut.
- **RTR Anti-Replay Safeguard:** Jika peretas mencoba mengirim kembali token refresh yang *sudah di-revoke*, database transaction locking (`SELECT ... FOR UPDATE`) mendeteksinya dan otomatis mencabut seluruh sesi aktif milik user tersebut untuk memaksa logout massal peretas.
- **Native Relational RBAC:** Relasi PostgreSQL User -> Role -> Permission many-to-many dinamis untuk otorisasi hak akses terpadu.
- **Activation & Password Flows:** Mengelola alur token aktivasi verifikasi email dan lupa password.

## 3. Quick start

1. **Setup Environment:**
   ```bash
   cp .env.example .env
   ```
2. **Start Database & Redis Containers:**
   ```bash
   docker-compose up -d
   ```
3. **Run Server Local:**
   ```bash
   go run cmd/server/main.go
   ```
   *RSA certificate files `certs/private.key` & `certs/public.key` otomatis dibuat saat inisialisasi booting pertama kali.*
4. **Jalankan Unit Test:**
   ```bash
   go test -v ./...
   ```

## 4. Documentation index

* [PRD.md](./PRD.md) - Definisi kebutuhan produk, delivery guarantee, dan sukses kriteria.
* [ROADMAP.md](./ROADMAP.md) - Peta jalan pembangunan fitur platform.
* [ARCHITECTURE.md](./ARCHITECTURE.md) - Diagram data-flow offline verification dan RTR anti-replay.
* [DATABASE.md](./DATABASE.md) - Skema ERD relasi join table RBAC.
* [API.md](./API.md) - Spesifikasi payload REST HTTP Auth, Introspect, dan RBAC.
* [SETUP.md](./SETUP.md) - Panduan cURL setup lokal dan simulasi security replay attacks.
* [DEPLOYMENT.md](./DEPLOYMENT.md) - Konfigurasi Docker multi-stage.
* [TESTING.md](./TESTING.md) - Strategi pengujian SQLite in-memory, RTR replay simulation, dan offline RSA tests.
* [CHANGELOG.md](./CHANGELOG.md) - Riwayat perubahan versi.
* [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) - Rencana migrasi retrofit downstream services (Projects 3, 4, 6).
* [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) - Retrospektif penanganan transactional rollback dan database locks.
* [ADR-001 (Token Asymmetric)](./adr/001-asymmetric-tokens.md) - Justifikasi pemilihan RS256 RSA dibanding HS256 HMAC.
* [ADR-002 (RBAC Schema)](./adr/002-rbac-schema.md) - Justifikasi relational database tables dibanding Casbin library.
