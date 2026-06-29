# File Management Service

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![MinIO](https://img.shields.io/badge/Storage-MinIO-blue.svg)](https://min.io)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://postgresql.org)

File Management Service adalah layanan backend untuk mengelola pengunggahan berkas secara aman menggunakan MinIO Object Storage (S3-compatible) dan database relasional PostgreSQL untuk pelacakan metadata. Proyek ini melatih integrasi SDK S3, penanganan HTTP Multipart Form-Data, pembuatan tautan unduhan privat menggunakan S3 Presigned URL, dan pengiriman berkas asinkron menggunakan chunked streaming I/O.

---

## 1. Tech stack

- **Core Engine:** Go (Golang) v1.20+
- **REST Framework:** Gin Web Framework
- **ORM:** GORM (v2)
- **Database:** PostgreSQL 15
- **Object Storage:** MinIO (RELEASE.2023-08-29)
- **Security:** JWT (`golang-jwt/jwt/v5`), Password Hashing (`bcrypt`)
- **Containerization:** Docker & Docker Compose

## 2. Key features

- **Auto-Bucket Creation:** Server secara otomatis mendeteksi ketersediaan bucket target (`user-files`) di MinIO saat booting dan membuatnya jika belum ada.
- **Multipart Upload Validator:** Memvalidasi berkas multipart form-data di tingkat routing (maksimal 10MB, tipe MIME diperbolehkan: JPEG, PNG, PDF).
- **Compensating Rollback:** Jika pengunggahan berkas fisik ke MinIO gagal, baris metadata `PENDING` di PostgreSQL otomatis dihapus (*rollback*) untuk menjamin konsistensi data.
- **S3 Presigned URL Download:** API `/files/:id/download` menghasilkan tautan bertanda tangan S3 dengan masa aktif 15 menit, membebaskan bandwidth server backend Go.
- **Direct Server Streaming:** API `/files/:id/view` mendukung render berkas langsung di tag HTML dengan mem-pipe byte stream MinIO ke client menggunakan `io.Copy`.

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
* [TESTING.md](./TESTING.md) - Strategi pengujian unit, validator multipart, dan compensating write.
* [CHANGELOG.md](./CHANGELOG.md) - Riwayat rilis perubahan.
* [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) - Rencana retrofit Auth terpusat dan chunked multipart upload.
* [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) - Retrospektif compensating write dan bandwidth offloading.
* [ADR-001 (Storage Strategy)](./adr/001-storage-strategy.md) - Justifikasi Object Storage dibanding Local Disk.
* [ADR-002 (Access Strategy)](./adr/002-access-strategy.md) - Justifikasi Presigned URL dibanding Direct Server Streaming.
