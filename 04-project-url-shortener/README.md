# URL Shortener Service

> Layanan API untuk memperpendek URL panjang menjadi kode unik yang pendek (short code), melacak jumlah klik kunjungan, dan mendukung waktu kedaluwarsa tautan.

**Status:** `Documented`
**Difficulty:** `Beginner`
**Part of:** [Golang Backend Roadmap](../README.md) — Project 1 of 8

---

## What this is

URL Shortener Service memecahkan masalah pembagian tautan web yang panjang agar menjadi ringkas dan mudah dibaca. Layanan ini memungkinkan pengguna menghasilkan URL pendek (baik secara acak maupun kustom alias) untuk mempermudah distribusi konten, serta menyediakan statistik sederhana berupa jumlah klik per tautan pendek.

## Why this project exists in the sequence

Proyek ini adalah fondasi pertama (Phase 1) dalam roadmap kita. Proyek ini mempraktikkan materi dasar dari Phase 0 (structs, interfaces, dan error handling Go) dan mengenalkan REST routing menggunakan Gin, ORM menggunakan GORM dengan PostgreSQL, konfigurasi lingkungan lokal dengan docker-compose, dan scaffolding Clean Architecture dasar. Selengkapnya lihat di [01-roadmap.md](../01-roadmap.md).

## Core features

- **Shorten URL:** Mengonversi URL panjang menjadi short code unik sepanjang 11 karakter berbasis timestamp nanodetik.
- **Custom Alias:** Menggunakan kata kustom sebagai short code (misal `/r/timur-portofolio`) selama alias tersebut masih tersedia.
- **Expiration Date:** Membatasi akses URL pendek dengan memberikan waktu kedaluwarsa.
- **Click Counter:** Melacak total kunjungan ke URL pendek secara real-time.
- **Health Check:** Menyediakan endpoint untuk pemantauan server dan koneksi database.

Detail persyaratan produk dapat dibaca di [PRD.md](./PRD.md).

## Tech stack

| Layer | Choice | Why (link to ADR if applicable) |
|---|---|---|
| Language | Go | Bahasa utama pembelajaran repository. |
| Framework | Gin | Router HTTP yang cepat dan populer di ekosistem Go. |
| Database | PostgreSQL 15 | Database relasional untuk menyimpan data URL pendek secara presisten. |
| ORM | GORM | Memudahkan auto-migration dan operasi CRUD dasar ([ADR-001](./adr/001-orm-decision.md)). |
| Containerization | Docker & Compose | Menyediakan environment PostgreSQL terisolasi untuk lokal development. |

## Documentation index

| Document | Purpose |
|---|---|
| [PRD.md](./PRD.md) | Persyaratan produk — apa yang dibangun dan mengapa |
| [ROADMAP.md](./ROADMAP.md) | Urutan pengembangan fitur dalam proyek ini |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Desain sistem, struktur Clean Architecture, alur data |
| [DATABASE.md](./DATABASE.md) | Skema database, auto-migration, relasi data |
| [API.md](./API.md) | Spesifikasi REST API endpoint |
| [SETUP.md](./SETUP.md) | Panduan instalasi dan menjalankan server lokal |
| [DEPLOYMENT.md](./DEPLOYMENT.md) | Cara membungkus aplikasi dengan Docker |
| [TESTING.md](./TESTING.md) | Strategi pengujian unit dan cakupan testing |
| [adr/](./adr/) | Architecture Decision Records (ADR-001, ADR-002) |
| [CHANGELOG.md](./CHANGELOG.md) | Riwayat versi aplikasi |
| [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md) | Rencana pengembangan fitur yang ditunda |
| [LESSONS-LEARNED.md](./LESSONS-LEARNED.md) | Retrospeksi dan hasil pembelajaran |

## Quick start

```bash
# 1. Jalankan PostgreSQL lokal
docker-compose up -d

# 2. Salin environment config
cp .env.example .env

# 3. Jalankan unit test
go test ./...

# 4. Jalankan server backend Go
go run cmd/server/main.go
```

## Status notes

Sasis backend, koneksi database, routing REST, generator short-code, serta seluruh unit test telah selesai diimplementasikan.
