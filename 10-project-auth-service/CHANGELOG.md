# Changelog: Authentication Service

Semua perubahan penting pada proyek Authentication Service akan dicatat di dokumen ini.

---

## [1.0.0] - 2026-06-29

### Added
- **Asymmetric RS256 Token Signing:** Implementasi keypair RSA 2048-bit auto-generator untuk token JWT yang mendukung verifikasi offline downstream services.
- **Refresh Token Rotation (RTR):** Mekanisme rotasi refresh token otomatis di database relasional untuk mengamankan sesi platform.
- **RTR Anti-Replay Safeguard:** Database row locking `SELECT ... FOR UPDATE` dan dynamic mass revocation untuk menangani deteksi serangan replay.
- **Native Relational RBAC Engine:** Modeling tabel many-to-many (`users`, `roles`, `permissions`, `role_permissions`, `user_roles`) di GORM PostgreSQL untuk authorization terpusat.
- **Lifecycle Activation & Reset Flows:** Endpoint verifikasi email aktivasi dan pemulihan forgot/reset password.
- **Token Introspection REST API:** Endpoint `/auth/introspect` dengan Redis caching untuk performa tinggi.
- **Suite Unit Testing:** Pengujian database SQLite in-memory, replay simulation, dan offline RSA signing tests.
- **Dockerfile & Compose:** Docker multi-stage build dan compose PostgreSQL 15 + Redis 7.
