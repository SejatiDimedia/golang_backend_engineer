# PRD: Digital Wallet API

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Aplikasi dompet digital (*digital wallet*) sangat sensitif terhadap konsistensi data finansial. Dua masalah paling kritis adalah **double-spending** (pengguna mengirim saldo melebihi kapasitas dompetnya akibat request paralel) dan **request duplikasi** (koneksi internet lambat memicu pengguna menekan tombol "Kirim" berkali-kali, menghasilkan transfer ganda). Sistem membutuhkan proteksi tingkat tinggi menggunakan kunci idempotensi, distributed lock, serta pencatatan akuntansi double-entry ledger yang ketat.

## 2. Goals

- Memungkinkan pengguna mendaftar dan masuk (login JWT ad-hoc) dan secara otomatis mendapatkan satu dompet digital (`Wallet`).
- Menyediakan API untuk melakukan isi saldo (`Top-up`), penarikan (`Withdrawal`), dan transfer saldo antar pengguna (`Transfer`).
- Menerapkan **Double-Entry Ledger** (tabel mutasi debit/kredit) untuk menjamin akurasi saldo finansial.
- Mencegah double-spending konkuren menggunakan **Redis Distributed Locks**.
- Mencegah request ganda menggunakan **Idempotency Keys** yang disimpan sementara di Redis cache.
- Menggunakan Redis cache untuk menyajikan data saldo dompet secara instan (`GET /wallet/balance`).

## 3. Non-goals

- **Koneksi Payment Gateway Riil:** Seluruh top-up dan withdrawal dilakukan secara simulasi API instan (tidak menggunakan bank transfer/Qris sungguhan).
- **Multi-Currency:** Dompet digital ini hanya melayani mata uang tunggal (Rupiah/IDR).

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Wallet User | Melakukan pengisian saldo, mengirim uang ke sesama user dengan aman, dan melihat riwayat mutasi rekening. | Beberapa kali sehari |
| System Auditor | Memeriksa seluruh catatan debit/kredit sistem untuk memastikan balance kliring bernilai nol. | Bulanan |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | Pengguna dapat melakukan registrasi dan masuk (login JWT) untuk mengamankan wallet. | Must |
| FR-2 | Setiap user otomatis memiliki satu `Wallet` dengan nomor akun dompet unik. | Must |
| FR-3 | Pengguna dapat melakukan `Top-up` dan `Withdraw` saldo wallet mereka sendiri. | Must |
| FR-4 | **Transfer Saldo:** Pengguna dapat mengirim saldo ke pengguna lain berdasarkan nomor akun dompet target. | Must |
| FR-5 | **Double-Entry Ledger:** Setiap perubahan saldo (top-up, withdraw, transfer) wajib dicatat sebagai entri jurnal debit/kredit yang seimbang di database. | Must |
| FR-6 | **Idempotency Guard:** API transfer dan top-up wajib menyertakan header `X-Idempotency-Key`. Jika kunci yang sama dikirim kembali dalam waktu 1 jam, server harus mengembalikan respons yang sama tanpa mengeksekusi ulang transaksi. | Must |
| FR-7 | Pengguna dapat melihat daftar histori mutasi saldo mereka. | Must |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Financial Consistency | Nilai saldo tidak boleh disimpan sebagai kolom mentah yang bebas di-update begitu saja (`UPDATE wallets SET balance = balance + X`). Saldo harus dikalkulasikan atau dikonfirmasi berdasarkan hasil penjumlahan entri mutasi ledger historis untuk mencegah kebocoran saldo. |
| Concurrency (Redis Lock) | Sebelum mengeksekusi transfer saldo, sistem wajib mengunci ID wallet pengirim dan penerima secara terdistribusi di Redis untuk memblokir transaksi transfer paralel. |
| Idempotency Storage | Idempotency key beserta salinan respons HTTP-nya disimpan di Redis cache dengan masa kedaluwarsa (TTL) selama 1 jam. |
| Caching | Kueri saldo wallet (`GET /wallet/balance`) diarahkan ke Redis cache. Cache otomatis dihapus (*invalidated*) setiap kali mutasi transaksi baru sukses disimpan ke database PostgreSQL. |

## 7. Constraints

- **Teknologi:** Go, PostgreSQL, Redis, GORM, Gin, Docker.
- **Library Lock terdistribusi:** Redis distributed locks (menggunakan pustaka `bsm/redislock` atau `go-redsync/redsync`).
- **Idempotency:** Custom Gin middleware atau service wrapper untuk memeriksa dan menyimpan idempotency key.

## 8. Success criteria

- Lulus pengujian konkurensi (simulasi race condition transfer paralel) tanpa menghasilkan saldo minus atau nilai balance yang salah.
- Pengiriman request transfer dengan idempotency key yang sama berulang kali menghasilkan respons yang identik tanpa memotong saldo ganda.
- Kecepatan respons saldo instan karena ter-cache di Redis.

## 9. Open questions

- **Ledger Model:** Memilih **Opsi A (Running Balance + Ledger)**. Sistem menyimpan running balance di kolom `balance` tabel `wallets` untuk performa kueri baca, namun dilindungi oleh transaksi debit/kredit yang tercatat secara atomis di tabel ledger (didokumentasikan di [ADR-001](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/07-project-digital-wallet/adr/001-ledger-model-strategy.md)).
- **Redis Lock Client:** Memilih **Opsi B (Hand-Rolled Redis Lock)**. Kita akan mengimplementasikan distributed lock utilitas kita sendiri menggunakan perintah `SET key value NX PX` dan pelepasan lock menggunakan script Lua di Go (didokumentasikan di [ADR-002](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/07-project-digital-wallet/adr/002-redis-lock-strategy.md)).

---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
