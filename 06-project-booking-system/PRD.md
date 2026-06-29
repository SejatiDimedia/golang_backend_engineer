# PRD: Booking Management System (Coworking Space)

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Pengguna coworking space sering kali mengalami kesulitan memesan ruangan rapat atau meja kerja (*desk*) secara terencana. Di sisi pengelola, bentrokan waktu pemesanan (double-booking) sering terjadi karena sistem pemesanan tidak memvalidasi slot waktu secara real-time dan aman secara konkurensi. Selain itu, penanganan zona waktu (*timezone*) yang tidak konsisten antara aplikasi frontend (client) dan database backend sering mengakibatkan pergeseran jam pemesanan yang merugikan pengguna.

## 2. Goals

- Memungkinkan pengguna mendaftar dan masuk (login) untuk mengelola pemesanan mereka sendiri.
- Menyediakan katalog Ruangan/Meja yang dapat dipesan berdasarkan rentang waktu tertentu.
- Mencegah bentrokan waktu pemesanan (*double-booking*) pada ruangan/meja yang sama untuk rentang waktu yang tumpang tindih.
- Mengirimkan simulasi notifikasi konfirmasi (stub) ketika pemesanan berhasil dibuat atau dibatalkan.
- Memperkenalkan penanganan otentikasi JWT ad-hoc di Go, middleware HTTP, serta penanganan tipe data `time.Time` secara konsisten dalam zona UTC.

## 3. Non-goals

- **Sistem Pembayaran Terintegrasi:** Tidak ada integrasi dengan payment gateway (seperti Midtrans/Stripe) di rilis ini. Seluruh pemesanan dianggap disetujui secara instan atau bayar di tempat (*pay at location*).
- **Layanan Notifikasi Riil:** Pengiriman email/SMS sungguhan ditiadakan di rilis ini. Konfirmasi pemesanan hanya memicu output log konsol sederhana (stub) sebagai placeholder untuk *Notification Service* (Project 6).
- **Shared Auth Service:** Manajemen otentikasi JWT ditulis langsung di dalam sasis proyek ini secara ad-hoc (tidak menggunakan auth service eksternal, untuk merasakan batasan/pain point duplikasi kode auth sebelum masuk ke Project 7).

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Customer (Member) | Mencari meja/ruangan kosong pada tanggal tertentu, melakukan booking, dan melihat riwayat pesanan mereka. | Beberapa kali seminggu |
| Space Admin | Mengelola ketersediaan ruangan, memantau seluruh jadwal booking harian, dan melakukan pembatalan booking jika diperlukan. | Setiap hari |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | Pengguna dapat melakukan registrasi akun baru dan login menggunakan email & password. | Must |
| FR-2 | Sistem menghasilkan JSON Web Token (JWT) yang valid saat login sukses. Token tersebut wajib dikirim di header HTTP Auth untuk mengakses resource terproteksi. | Must |
| FR-3 | Pengguna dapat melihat daftar Ruangan/Meja yang aktif beserta status ketersediaannya. | Must |
| FR-4 | Pengguna dapat membuat pemesanan (`Booking`) baru dengan menentukan ID Ruangan/Meja, tanggal, waktu mulai (`start_time`), dan waktu selesai (`end_time`). | Must |
| FR-5 | **Pencegahan Double-Booking:** Sistem wajib menolak pemesanan baru jika slot waktu ruangan/meja yang diminta tumpang tindih (*overlap*) dengan pemesanan aktif lainnya. | Must |
| FR-6 | Pengguna dapat melihat daftar riwayat pemesanan mereka sendiri, sedangkan Admin dapat melihat seluruh daftar pemesanan di coworking space. | Must |
| FR-7 | Pengguna dapat membatalkan pemesanan mereka sendiri (`Cancel Booking`) maksimal 2 jam sebelum waktu mulai pemesanan. | Should |
| FR-8 | Sistem mencetak log simulasi notifikasi (stub) ke stdout setiap kali pemesanan sukses dibuat atau dibatalkan. | Must |
| FR-9 | Menyediakan endpoint `GET /health` untuk memantau status server dan koneksi database. | Must |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Timezone Handling | Semua data tanggal & waktu disimpan di PostgreSQL dalam format **UTC**. API menerima format standar RFC3339 dengan offset zona waktu dan mengonversinya secara aman ke UTC sebelum disimpan. |
| Security | Password pengguna wajib di-hash menggunakan algoritma **bcrypt** sebelum disimpan ke database. |
| Concurrency | Validasi bentrokan slot waktu pemesanan harus dilindungi dari kondisi balapan konkurensi (misal dua user memesan slot yang sama di milidetik yang sama) menggunakan locking transaksional atau unique constraints. |
| Portability | Berjalan di lingkungan kontainer terisolasi menggunakan Docker Compose (Go app + PostgreSQL). |

## 7. Constraints

- **Teknologi:** Go, PostgreSQL, GORM, Gin, JWT (golang-jwt/jwt), bcrypt, Docker.
- **Domain:** Coworking Space Booking System. (Bisa disesuaikan jika pengguna memilih domain lain seperti Barbershop, Clinic, atau Photography Studio).
- **Tujuan Belajar:** Memahami ad-hoc JWT auth, rantai HTTP middleware Gin, penanganan timezone di backend, algoritma pencocokan rentang waktu overlap, dan pembuatan kode stub.

## 8. Success criteria

- API JWT Authentication berhasil mengamankan endpoint pemesanan (hanya pemilik token valid yang bisa memesan).
- Bentrokan waktu pemesanan terbukti ditolak oleh sistem dan tidak menghasilkan tumpang tindih slot di database.
- Database menyimpan waktu mulai dan selesai dalam format UTC secara konsisten.

## 9. Open questions

- **Domain Selection:** Memilih domain **Coworking Space** karena pemodelan meja/ruangan kerja sangat representatif untuk simulasi bentrokan slot waktu yang dinamis.
- **Pencegahan Overlap Booking di Database:** Menggunakan query transaksional `SELECT ... FOR UPDATE` (pessimistic lock) untuk mengambil data booking aktif dan memvalidasi irisan waktu (*overlap checking*) di tingkat Go. Hal ini menjaga kode tetap portabel serta melatih penulisan logika waktu di Go secara eksplisit (didokumentasikan di [ADR-002](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/06-project-booking-system/adr/002-overlap-checking-strategy.md)).

---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
