# PRD: Authentication Service (Auth Service)

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Dalam arsitektur berorientasi layanan (SOA) / microservices, implementasi otentikasi dan otorisasi secara terdesentralisasi (ad-hoc) di setiap service individu memicu duplikasi data user, meningkatkan risiko kebocoran credential, dan menyulitkan pemeliharaan kebijakan keamanan. Sistem memerlukan layanan Authentication Service (Identity Provider) terpusat yang menangani pendaftaran, login aman, rotasi token untuk memitigasi replay attacks, pemulihan sandi, serta otorisasi berbasis Role-Based Access Control (RBAC) secara konsisten yang dapat dikonsumsi oleh seluruh microservices lainnya.

## 2. Goals

- Menyediakan Identity Provider tunggal untuk pendaftaran, login, dan manajemen kredensial pengguna.
- Mengimplementasikan alur token berumur pendek (Access Token, 15 menit) didampingi token berumur panjang (Refresh Token, 7 hari) dengan mekanisme **Refresh Token Rotation (RTR)** guna mencegah penyalahgunaan token yang dicuri.
- Menyediakan alur **Email Verification** dan **Forgot/Reset Password** yang aman menggunakan secure validation tokens.
- Mengimplementasikan **Role-Based Access Control (RBAC)** dinamis di mana setiap user terasosiasi ke role yang memiliki hak akses (*permissions*) eksplisit.
- Menyediakan endpoint introspeksi token (`POST /auth/introspect`) agar microservices lain dapat memvalidasi token JWT secara terpusat tanpa harus mengakses database user secara langsung.

## 3. Non-goals

- **Integration with Real Third-Party OAuth providers (e.g., Google, Apple OAuth):** Kami akan mem-mock atau mensimulasikan OAuth callback flow untuk mempelajari state exchange, tetapi fokus utamanya adalah perancangan custom OAuth/OIDC client & auth flow lokal.
- **Real SMTP/Mail Services:** Pengiriman email verifikasi dan reset password akan di-mock dengan menulis tautan token ke stdout log/file guna menyederhanakan pengujian lokal.

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Client User | Melakukan registrasi, verifikasi email, login, reset password, dan memperbarui token sesi. | Setiap kali berinteraksi dengan platform |
| Downstream Microservices (e.g. Wallet, Booking) | Memanggil endpoint introspeksi untuk memverifikasi keabsahan JWT client dan mendapatkan role/permission pengguna bersangkutan. | Setiap request masuk |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | **Secure Register & Login:** Registrasi user baru dan autentikasi password menggunakan enkripsi `bcrypt`. | Must |
| FR-2 | **Token Rotation (RTR):** Mengembalikan Access Token (JWT) dan Refresh Token (UUID/random string) saat login. Setiap permintaan token baru menggunakan Refresh Token wajib me-rotate (menghapus refresh token lama dan mengembalikan refresh token baru). Jika refresh token yang sudah hangus (expired/revoked) digunakan kembali, seluruh sesi user tersebut harus otomatis dinonaktifkan (anti-replay attack). | Must |
| FR-3 | **Token Introspection:** API `POST /auth/introspect` menerima token JWT, memvalidasi tanda tangan, dan mengembalikan claims (`user_id`, `email`, `role`, `permissions`). | Must |
| FR-4 | **RBAC Engine:** Database PostgreSQL melacak relasi User, Roles, Permissions, dan UserRoles. | Must |
| FR-5 | **Email Verification:** Mengirimkan token verifikasi berumur pendek (15 menit) ke email saat registrasi. User tidak diizinkan login atau mengakses data tertentu sebelum status email berubah menjadi `verified`. | Must |
| FR-6 | **Forgot & Reset Password:** Menyediakan endpoint `/auth/forgot-password` (mengirim token reset) dan `/auth/reset-password` (mengganti sandi lama menggunakan token valid). | Must |
| FR-7 | **Token Revocation:** Endpoint `/auth/logout` secara eksplisit mencabut validitas refresh token di database. | Should |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Security | Password di-hash dengan `bcrypt` cost minimal 10. Access Token JWT harus ditandatangani dengan algoritma asymmetric key (RSA `RS256` / ECDSA) di mana Auth Service memegang private key untuk tanda tangan dan service lain hanya membutuhkan public key untuk verifikasi. |
| Performance | Endpoint `/auth/introspect` harus memiliki latency $<5\text{ms}$ (menggunakan Redis cache untuk menyimpan validitas token JWT aktif). |
| Resiliency | Jika database PostgreSQL mengalami down, microservice lain tetap dapat memverifikasi JWT secara offline selama public key RSA di-cache lokal di service masing-masing. |
| Data Consistency | Penanganan rotasi refresh token di DB harus menggunakan database transaction locking (`SELECT ... FOR UPDATE`) untuk mencegah race conditions saat refreshing token konkuren. |

## 7. Constraints

- **Teknologi:** Go, PostgreSQL (v15), Redis (v7), GORM, Gin, Docker & Docker Compose.
- **Asymmetric Encryption:** Menggunakan key pair RSA (2048-bit) untuk penandatanganan JWT.

## 8. Success criteria

- Alur Refresh Token Rotation berjalan mulus (token lama langsung tidak aktif, percobaan replay token memicu logout massal).
- Client verifikasi email dan reset password berjalan sukses secara asinkron.
- Microservices downstream (misal: client dummy) dapat memvalidasi token JWT secara offline menggunakan public key, atau secara online via introspect.
- Struktur RBAC membatasi hak akses berdasarkan izin permission secara dinamis.

## 9. Open questions

- **Asymmetric Signature Key Management:** Memilih **Asymmetric RS256**. Private key (`private.key`) dan public key (`public.key`) RSA 2048-bit akan otomatis dibuat saat bootstrapping awal di folder `config/certs` jika file belum ada. Downstream services memverifikasi token offline menggunakan public key secara lokal.
- **RBAC Policy Model:** Memilih **PostgreSQL Tables**. Model RBAC diimplementasikan secara native melalui skema database relasional (`users`, `roles`, `permissions`, `role_permissions`, `user_roles`) menggunakan GORM, guna menghindari kompleksitas engine eksternal seperti Casbin.


---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
