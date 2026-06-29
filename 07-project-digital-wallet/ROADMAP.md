# Roadmap: Digital Wallet API

**Status:** `Planning`

Dokumen ini memetakan urutan pembangunan fitur di Digital Wallet API dengan fokus utama pada konsistensi data finansial double-entry ledger, distributed locking Redis, dan proteksi idempotensi.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Go module, docker-compose PostgreSQL 15 + Redis Alpine, config load, utilitas JWT helper, dan inisialisasi koneksi redis client. | — | `Planned` |
| 2 — Redis Lock Utility | Implementasi lock manager kustom (`AcquireLock`/`ReleaseLock`) menggunakan perintah `SET NX PX` dan Lua script pelepasan kunci. | Phase 1 | `Planned` |
| 3 — User & Wallet Domain | Entitas `User` dan `Wallet` (dengan `balance` kolom), registrasi akun otomatis menginisialisasi wallet dengan saldo 0. | Phase 2 | `Planned` |
| 4 — Ledger Transactions | Entitas `Transaction` (Double-entry debit/kredit ledger), REST API `POST /top-up` dan `POST /withdraw` terintegrasi transaksional atomic. | Phase 3 | `Planned` |
| 5 — Concurrency-Safe Transfer | REST API `POST /transfer` yang mengunci ID wallet pengirim dan penerima secara terdistribusi (Redis Lock) sebelum memotong/menambah saldo secara atomic. | Phase 4 | `Planned` |
| 6 — Idempotency Key Guard | Middleware/helper validasi `X-Idempotency-Key` menggunakan Redis cache TTL 1 jam untuk menangani request retries ganda dari client. | Phase 5 | `Planned` |
| 7 — Balance Cache & Invalidation | REST API `GET /wallet/balance` terintegrasi cache Redis. Penghapusan cache otomatis dipicu saat mutasi ledger baru ditulis ke database. | Phase 6 | `Planned` |
| 8 — Hardening & Testing | Simulasi race condition konkruen menggunakan Goroutines paralel untuk membuktikan sistem kebal double-spending, dan test idempotensi. | Phase 7 | `Planned` |
| 9 — Deployment & Docs | Dockerfile multi-stage, perbaruan root README, dan penyusunan 12 dokumen standar rilis. | Phase 8 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| Redis Client & Lock | — | Infrastruktur konkurensi dasar. Harus selesai sebelum logika transfer mulai ditulis. |
| User Register & Login | FR-1 | Otentikasi user dan pembuatan dompet digital otomatis. |
| Wallet Initialization | FR-2 | Setiap user terdaftar wajib mendapatkan alamat dompet unik untuk menerima/mengirim saldo. |
| Top-up & Withdrawal | FR-3 | Pintu masuk dan keluar saldo dompet. Bergantung pada skema ledger. |
| Mutasi Ledger (Double-Entry) | FR-5 | Catatan kebenaran data finansial. Ditulis transaksional bersama top-up/withdraw. |
| Transfer Saldo | FR-4 | Fitur transfer konkuren. Menggabungkan Redis Lock dan Transaksi SQL database. |
| Idempotency Key | FR-6 | Middleware pencegah request duplikasi transfer. Menyaring request di level HTTP routing. |
| Balance Cache | FR-7 | Membaca saldo dari cache Redis, dibersihkan saat ada transaksi sukses baru. |

## 3. Concepts this project is exercising

- **Redis Distributed Locking (Primitive):** Menjamin eksklusivitas pemrosesan data wallet konkuren.
- **Double-Entry Bookkeeping:** Menulis struktur data keuangan ter-audit di mana mutasi masuk (credit) dan keluar (debit) seimbang.
- **Idempotency API:** Membangun API handal yang toleran terhadap retry request jaringan tanpa duplikasi efek samping.
- **Cache Invalidation Pattern:** Memastikan data di cache Redis sinkron dengan database PostgreSQL (Write-Through/Invalidation).
- **Concurrency Testing di Go:** Menggunakan package `sync.WaitGroup` dan Goroutines di Go testing untuk menyimulasikan ratusan request paralel menyerang saldo wallet.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek Digital Wallet API |
