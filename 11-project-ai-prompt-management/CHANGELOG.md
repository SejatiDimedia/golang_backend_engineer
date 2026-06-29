# Changelog: AI Prompt Management API

Semua perubahan penting pada proyek AI Prompt Management API akan dicatat di dokumen ini.

---

## [1.0.0] - 2026-06-29

### Added
- **Asymmetric Offline JWT verification:** Dashboard admin memvalidasi login JWT secara offline menggunakan berkas `public.key` milik Auth Service (Project 7).
- **Secure API Key System:** Otentikasi server-to-server API Key (`prompt_live_...`) dengan pengamanan hash SHA-256 di PostgreSQL.
- **Redis Cache-Aside Middleware:** Caching hasil verifikasi API Key (`apikey:<hash> -> workspace_id`) dengan TTL 1 jam untuk latensi $<2\text{ms}$.
- **Prompt Version Snapshots:** Version-control immutable berbasis full snapshot (v1, v2, v3) untuk audit trail instan.
- **Template Compiler Engine:** Parser regex untuk dynamic placeholders `{{var}}` dan token estimation formula.
- **Async Analytics Logging Daemon:** Pencatatan log latensi dan estimasi token secara asinkron menggunakan buffered channels dan background worker.
- **SQLite In-Memory Test Suites:** Unit testing komprehensif menguji middleware caching, RS256 offline verifier, compiler engine, dan multi-tenant workspaces.
- **Dockerfile & Compose:** Docker multi-stage build dan compose PostgreSQL 15 + Redis 7.
