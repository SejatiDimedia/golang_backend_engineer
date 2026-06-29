# Booking Management System (Coworking Space)

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://postgresql.org)

Booking Management System adalah sistem pemesanan ruangan dan meja kerja (*desk*) untuk coworking space terpadu. Proyek ini memvalidasi bentrokan waktu pemesanan (*double-booking*) secara aman menggunakan transaksi database PostgreSQL (`FOR UPDATE`) dan melatih pengiriman otentikasi JWT ad-hoc serta runtime middleware Gin di Go.

---

## 1. Tech stack

- **Core Engine:** Go (Golang) v1.20+
- **REST Framework:** Gin Web Framework
- **ORM:** GORM (v2)
- **Database:** PostgreSQL 15
- **Security:** JWT (`golang-jwt/jwt/v5`), Password Hashing (`bcrypt`)
- **Containerization:** Docker & Docker Compose

## 2. Key features

- **Otentikasi Ad-Hoc:** Registrasi & Login terproteksi password hashing, menghasilkan JWT token.
- **Katalog Meja/Ruangan:** Manajemen CRUD Meja/Ruangan coworking space.
- **Transaksi Booking & Overlap Check:** Validasi tumpang tindih slot waktu ($S_1 < E_2 \ \text{and} \ S_2 < E_1$) secara konkuren di bawah PostgreSQL pessimistic row-locking (`SELECT ... FOR UPDATE`).
- **Standardisasi UTC:** Seluruh tanggal pemesanan dikonversi secara ketat ke UTC saat masuk dan diproses dalam UTC.
- **Batasan Pembatalan:** User hanya dapat membatalkan pemesanan maksimal 2 jam sebelum waktu mulai.
- **Stub Notifikasi:** Log stdout placeholder sebagai stub untuk Notification Service (Project 6).

## 3. Quick start

Pastikan Docker Desktop aktif di komputer Anda.

1. **Setup Environment:**
   ```bash
   cp .env.example .env
   ```
2. **Start Database Container:**
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

Sebelum menguji secara lengkap, tinjau dokumen arsitektur dan spesifikasi teknis berikut:

* [PRD.md](./PRD.md) - Kebutuhan produk dan sukses kriteria.
* [ROADMAP.md](./ROADMAP.md) - Peta jalan urutan pengerjaan fitur.
* [ARCHITECTURE.md](./ARCHITECTURE.md) - Diagram data-flow dan struktur Clean Architecture.
* [DATABASE.md](./DATABASE.md) - Skema ERD dan indeks tabel.
* [API.md](./API.md) - Spesifikasi lengkap payload HTTP Endpoint.
* [SETUP.md](./SETUP.md) - Panduan instalasi dan pengujian cURL lokal.
* [DEPLOYMENT.md](./DEPLOYMENT.md) - Konfigurasi Docker multi-stage.
* [TESTING.md](./TESTING.md) - Strategi pengujian unit dan cakupan mock data.
* [CHANGELOG.md](./CHANGELOG.md) - Riwayat rilis perubahan.
* [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) - Rencana retrofit Auth terpusat dan Redis locking.
* [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) - Retrospektif penanganan UTC time-drift di Go.
* [ADR-001 (JWT Auth)](./adr/001-auth-strategy.md) - Justifikasi implementasi auth lokal.
* [ADR-002 (Overlap Checking)](./adr/002-overlap-checking-strategy.md) - Justifikasi database locking dibanding exclusion constraint.
