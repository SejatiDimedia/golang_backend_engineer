# Changelog: Digital Wallet API

Format penulisan didasarkan pada [Keep a Changelog](https://keepachangelog.com/).

---

## [1.0.0] — 2026-06-29 — Initial Release

### Added
- **Hand-Rolled Redis Lock Manager:** Utilitas pengunci terdistribusi (`internal/utils/lock.go`) menggunakan command Redis primitive `SET NX PX` dan Lua script pelepasan lock atomik.
- **Double-Entry Ledger Design:** Skema database PostgreSQL yang memisahkan tabel `wallets` (untuk running balance numeric precision) dan tabel `transactions` (debit/kredit mutasi ledger) yang diproses secara atomic.
- **Idempotency Key Middleware:** Filter `X-Idempotency-Key` (`internal/middleware/idempotency.go`) untuk mencegat request retries ganda dengan menyimpan salinan body respons sukses di Redis cache (TTL 1 jam).
- **Balance Caching:** Pola cache-aside untuk saldo dompet digital di Redis dengan invalidasi cache otomatis setelah mutasi ledger ditulis.
- **User & Wallet Domain:** Registrasi user (`POST /register`) secara transaksional otomatis menginisialisasi dompet digital baru dengan nomor rekening unik (contoh: `W-10001`).
- **Transfer API:** API transfer antar dompet terproteksi distributed locking berurutan untuk menghindari deadlock.
- **Unit Testing:** Unit test komprehensif menguji kalkulasi saldo, dan simulasi transfer konkuren 10 Goroutines paralel.
- **Dockerization:** Berkas Dockerfile multi-stage dan docker-compose.yml (Postgres + Redis).
