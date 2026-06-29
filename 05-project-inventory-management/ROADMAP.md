# Roadmap: Inventory Management API

**Status:** `Planning`

Dokumen ini mengurutkan urutan pengerjaan fitur dalam proyek Inventory Management API untuk mengendalikan kompleksitas transaksi database relasional.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Inisialisasi modul Go, docker-compose PostgreSQL 15 lokal, konfigurasi env, dan setup Gin router dasar. | — | `Planned` |
| 2 — Transaction Infrastructure | Implementasi `TransactionManager` berbasis context wrapper untuk mengelola transaksi database GORM secara transparan. | Phase 1 | `Planned` |
| 3 — Master Data CRUD | REST API CRUD untuk entitas `Category` dan `Supplier` beserta validasi payload dasar. | Phase 2 | `Planned` |
| 4 — Product Management | CRUD entitas `Product` dengan relasi database (foreign keys ke Category dan Supplier) dan penegasan RESTRICT constraint saat penghapusan master data. | Phase 3 | `Planned` |
| 5 — Stock Mutation (Transactions) | Implementasi API `Stock In` dan `Stock Out` yang berjalan di dalam transaksi database atomis untuk memperbarui kuantitas produk dan membuat catatan di `stock_movements`. | Phase 4 | `Planned` |
| 6 — History & Export | Endpoint riwayat mutasi stok berpaginasi dengan pemfilteran, serta endpoint untuk mengunduh laporan stok dalam format file CSV. | Phase 5 | `Planned` |
| 7 — Hardening | Penulisan unit test komprehensif pada Service Layer dan Handler, validasi parameter bilangan positif, dan standardisasi format API error. | All above | `Planned` |
| 8 — Deployment & Docs | Penyiapan multi-stage Dockerfile, docker-compose full-stack, pengisian 12 berkas dokumentasi wajib. | Phase 7 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| Category & Supplier CRUD | FR-1, FR-2 | Harus ada terlebih dahulu karena produk membutuhkan ID kategori dan supplier saat dibuat. |
| Product CRUD | FR-3 | Fondasi data master inventaris. Harus ada sebelum kita bisa melakukan mutasi stok (in/out) terhadap produk terkait. |
| Restrict Delete Validation | FR-4 | Validasi keamanan data relasional. Menghalangi terhapusnya kategori/supplier yang masih dirujuk oleh produk aktif. |
| Stock In (Transaction) | FR-5, FR-7 | Menguji fungsionalitas transaksi database atomis pertama kali dengan menambah stok dan mencatat mutasi masuk. |
| Stock Out (Transaction) | FR-6, FR-7 | Menguji logika transaksi dengan tambahan validasi bisnis: kuantitas stok tidak boleh menjadi negatif. |
| Stock Movement History | FR-8 | Membaca log histori mutasi yang dihasilkan dari transaksi Stock In/Out. Membutuhkan paginasi agar performa kueri optimal. |
| CSV Export | FR-9 | Mengekspor status inventaris produk saat ini ke format CSV menggunakan stream writer untuk performa tinggi. |
| Health Check | FR-10 | Untuk memantau keaktifan service dan koneksi database relasional. |

## 3. Concepts this project is exercising

- **SQL Transactions di Go:** Menjalankan transaksi secara deklaratif dari Service Layer tetapi terisolasi secara infrastruktur melalui context propagation.
- **Relational Integrity (Foreign Keys):** Menggunakan constraint database relasional secara disiplin (RESTRICT/NO ACTION) di tingkat database.
- **Kueri Gabungan (Joins) & Preloading:** Memanfaatkan fitur `Preload` GORM untuk menyajikan data relasi secara optimal.
- **Pagination & Query Filtering:** Menghindari pengambilan seluruh baris database ke memori dengan membatasi kueri menggunakan `Limit` dan `Offset`.
- **CSV Data Streaming:** Menghasilkan dokumen teks CSV secara langsung ke HTTP response writer guna menghemat penggunaan RAM pada data besar.

## 4. Known risks / unknowns at planning time

- **Kebocoran Transaksi (Transaction Leak):** Kegagalan menutup transaksi (tidak memanggil Rollback saat terjadi error, atau lupa Commit) dapat menyebabkan kebocoran koneksi database (*database pool exhaustion*). Masalah ini akan diantisipasi dengan penanganan blok `defer` rollback otomatis pada implementasi `TransactionManager`.
- **Dirty Reads / Concurrent Mutation:** Jika ada dua request mutasi stok bersamaan untuk produk yang sama, stock quantity bisa tidak akurat. Kita harus menggunakan baris database lock (klausul `SELECT ... FOR UPDATE` atau update atomic) untuk memastikan data konsisten.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek Inventory Management API |
