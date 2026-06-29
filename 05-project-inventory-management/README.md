# Inventory Management API

> Layanan API untuk mengelola data master persediaan produk, kategori, supplier, mutasi stok masuk/keluar secara transaksional (all-or-nothing), serta ekspor ringkasan persediaan ke file CSV.

**Status:** `Documented`
**Difficulty:** `Beginner→Intermediate`
**Part of:** [Golang Backend Roadmap](../README.md) — Project 2 of 8

---

## What this is

Inventory Management API memecahkan masalah pencatatan stok fisik di gudang. Layanan ini memastikan setiap mutasi stok (stok masuk dan stok keluar) tercatat riwayatnya secara sinkron dan konsisten di database PostgreSQL. Kuantitas stok produk dilindungi oleh validasi logika agar tidak pernah bernilai negatif, dan penguncian baris (*row locking*) digunakan untuk mencegah kondisi balapan (*race condition*) ketika beberapa transaksi terjadi secara bersamaan.

## Why this project exists in the sequence

Proyek ini berada pada Phase 1 (Foundations) dan membangun di atas REST API dasar dari Project 1. Proyek ini memperkenalkan relasi database yang sesungguhnya (kunci asing/foreign keys), operasi SQL Transactions terisolasi di tingkat Repository menggunakan pola context manager, paginasi data mutasi yang efisien, dan pembentukan data streaming file (.csv) langsung ke HTTP writer. Selengkapnya lihat di [01-roadmap.md](../01-roadmap.md).

## Core features

- **Master Data CRUD:** Mengelola data master `Category`, `Supplier`, dan `Product` (dengan relasi kategori dan supplier).
- **Relational Restriction:** Menghalangi penghapusan kategori atau supplier yang masih dirujuk oleh produk aktif.
- **Transactional Stock In/Out:** Melakukan penambahan/pengurangan stok produk secara aman di dalam kueri transaksi terisolasi PostgreSQL.
- **Stock Validation:** Memblokir transaksi stok keluar jika kuantitas stok saat ini kurang dari jumlah yang diminta (stok tidak boleh negatif).
- **Movement History & Pagination:** Melihat riwayat mutasi berpaginasi lengkap dengan filter id produk dan tipe mutasi.
- **CSV Laporan Persediaan:** Mengunduh ringkasan produk dan jumlah stok saat ini dalam format CSV secara streaming.

Detail persyaratan produk dapat dibaca di [PRD.md](./PRD.md).

## Tech stack

| Layer | Choice | Why (link to ADR if applicable) |
|---|---|---|
| Language | Go | Bahasa utama pembelajaran repository. |
| Framework | Gin | Router HTTP untuk API REST. |
| Database | PostgreSQL 15 | Menyimpan data terstruktur dan mendukung ACID Transactions. |
| ORM | GORM | Memudahkan pemetaan relasi data master dan migrasi. |
| Transaction | Context Manager | Mengisolasi transaksi tanpa mencemari Service Layer ([ADR-001](./adr/001-transaction-management.md)). |
| Containerization | Docker & Compose | Menjalankan instance database PostgreSQL 15 terisolasi. |

## Documentation index

| Document | Purpose |
|---|---|
| [PRD.md](./PRD.md) | Persyaratan produk — apa yang dibangun dan mengapa |
| [ROADMAP.md](./ROADMAP.md) | Urutan pengembangan fitur dalam proyek ini |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Desain arsitektur, pola transaction manager, alur data |
| [DATABASE.md](./DATABASE.md) | Skema database, relasi foreign keys, atomic locking |
| [API.md](./API.md) | Spesifikasi REST API endpoint |
| [SETUP.md](./SETUP.md) | Panduan instalasi lokal dan konfigurasi database |
| [DEPLOYMENT.md](./DEPLOYMENT.md) | Cara membungkus aplikasi dengan Docker |
| [TESTING.md](./TESTING.md) | Strategi pengujian unit transaksi dan rollback |
| [adr/](./adr/) | Architecture Decision Records (ADR-001) |
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

Semua endpoint fungsional, logika transaksi, penguncian baris database, streaming CSV, dan unit test telah lengkap ditulis.
