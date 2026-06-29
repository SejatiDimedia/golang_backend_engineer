# Roadmap: Booking Management System (Coworking Space)

**Status:** `Planning`

Dokumen ini mengurutkan langkah pengerjaan fitur dalam proyek Booking Management System untuk memastikan otentikasi JWT ad-hoc dan pencegahan double-booking konkruen terimplementasi secara aman.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Go module, docker-compose PostgreSQL 15, config load, utilitas JWT helper (Generate/Validate), dan bcrypt password hashing. | — | `Planned` |
| 2 — User Domain & Auth API | Entitas `User` (GORM), REST API Registrasi dan Login, dan validasi input email/password. | Phase 1 | `Planned` |
| 3 — Auth Middleware | Gin HTTP Middleware (`AuthMiddleware`) untuk memverifikasi token JWT di Authorization header dan menyuntikkan data user aktif ke request context. | Phase 2 | `Planned` |
| 4 — Catalog Management | REST API CRUD untuk entitas Meja/Ruangan (`Desk` / `Room`) sebagai aset coworking space yang dapat dipesan. | Phase 3 | `Planned` |
| 5 — Booking Core (Transactions) | Entitas `Booking`, REST API pembuatan booking dengan pengubahan timezone ke UTC, locking transaksional, dan algoritma pendeteksi tumpang tindih waktu (overlap). | Phase 4 | `Planned` |
| 6 — Booking Cancellation | Logika pembatalan booking (`POST /bookings/:id/cancel`) dengan batasan waktu minimal 2 jam sebelum pemesanan dimulai. | Phase 5 | `Planned` |
| 7 — Notifications & Stats | Log simulasi notifikasi (stub) ke stdout saat booking dibuat/dibatalkan, dan API dashboard admin untuk mereview total pemesanan. | Phase 6 | `Planned` |
| 8 — Hardening & Testing | Unit test algoritma overlap checker, unit test JWT middleware menggunakan router recorder, dan validasi input datetime RFC3339. | Phase 7 | `Planned` |
| 9 — Deployment & Docs | Dockerfile multi-stage, perbaruan root README, dan pemenuhan 12 berkas dokumentasi rilis proyek. | Phase 8 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| User Register & Login | FR-1 | User akun harus ada terlebih dahulu agar token JWT dapat dihasilkan untuk menguji otentikasi endpoint lainnya. |
| JWT Authentication | FR-2 | Middleware keamanan utama. Harus terimplementasi dan teruji sebelum endpoint pemesanan mulai ditulis. |
| Desk/Room Catalog | FR-3 | Katalog aset yang dapat dipesan. Harus ada agar ID aset valid saat melakukan pemesanan. |
| Booking Creation | FR-4 | Fitur inti. Bergantung pada data User (Auth) dan Aset (Catalog). |
| Overlap Validation | FR-5 | Penjaga integritas data. Logika overlap check disisipkan langsung di dalam alur pembuatan booking. |
| User/Admin Bookings List | FR-6 | Memungkinkan user melihat pemesanan mereka (dan admin melihat semua). Membaca data yang sukses dibuat di FR-4. |
| Cancellation | FR-7 | Fitur pembatalan. Menguji logika manipulasi status booking (`CANCELLED`) dan validasi rentang selisih jam. |
| Notification Log (Stub) | FR-8 | Stub log tercetak tepat setelah booking/cancellation selesai diproses di service layer. |
| Health Check | FR-9 | Memantau keaktifan service dan database relasional. |

## 3. Concepts this project is exercising

- **Ad-Hoc JWT Authentication:** Mengamankan API menggunakan token JWT yang didekode secara manual di Gin Middleware.
- **Waktu & Zona Waktu (Go time.Time):** Mengonversi input datetime lokal ke format UTC secara konsisten sebelum disimpan, serta menggunakan time helper untuk menghitung selisih durasi.
- **Pessimistic Locking (FOR UPDATE):** Menerapkan lock pada baris aset (Desk/Room) untuk mengamankan kueri validasi ketersediaan waktu dari *race condition* konkruen.
- **Algoritma Overlap Rentang Waktu:** Menerapkan kueri deteksi irisan waktu $S_1 < E_2 \ \text{and} \ S_2 < E_1$.
- **Pola Kode Stub:** Membuat fungsionalitas placeholder yang terisolasi (notifikasi log) agar siap direnovasi di masa depan.

## 4. Known risks / unknowns at planning time

- **Kebocoran Waktu (Time Drift):** Jika server host dan database postgres berjalan pada zona waktu yang berbeda, query pembandingan waktu PostgreSQL bisa bergeser. Kami mengantisipasi ini dengan memaksakan parameter koneksi PostgreSQL menggunakan `TimeZone=UTC` dan memastikan aplikasi Go selalu melakukan `.UTC()` sebelum mengirim waktu ke database.
- **JWT Secret Management:** Secret key JWT disimpan di berkas `.env`. Jika secret bocor, keamanan sistem runtuh. Kita harus menegaskan penggunaan kunci yang rumit dan tidak mempublikasikannya ke Git.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek Booking Management System |
