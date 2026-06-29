# Roadmap: Authentication Service (Auth Service)

**Status:** `Planning`

Dokumen ini memetakan urutan pembangunan fitur pada Authentication Service untuk mengamankan platform menggunakan asymmetric JWT tokens, rotasi refresh token, verifikasi email, forgot password, dan RBAC terpusat.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Go module, docker-compose PostgreSQL + Redis, config load, dan RSA key generator. | — | `Planned` |
| 2 — DB Schema & Models | Modeling entitas `User`, `Role`, `Permission`, `RefreshToken`, `VerificationToken` di GORM. | Phase 1 | `Planned` |
| 3 — RSA Signer & RTR Logic | Helper token RS256 dan logic database transaction untuk Refresh Token Rotation (RTR). | Phase 2 | `Planned` |
| 4 — Auth REST APIs | Endpoint REST API `/auth/register`, `/auth/login`, `/auth/refresh`, `/auth/logout`. | Phase 3 | `Planned` |
| 5 — Verification & Reset | Alur verifikasi email dan lupa sandi (forgot/reset password) dengan mock email sender. | Phase 4 | `Planned` |
| 6 — Introspection API | REST API `/auth/introspect` untuk microservices lain didukung Redis caching token aktif. | Phase 5 | `Planned` |
| 7 — RBAC Management | API endpoint admin untuk mengelola dynamic roles & permissions. | Phase 6 | `Planned` |
| 8 — Hardening & Testing | Unit test key generation, RTR replay attacks detection, offline verification, dan RBAC policy validation. | Phase 7 | `Planned` |
| 9 — Deployment & Retrofit | Dockerfile, dokumen rilis lengkap, dan rencana migrasi retrofit Booking/Wallet/Notification. | Phase 8 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| RSA Key Generator | FR-3 | Kunci tanda tangan harus siap sebelum token JWT RS256 dapat dibuat. |
| User & Role Schemas | FR-4 | Model relasional database terpusat harus terdefinisi untuk registrasi awal. |
| RTR Replay Detection | FR-2 | Inti dari session management yang aman untuk mengantisipasi kebocoran token. |
| Token Introspection | FR-3 | Menyediakan gerbang validasi online/offline bagi service lain di platform. |
| Email & Password Tokens | FR-5, FR-6 | Alur lifecycle pengamanan kredensial dan aktivasi akun user. |

## 3. Concepts this project is exercising

- **Asymmetric Cryptography:** Menggunakan pasangan kunci privat/publik (RSA) untuk enkripsi JWT.
- **Refresh Token Rotation (RTR):** Melindungi sesi pengguna dari pencurian token dengan deteksi replay serangan asinkron.
- **Role-Based Access Control (RBAC):** Merancang arsitektur join table many-to-many untuk mengelola permissions secara dinamis.
- **Platform Microservices Introspection:** Merancang API modular agar dikonsumsi secara andal oleh downstream services.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek Authentication Service |
