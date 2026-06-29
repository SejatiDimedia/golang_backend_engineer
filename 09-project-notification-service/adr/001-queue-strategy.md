# ADR-001: Pilihan Strategi Broker Antrean Pesan

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Layanan Notification Service membutuhkan broker antrean pesan (message queue broker) untuk menyimpan tugas pengiriman secara sementara dan memprosesnya secara asinkron. Kami harus memilih apakah akan menulis pengelola antrean (queue manager) secara manual menggunakan primitif Redis atau memanfaatkan pustaka siap pakai seperti `Asynq`.

## Decision

Kami memutuskan menggunakan **Hand-Rolled Redis List & Sorted Set Queue Manager** yang ditulis manual di Go.

Rincian antrean:
1. **Instant Queue:** Menggunakan struktur data Redis **List** dengan perintah `LPUSH` untuk memasukkan tugas dan worker menggunakan `BRPOP` (Blocking Remove Pop) untuk mengambil tugas secara real-time tanpa polling berulang.
2. **Scheduled & Retry Queue:** Menggunakan struktur data Redis **Sorted Set (ZSET)** dengan perintah `ZADD` di mana skornya adalah timestamp UTC kapan notifikasi harus diproses. Worker poller berkala mengambil data dengan `ZRANGEBYSCORE` dan memindahkannya ke List instan setelah waktu jatuh tempo tiba.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Hand-Rolled Redis (Chosen)** | - Memberikan pemahaman mendalam tentang siklus asinkron worker pool di Go.<br>- Mengajari penggunaan perintah Redis berkinerja tinggi (`BRPOP`, `ZADD`, `ZREMRANGEBYSCORE`).<br>- Bebas dependensi pustaka luar yang berat. | - Harus menangani serialisasi payload, error handling, dan daemon polling secara manual. |
| **B. Pustaka Asynq** | - Sangat matang, kaya fitur (auto-retry, monitoring dashboard, concurrency limit). | - Mengabstraksi seluruh proses antrean sehingga mengurangi nilai pembelajaran arsitektur tingkat rendah. |

## Reasoning

Notification Service dirancang di dalam roadmap belajar ini untuk memberikan pemahaman mendalam tentang arsitektur infrastruktur asinkron. Dengan membangun pengelola antrean manual (Opsi A), kami dilatih merancang background daemon loop Go, mengatur sinkronisasi goroutine secara aman, dan mengelola Redis commands secara langsung. Ini akan memberikan pondasi yang jauh lebih kokoh sebelum beralih ke library production-grade di industri.

## Consequences

- **Positif:** Mengerti cara kerja engine antrean asinkron, performa sangat cepat, kontrol penuh atas log audit.
- **Negatif:** Penulisan kode asinkron worker loop dan scheduler poller membutuhkan penanganan edge-case concurrency Go secara teliti (seperti penanganan OS signals dan context cancellation).
