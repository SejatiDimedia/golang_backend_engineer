# PRD: Inventory Management API

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Staff gudang dan pemilik toko sering mengalami kesulitan melacak persediaan fisik barang secara real-time. Tanpa pencatatan yang disiplin, stok barang bisa habis tanpa diketahui (*out of stock*) atau justru menumpuk terlalu banyak. Selain itu, jika perubahan stok (penambahan/pengurangan) tidak dicatat bersama riwayatnya secara atomis di database, data kuantitas barang saat ini dengan riwayat mutasi barang bisa tidak sinkron, menyebabkan kebingungan laporan keuangan.

## 2. Goals

- Memungkinkan staff gudang untuk mengelola data master Produk, Kategori, dan Supplier.
- Menyediakan mekanisme mutasi stok (stok masuk dan stok keluar) yang aman secara transaksional di database PostgreSQL.
- Menyediakan pencarian data inventaris yang mendukung paginasi dan filter kategori.
- Menyediakan fitur ekspor laporan ringkasan stok dalam format file CSV.
- Memperdalam pemikiran relasional database (foreign keys) dan manajemen transaksi SQL (commit/rollback) di Go.

## 3. Non-goals

- **Manajemen Akun Gudang:** Tidak ada autentikasi pengguna atau pembagian hak akses (role) staff vs manajer di versi awal ini (disimpan untuk Project 3 & 7).
- **Multi-Warehouse (Banyak Gudang):** Layanan ini mengasumsikan hanya ada satu lokasi gudang penyimpanan fisik tunggal.
- **Prediksi Stok Otomatis (AI Forecasting):** Layanan ini hanya mencatat data mutasi saat ini secara presisi, tidak memprediksi kebutuhan stok di masa depan.

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Warehouse Staff | Melakukan pencatatan ketika barang baru datang dari supplier (Stock In) atau barang keluar ke pelanggan (Stock Out). | Setiap hari |
| Inventory Manager | Memantau persediaan saat ini, memeriksa riwayat mutasi stok, dan mengekspor laporan CSV untuk evaluasi. | Mingguan |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | Pengguna dapat melakukan operasi CRUD untuk data master Kategori (`Category`). | Must |
| FR-2 | Pengguna dapat melakukan operasi CRUD untuk data master Supplier (`Supplier`). | Must |
| FR-3 | Pengguna dapat melakukan operasi CRUD untuk data master Produk (`Product`), termasuk menghubungkannya ke Kategori dan Supplier (menggunakan relasi database). | Must |
| FR-4 | Sistem tidak boleh membolehkan penghapusan Kategori atau Supplier jika masih ada Produk aktif yang terhubung dengannya (Constraint RESTRICT/NO ACTION). | Must |
| FR-5 | Pengguna dapat melakukan penambahan stok (`Stock In`) dengan menentukan jumlah barang dan mencatat supplier asal. | Must |
| FR-6 | Pengguna dapat melakukan pengurangan stok (`Stock Out`) dengan menentukan jumlah barang. Sistem harus menolak transaksi jika jumlah stok saat ini tidak mencukupi (stok tidak boleh negatif). | Must |
| FR-7 | Setiap mutasi stok (Stock In/Out) wajib mencatat entri riwayat baru ke tabel `stock_movements` dan mengupdate `stock_quantity` di tabel `products` secara atomis di dalam satu transaksi database (All or Nothing). | Must |
| FR-8 | Pengguna dapat melihat daftar histori mutasi stok yang mendukung filter berdasarkan produk, tipe mutasi (in/out), serta mendukung paginasi (`page` & `limit`). | Must |
| FR-9 | Pengguna dapat mengunduh ringkasan persediaan produk saat ini dalam bentuk file CSV (`GET /products/export`). | Should |
| FR-10 | Menyediakan endpoint `GET /health` untuk memantau status aplikasi dan koneksi database. | Must |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Performance | Waktu eksekusi mutasi stok < 150ms. Pembacaan daftar produk dengan paginasi harus memiliki indeks yang sesuai. |
| Security | Validasi input kuantitas barang (harus bernilai bilangan bulat positif > 0) untuk mencegah manipulasi data stok negatif. |
| Data consistency | Menggunakan isolasi transaksi PostgreSQL (default: Read Committed) yang aman untuk memastikan keakuratan kuantitas produk saat terjadi request mutasi stok secara simultan (concurrent). |
| Portability | Konfigurasi database relasional lokal dan kode server dijalankan menggunakan Docker Compose. |

## 7. Constraints

- **Teknologi:** Go, PostgreSQL, GORM, Gin, Docker.
- **Tujuan Belajar:** Memahami implementasi SQL Transactions (Begin, Commit, Rollback) di Go, memodelkan relasi database One-to-Many dan Many-to-One di GORM, serta memahami implementasi paginasi kueri dan file streaming (CSV generator).

## 8. Success criteria

- API transaksi stok teruji secara fungsional (FR-5, FR-6, FR-7) dan tidak membiarkan kuantitas stok produk bernilai negatif.
- Laporan CSV dapat diunduh langsung via browser atau HTTP client dengan header `Content-Type: text/csv`.
- Memenuhi 12 dokumen standar proyek yang diisi secara spesifik.

## 9. Open questions

- **Transaction Management di Clean Architecture:** Ditangani di Repository Layer dengan menyediakan Transaction Wrapper/Callback agar tidak membocorkan implementasi GORM ke Service Layer (didokumentasikan di [ADR-001](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/adr/001-transaction-management.md)).
- **ORM vs sqlx:** Memilih untuk tetap menggunakan GORM karena kompleksitas *preloading* relasi Category dan Supplier di Project 2 lebih efisien diselesaikan dengan GORM, sementara integritas transaksi dapat dipelajari secara eksplisit melalui blok transaksi GORM.

---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
