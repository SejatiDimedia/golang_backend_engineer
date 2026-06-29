# Roadmap: AI Prompt Management API

**Status:** `Planning`

Dokumen ini memetakan urutan pembangunan fitur pada AI Prompt Management Service untuk melatih modular multi-tenancy, version-control prompt snapshots, otentikasi API Key terdistribusi, dan offline JWT verification.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Go module, docker-compose PG + Redis, config load, and public key cert setup. | — | `Planned` |
| 2 — DB Schema & Models | GORM modeling `Workspace`, `Prompt`, `PromptVersion`, `ApiKey`, dan `AnalyticsLog`. | Phase 1 | `Planned` |
| 3 — Offline JWT Verifier | Integrasi public key RSA dari Auth Service untuk otentikasi dashboard admin secara offline. | Phase 2 | `Planned` |
| 4 — API Key Generator | Utilitas pembuat API key format `sk_live_...` dengan penyimpanan hash SHA-256. | Phase 3 | `Planned` |
| 5 — Workspace & Prompt CRUD | REST API management workspaces dan prompt versions (Draft vs Active). | Phase 4 | `Planned` |
| 6 — Template Compiler Engine | Regex parser untuk mengompilasi string template dinamis berbasis variabel `{{var}}`. | Phase 5 | `Planned` |
| 7 — API Key Middleware & Cache | Filter otentikasi server-to-server didukung Redis caching `<hash>` untuk performa tinggi. | Phase 6 | `Planned` |
| 8 — Usage Analytics | Logging otomatis latensi, hit count, dan estimasi token length pasca kompilasi. | Phase 7 | `Planned` |
| 9 — Hardening & Testing | Unit test compiler regex, testing API Key cache hit/miss, dan SQLite in-memory tests. | Phase 8 | `Planned` |
| 10 — Deployment & Capstone | Dockerfile multi-stage, penulisan 12 dokumen, dan checklist kesiapan integrasi capstone. | Phase 9 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| Asymmetric JWT Verifier | FR-6 | Gerbang keamanan akses dashboard admin untuk mengelola workspaces. |
| API Key Hashing | FR-5 | Autentikasi server-to-server yang aman dari pencurian database. |
| Prompt Version Snapshots | FR-3 | Menyimpan salinan teks penuh agar tidak terjadi kegagalan runtime saat prompt diubah. |
| Regex Parser Compiler | FR-4 | Mesin rendering pengisi nilai variabel dinamis. |
| Usage Logger | FR-7 | Mengukur pemakaian API key secara analitik. |

## 3. Concepts this project is exercising

- **Offline Token Authentication:** Memanfaatkan public key RSA untuk validasi JWT tanpa panggilan database/jaringan.
- **API Key Security:** Menerapkan hash SHA-256 dan prefix token format industri.
- **Domain Versioning (Snapshots):** Menerapkan konsep immutability pada domain data model.
- **Redis Cache-Aside Pattern:** Mengoptimalkan throughput validasi otentikasi.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek AI Prompt Management API |
