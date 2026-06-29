# Testing Strategy: Digital Wallet API

---

## 1. Scope of testing for this project

Pada Project 4, fokus testing ditingkatkan untuk menguji **konsistensi saldo di bawah beban konkuren (race conditions)** dan keandalan sistem penyaringan **idempotensi request**.
Kami menguji:
- Operasi Top-up, Withdraw, dan perhitungan saldo ledger di Service Layer.
- Uji Konkurensi Transfer: Menjalankan Goroutines paralel secara bersamaan untuk melakukan transfer saldo dari akun A ke akun B guna memvalidasi efektivitas **Redis lock** dalam melindungi saldo dari double-spending.
- Uji Idempotensi Middleware: Memastikan request dengan header key yang sama mengembalikan respons identik dari cache Redis tanpa memicu eksekusi ganda di database.

## 2. Test types in use

| Type | Used? | Tooling | Scope |
|---|---|---|---|
| **Unit tests (Service)** | Yes | Go standard `testing` | Menguji logika bisnis, kalkulasi topup/withdraw, dan simulasi transfer konkuren (10 paralel workers) memanfaatkan mock local mutex lock manager. |
| **Unit tests (Middleware)** | Yes | Gin HTTP recorder `httptest` | Menguji parser middleware JWT. |
| **Integration tests** | No | — | Integrasi database relasional nyata PostgreSQL dan cache Redis diuji secara manual. |

## 3. What is covered

### 1. `WalletService` Tests (`internal/service/wallet_test.go`)
- `TestWalletService_TopUpAndWithdraw`: Validasi kalkulasi tambah/kurang saldo dan deteksi error saldo tidak cukup.
- `TestWalletService_ConcurrentTransferSafety`: Menjalankan 10 Goroutines paralel secara bersamaan mentransfer Rp 10.000 dari Wallet 1 ke Wallet 2. Memverifikasi saldo akhir Wallet 1 berkurang tepat Rp 100.000 dan saldo Wallet 2 bertambah tepat Rp 100.000 tanpa deadlock.

---

## 4. Running tests

Jalankan pengujian unit lokal:
```bash
go test -v ./...
```

Periksa persentase cakupan pengujian unit:
```bash
go test -cover ./...
```

**Hasil Pengujian Saat Ini:**
- `internal/service`: **100% PASS** (Mencakup pengetesan saldo konkuren & deadlock lock-ordering).

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen pengujian saldo konkuren dan distributed locking |
