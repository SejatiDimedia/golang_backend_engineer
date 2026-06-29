# Roadmap: Notification Service

**Status:** `Planning`

Dokumen ini memetakan urutan pembangunan fitur pada Notification Service untuk melatih pemrosesan antrean pesan asinkron, worker pool, scheduler, dan kebijakan retry exponential backoff.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Go module, docker-compose PostgreSQL 15 + Redis 7, config load, dan JWT helper. | — | `Planned` |
| 2 — Redis Queue Manager | Utilitas antrean asinkron manual (`internal/queue/manager.go`) menggunakan primitive List (`BRPOP`/`LPUSH`) dan Sorted Set (`ZADD`). | Phase 1 | `Planned` |
| 3 — DB Schema & Domain | Entitas `Notification` & `NotificationLog` (PostgreSQL), login/register JWT ad-hoc, dan repository metadata. | Phase 2 | `Planned` |
| 4 — REST API Gateway | REST API `POST /notifications` (instan & terjadwal), validasi payload (MIME type target email/phone/webhook), JWT Auth. | Phase 3 | `Planned` |
| 5 — Worker Pool Daemon | Daemon background loop (`internal/worker/pool.go`) dengan konfigurasi concurrency pool paralel membaca dari Redis blocking queue. | Phase 4 | `Planned` |
| 6 — Scheduler Poller Ticker | Daemon poller interval 1 detik (`internal/worker/poller.go`) memantau Sorted Set, memindahkan tugas jatuh tempo ke antrean utama. | Phase 5 | `Planned` |
| 7 — Exponential Backoff Retry | Implementasi backoff matematika $2^{\text{attempt}} \times 2$ detik, otomatis requeueing ke Sorted Set, audit log PostgreSQL, batas 5 kali max retry (DLQ). | Phase 6 | `Planned` |
| 8 — Hardening & Testing | Unit test queue manager, pengujian backoff, simulasi physical provider failure rates (30% error simulation). | Phase 7 | `Planned` |
| 9 — Deployment & Revisit | Dockerfile multi-stage, perbaruan status root README, penulisan 12 dokumen, dan draf revisi Project 3. | Phase 8 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| Redis Queue Client | — | Komponen inti antrean. Tanpa ini, worker dan API tidak bisa bertukar tugas. |
| Instant Enqueue | FR-2 | Endpoint API melempar tugas ke Redis List secara instan. |
| Background Worker Pool | FR-3 | Utas paralel yang membongkar antrean dan mengirimkan data ke provider dummy. |
| Scheduled Ticker | FR-6 | Memfasilitasi delay pengiriman berkas hingga waktu `send_at`. |
| Exponential Backoff | FR-4 | Kebijakan cerdas untuk menunda retry agar tidak membebani server/provider. |
| Audit Trail Logging | FR-5 | Rekaman lengkap histori status dan error per notifikasi di PostgreSQL. |
| DLQ (Dead-Letter Queue) | FR-7 | Mengamankan antrean dari looping error tak terbatas pada notifikasi rusak. |

## 3. Concepts this project is exercising

- **Asynchronous Message Queueing:** Memahami pola pub-sub dan produsen-konsumen di luar utas HTTP utama.
- **Worker Pools Concurrency:** Membatasi jumlah utas background Go (`go worker(id)`) untuk mencegah over-resource.
- **Sorted Set scheduling:** Mengurutkan waktu tunda di Redis untuk eksekusi terjadwal yang akurat.
- **Distributed Compensating Writes:** Mengatur agar state di Redis dan PostgreSQL tetap sinkron meskipun salah satu servis crash.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek Notification Service |
