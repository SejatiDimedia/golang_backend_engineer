# Changelog: Inventory Management API

Semua perubahan penting pada proyek Inventory Management API dicatat di sini. Format penulisan didasarkan pada [Keep a Changelog](https://keepachangelog.com/).

---

## [1.0.0] — 2026-06-29 — Initial Implementation

### Added
- **Scaffolding Proyek:** Struktur folder Clean Architecture (cmd, internal/config, internal/entity, internal/handler, internal/repository, internal/service).
- **Transaction Manager Wrapper:** Pola `TransactionManager` berbasis context manager untuk mengisolasi koordinasi transaksi SQL (GORM `Begin`/`Commit`/`Rollback`) di tingkat repository, menjaga kesucian Service Layer ([ADR-001](./adr/001-transaction-management.md)).
- **Relational Schema:** Struktur 4 tabel terelasi PostgreSQL 15 (`categories`, `suppliers`, `products`, `stock_movements`) dengan constraint kunci asing (foreign keys) `RESTRICT` on delete pada Kategori dan Supplier.
- **REST API Endpoints:**
  - `GET /health` untuk pemeriksaan kesehatan backend & database.
  - CRUD Kategori (`POST`, `GET`, `PUT`, `DELETE /categories`).
  - CRUD Supplier (`POST`, `GET`, `PUT`, `DELETE /suppliers`).
  - CRUD Produk (`POST`, `GET`, `PUT`, `DELETE /products`) dengan validasi foreign keys dan paginasi data.
  - REST API Mutasi Stok (`POST /products/:id/stock-in` dan `/stock-out`) yang berjalan secara transaksional atomic.
  - REST API Riwayat Mutasi (`GET /stock-movements`) berpaginasi dengan filter.
- **Pessimistic Row Locking:** Penerapan locking database `SELECT ... FOR UPDATE` saat mutasi stok keluar (Stock Out) untuk memblokir race condition kuantitas produk.
- **Stock Validation:** Pencegahan kuantitas stok produk bernilai negatif.
- **CSV Data Streaming:** REST API endpoint `GET /products/export` untuk mengunduh laporan persediaan produk dalam format CSV secara streaming langsung ke socket writer HTTP.
- **Unit Testing:** Pengujian unit komprehensif (100% pass) untuk transaksi mutasi di `MovementService` dan HTTP parameter binding di `MovementHandler`.
- **Dockerization:** Multi-stage Dockerfile untuk runtime container aplikasi yang efisien dan aman.
- **Dokumentasi Standar:** Pembentukan 12 berkas dokumentasi standar lengkap.

### Notes
Seluruh pengerjaan diselesaikan secara terstruktur setelah memperoleh persetujuan rencana implementasi awal.
