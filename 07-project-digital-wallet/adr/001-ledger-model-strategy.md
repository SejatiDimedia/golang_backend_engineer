# ADR-001: Model Desain Akuntansi Dompet Digital

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Untuk aplikasi finansial seperti dompet digital (*digital wallet*), konsistensi saldo adalah prioritas utama. Kami perlu menentukan struktur tabel database untuk menyimpan dan melacak saldo pengguna agar terhindar dari ketidakakuratan saldo akibat kegagalan transaksi.

## Decision

Kami memutuskan untuk mengadopsi **Model Akuntansi Hibrida (Running Balance + Double-Entry Ledger)**.

Dalam desain ini:
1. Tabel `wallets` menyimpan kolom `balance` (running balance) saat ini untuk melayani kueri baca saldo dengan sangat cepat ($O(1)$).
2. Tabel `wallet_transactions` bertindak sebagai buku besar pembantu (*ledger*) yang mencatat mutasi debit/kredit setiap transaksi (top-up, withdraw, transfer).
3. Setiap mutasi saldo wajib dieksekusi di dalam transaksi database SQL tunggal yang memperbarui tabel `wallets` dan menambahkan baris mutasi ke `wallet_transactions` secara atomik (jika salah satu gagal, seluruh pembaruan dibatalkan).

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Running Balance & Ledger (Chosen)** | - Kueri baca saldo instan ($O(1)$).<br>- Memiliki catatan audit trail historis lengkap di tabel ledger untuk proses rekonsiliasi. | - Memerlukan penulisan kode transaksi SQL yang sangat hati-hati di server untuk sinkronisasi dua tabel. |
| **B. Pure Ledger (Murni Kueri SUM)** | - Sangat aman dari kebocoran saldo karena tidak ada update manual pada kolom saldo.<br>- Audit trail adalah satu-satunya sumber kebenaran data. | - Degradasi performa yang parah seiring bertambahnya jutaan entri transaksi historis user (kueri `SUM` menjadi sangat lambat). |

## Reasoning

Dalam sistem produksi finansial nyata, kueri baca saldo pengguna dipanggil sangat sering (hampir setiap kali aplikasi dibuka atau halaman dimuat). Melakukan kueri agregasi `SUM(amount)` dari tabel transaksi historis jutaan baris (Opsi B) akan membebani database PostgreSQL secara drastis. 

Opsi A menyelesaikan masalah performa ini dengan menaruh running balance di tabel `wallets`, namun melindunginya dengan transaksi atomic. Jika auditor mendeteksi adanya kecurigaan selisih saldo, kami dapat dengan mudah merekonstruksi saldo valid dengan menjumlahkan data historis di tabel ledger dan membandingkannya dengan running balance saat itu (*recon process*).

## Consequences

- **Positif:** Kecepatan kueri baca sangat tinggi, audit log finansial lengkap, data terjamin konsisten.
- **Negatif:** Skema database sedikit lebih kompleks karena melibatkan dua entitas yang saling terikat erat.
