# Changelog: Booking Management System

Format penulisan didasarkan pada [Keep a Changelog](https://keepachangelog.com/).

---

## [1.0.0] — 2026-06-29 — Initial Release

### Added
- **JWT Authentication Helper:** Pembuatan utility JWT generator dan validator (`internal/utils/jwt.go`).
- **HTTP Middlewares:**
  - `AuthMiddleware` untuk parsing JWT Bearer token dan validasi masa aktif token.
  - `RequireRole` guard middleware untuk membatasi hak akses resource berdasarkan hak role (customer vs admin).
- **User Domain:** Entitas `User` dengan hashing password aman berbasis `bcrypt`. REST API `/register` dan `/login`.
- **Desk (Aset) Catalog:** CRUD endpoint master data meja/ruangan kerja coworking space (`/desks`), terproteksi role admin.
- **Booking Core & Transactions:**
  - Entitas `Booking` terelasi ke User dan Desk.
  - Transaksi pembuatan booking transaksional (`POST /bookings`) terproteksi locking `SELECT ... FOR UPDATE` baris meja bersangkutan.
  - Algoritma pencarian overlap booking di database menggunakan perbandingan UTC: `start_time < endTime AND end_time > startTime`.
- **Cancellation Window Guard:** Endpoint `/bookings/:id/cancel` yang memvalidasi pembatalan minimal 2 jam sebelum pemesanan dimulai (dengan kemampuan bypass khusus admin).
- **Log Stub Notification:** Output stdout logger simulasi notifikasi saat booking sukses dibuat atau dibatalkan.
- **Unit Testing:** Suite pengujian unit lengkap untuk verifikasi error overlap, bypass pembatalan admin, dan token middleware.
- **Dockerization:** Berkas Dockerfile multi-stage dan docker-compose.yml.
