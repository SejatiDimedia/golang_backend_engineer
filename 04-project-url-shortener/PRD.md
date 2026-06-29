# PRD: URL Shortener Service

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Ketika membagikan tautan web (URL) yang panjang, sering kali tautan tersebut menjadi sulit dibaca, sulit diketik secara manual, dan rentan rusak karena pemotongan karakter saat dikirim. Selain itu, pemilik tautan tidak memiliki cara mudah untuk melacak berapa banyak orang yang telah mengklik tautan tersebut atau membatasi masa berlaku tautan tersebut secara otomatis.

## 2. Goals

- Memungkinkan pengguna untuk memperpendek URL panjang menjadi URL pendek yang acak atau menggunakan alias kustom.
- Menyediakan pelacakan dasar (jumlah klik) dan pengaturan masa berlaku (kedaluwarsa) untuk setiap tautan pendek.
- Menjadi sarana pembelajaran pertama mengenai ekosistem backend Go (REST API, routing, koneksi database relasional, dan dockerization).

## 3. Non-goals

- **Autentikasi Pengguna (User Auth):** Tidak ada registrasi, login, atau pembatasan URL per user untuk versi awal ini (disimpan untuk Project 3 & 7). Semua fitur dapat diakses secara publik.
- **Tampilan Antarmuka (UI Frontend):** Proyek ini hanya menyediakan REST API. Tidak ada halaman web HTML/CSS yang kompleks.
- **Skalabilitas Multi-Region atau High Performance Caching:** Caching (seperti Redis) sengaja ditiadakan di versi awal agar fokus pada CRUD PostgreSQL dan GORM/sqlx.

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Public User | Memperpendek tautan panjang dengan cepat dan membagikannya secara mudah | Sesekali (ad-hoc) |
| Content Creator | Membuat alias tautan yang mudah diingat (misal `/kado-timur`) dan melacak total klik | Sering |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | Pengguna dapat mengirimkan URL panjang dan mendapatkan short code acak (6-8 karakter). | Must |
| FR-2 | Pengguna dapat mengajukan kustom alias (misal `/custom-name`). Sistem harus memvalidasi keunikan alias tersebut. | Must |
| FR-3 | Saat mengakses `/r/{short_code}`, sistem melakukan pengalihan (redirect HTTP 302/301) ke URL panjang target. | Must |
| FR-4 | Pengguna dapat mengatur waktu kedaluwarsa (`expires_at`) opsional saat pembuatan. URL pendek tidak dapat diakses setelah melewati waktu tersebut. | Should |
| FR-5 | Setiap kali pengalihan sukses dilakukan, sistem menambah hitungan klik (`click_count`) pada tautan terkait secara real-time. | Must |
| FR-6 | Pengguna dapat melihat metadata URL (kapan dibuat, expires_at, target_url, dan click_count) via endpoint statistik. | Must |
| FR-7 | Layanan memiliki endpoint `GET /health` untuk memantau status aplikasi dan koneksi ke database. | Must |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Performance | Waktu respon pengalihan (redirect) < 100ms dalam kondisi normal. |
| Security | Validasi input URL target (harus berformat URL valid) untuk mencegah eksploitasi keamanan dasar. |
| Availability | Ketersediaan mandiri (self-contained) di lingkungan lokal menggunakan Docker Compose. |
| Scalability | Skala terbatas untuk single-instance deployment lokal. PostgreSQL menjadi bottleneck pertama saat beban concurrent tinggi. |
| Data consistency | Data statistik klik dan masa berlaku harus disimpan secara presisten di PostgreSQL. |

## 7. Constraints

- **Teknologi:** Harus menggunakan bahasa pemrograman Go dengan framework REST Gin, PostgreSQL sebagai database relasional, dan Docker untuk containerization.
- **Tujuan Belajar:** Merupakan proyek pertama (Phase 1), sehingga fokus pada penulisan struktur Clean Architecture dasar (Handler -> Service -> Repository) dan implementasi interface serta penanganan error Go yang idiomatik.

## 8. Success criteria

- REST API berjalan dengan sukses menggunakan `docker-compose up`.
- Seluruh kebutuhan fungsional (FR-1 sampai FR-7) teruji dengan baik dan lulus pengujian manual/otomatis.
- Terpenuhinya 12 dokumen standar proyek (README, PRD, ARCHITECTURE, dll.) di folder proyek.

## 9. Open questions

- **GORM vs sqlx:** Memilih GORM karena kemudahan migrasi otomatis dan *speed-to-ship* untuk proyek pertama ini (didokumentasikan di [ADR-001](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/adr/001-orm-decision.md)).
- **Short code generation strategy:** Menggunakan URL-safe Base64 encoding dari timestamp di sisi aplikasi Go agar independen dari *database lock* (didokumentasikan di [ADR-002](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/adr/002-short-code-generation.md)).

---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
