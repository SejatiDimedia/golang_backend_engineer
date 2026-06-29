# Changelog: Notification Service

Semua perubahan penting pada proyek Notification Service akan dicatat di dokumen ini.

---

## [1.0.0] - 2026-06-29

### Added
- **Asynchronous Redis Queue Manager:** Implementasi antrean hand-rolled menggunakan Redis List (`LPUSH`/`BRPOP`) dan Sorted Set (`ZADD`).
- **Atomic Lua Scheduler:** Poller background 1 detik menggunakan script Lua di Redis untuk mencegah duplikasi polling pada konfigurasi multi-node.
- **Worker Pool Daemon:** Concurrency pool paralel di latar belakang untuk memproses tugas pengiriman secara non-blocking.
- **Exponential Backoff Retry Strategy:** Mekanisme penundaan retry otomatis ($2^{\text{attempt}} \times 2$ detik) saat pengiriman gagal dengan batas 5 kali max retries (Dead-Letter Queue).
- **PostgreSQL Audit Log:** Relasional database schemas untuk melacak status notifikasi (`PENDING`, `PROCESSING`, `SENT`, `FAILED`) dan audit log error detail per percobaan.
- **JWT Auth & API Gateways:** Route HTTP Gin dilindungi JWT ad-hoc, endpoint `/notifications` mengembalikan status `202 Accepted`.
- **Suite Unit Testing:** Pengujian miniredis in-memory, mock repository, dan pengetesan formula backoff math.
- **Docker Compose:** Konfigurasi docker postgres:15 dan redis:7.
