# Testing Strategy: Booking Management System

---

## 1. Scope of testing for this project

Fokus utama pengujian di Project 3 adalah memvalidasi:
1. **Algoritma Deteksi Overlap:** Memastikan irisan waktu yang tumpang tindih terdeteksi akurat dan pemesanan bentrok ditolak oleh sistem.
2. **Aturan Pembatalan Waktu (Cancellation Limits):** Memastikan batas waktu 2 jam dipatuhi untuk role `customer` dan dilewati (*bypass*) untuk role `admin`.
3. **JWT Authentication & Role Middleware:** Memastikan token JWT diparsing benar dan route khusus admin terproteksi dari akses customer.

## 2. Test types in use

| Type | Used? | Tooling | Scope |
|---|---|---|---|
| **Unit tests (Service)** | Yes | Go standard `testing` | Menguji logika bisnis `BookingService` menggunakan in-memory map repository mock. |
| **Unit tests (Middleware)** | Yes | Gin HTTP recorder `httptest` | Menguji parser middleware JWT dan otorisasi role. |
| **Integration tests** | No | — | Pengujian database relasional penuh dan constraint diuji secara manual. |

## 3. What is covered

### 1. `BookingService` Tests (`internal/service/booking_test.go`)
- `TestCreateBooking_Success`: Booking baru pada slot kosong berhasil disimpan dengan status `CONFIRMED`.
- `TestCreateBooking_OverlapConflict`: Pendaftaran booking baru pada jam yang tumpang tindih dengan booking aktif diblokir dan mengembalikan error `ErrDoubleBooking`.
- `TestCancelBooking_CancellationWindowClosed`:
  - Membatalkan booking <2 jam oleh customer -> Ditolak (`ErrCancellationWindowClosed`).
  - Membatalkan booking <2 jam oleh admin -> Diizinkan (bypass sukses).

### 2. `AuthMiddleware` Tests (`internal/middleware/auth_test.go`)
- `TestAuthMiddleware_Success`: Token JWT valid didekode dan data user dimasukkan ke Gin Context.
- `TestAuthMiddleware_MissingHeader`: Request tanpa header Authorization diblokir dengan HTTP 401.
- `TestRequireRole_Forbidden`: Customer diblokir dengan HTTP 403 saat mengakses route ber-role `admin`.

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
- `internal/middleware`: **100% PASS**
- `internal/service`: **100% PASS**

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen pengetesan unit middleware dan overlap validation |
