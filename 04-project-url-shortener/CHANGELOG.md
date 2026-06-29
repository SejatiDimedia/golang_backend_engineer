# Changelog: URL Shortener Service

Semua perubahan penting pada proyek URL Shortener Service dicatat di sini. Format penulisan didasarkan pada [Keep a Changelog](https://keepachangelog.com/).

---

## [1.0.0] — 2026-06-29 — Initial Implementation

### Added
- **Scaffolding Proyek:** Struktur folder Clean Architecture (cmd, internal/config, internal/entity, internal/handler, internal/repository, internal/service).
- **Inisiasi Database:** Integrasi koneksi GORM ke PostgreSQL 15 dan docker-compose local dev, dilengkapi dengan AutoMigration skema URL.
- **REST API Endpoints:** 
  - `GET /health` untuk pemeriksaan kesehatan backend & database.
  - `POST /shorten` untuk memperpendek URL panjang dengan validasi HTTP body, kustom alias, dan waktu kedaluwarsa.
  - `GET /r/:short_code` untuk pengalihan URL (HTTP 302) dengan *concurrency-safe click counting* dan pengecekan masa berlaku.
  - `GET /stats/:short_code` untuk memperoleh statistik klik dan metadata URL.
- **Algoritma Short Code:** Pustaka pembuatan short code menggunakan URL-safe Base64 encoding dari timestamp nanodetik (independen dari round-trip database).
- **Unit Testing:** Pengujian unit terisolasi (100% pass) untuk logika bisnis Service layer dan fungsionalitas HTTP Handler layer menggunakan HTTP recorder mock.
- **Dockerization:** Multi-stage Dockerfile untuk runtime container aplikasi yang efisien dan aman.
- **Dokumentasi Standar:** Pembentukan dokumen dasar: `README.md`, `PRD.md`, `ROADMAP.md`, `ARCHITECTURE.md`, `DATABASE.md`, `API.md`, `SETUP.md`, `DEPLOYMENT.md`, `TESTING.md`, dan `CHANGELOG.md`.

### Notes
Pengerjaan proyek ini diselesaikan dalam satu siklus pengerjaan terintegrasi setelah mendapatkan persetujuan terhadap rencana implementasi awal.
