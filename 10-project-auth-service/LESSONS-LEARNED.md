# Lessons Learned: Authentication Service

Retrospektif pembelajaran dari arsitektur platform Authentication Service.

---

## 1. Keberhasilan Offline JWT Verification RS256
- **Temuan:** Di awal perancangan microservices platform, otentikasi terpusat sering kali memicu kekhawatiran karena membebani satu service Auth Service dengan jutaan kueri validasi token online.
- **Pembelajaran:** Menggunakan asimmetris RS256 JWT adalah solusi terbaik. Downstream services memvalidasi integritas token secara offline menggunakan kunci publik (`public.key`) secara lokal. Ini memotong latensi otentikasi dari milidetik jaringan menjadi mikrodetik kalkulasi CPU lokal ($O(1)$) dan memutus ketergantungan runtime.

## 2. Mengatasi Rollback Transaksi GORM untuk Menyimpan State Serangan (Replay Attack)
- **Temuan:** Awalnya, ketika replay attack dideteksi di dalam kueri database `RotateRefreshToken`, service layer langsung melempar error dan Transaction Manager (`WithTransaction`) me-rollback database. Namun, hal ini menyebabkan perintah pencabutan massal seluruh sesi aktif user (yang sengaja dilakukan untuk menghentikan akses peretas) ikut ter-rollback dan dibatalkan.
- **Pembelajaran:** Untuk menyimpan status darurat (seperti revoke massal) di database meskipun operasi bisnis utama gagal, kita harus mengembalikan `nil` (tanpa error) di dalam callback transaksi agar database men-commit status pencabutan ke disk, lalu di luar context transaksi barulah error otentikasi sesungguhnya dikembalikan ke client.

## 3. Pentingnya Penggunaan SQLite in-Memory untuk Unit Testing Relasional
- **Temuan:** Menulis mock sql manual (seperti sqlmock) untuk query JOIN yang kompleks pada RBAC tables many-to-many sangat rumit dan rentan patah saat ada perubahan schema kecil.
- **Pembelajaran:** Menggunakan driver SQLite in-memory (`sqlite.Open("file::memory:")`) terintegrasi GORM adalah strategi terbaik. Pengembang dapat menguji real SQL behaviors (JOIN many-to-many, row-level locking, foreign key constraints) secara instan dan tangguh tanpa perlu menyalakan kontainer PostgreSQL fisik saat unit testing dijalankan.
