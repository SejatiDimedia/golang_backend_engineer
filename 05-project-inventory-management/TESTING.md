# Testing Strategy: Inventory Management API

---

## 1. Scope of testing for this project

Pada Project 2 (Beginner->Intermediate), fokus testing ditingkatkan untuk menguji **Logika Transaksional Database (Begin/Commit/Rollback)** secara aman di tingkat unit test.
Kami menguji keandalan operasi mutasi stok masuk (Stock In) dan keluar (Stock Out) di Service Layer menggunakan Mock Repository. Uji unit ini memvalidasi:
- Transaksi berhasil melakukan commit perubahan kuantitas stok produk dan log mutasi.
- Kegagalan update database di salah satu tahap berhasil memicu kegagalan transaksi keseluruhan (*rollback simulation*).
- Mutasi barang keluar dibatalkan jika stok saat ini kurang dari yang diminta.

## 2. Test types in use

| Type | Used? | Tooling | Scope |
|---|---|---|---|
| **Unit tests (Service)** | Yes | Standard Go `testing` | Menguji logika transaksional mutasi stok di `internal/service/movement.go` secara terisolasi menggunakan unit test mocks. |
| **Unit tests (Handler)** | Yes | Gin `httptest` recorder | Menguji endpoint API mutasi stok (`POST /products/:id/stock-in` dan `/stock-out`) dengan merekam HTTP request/response simulasi. |
| **Integration tests** | No | — | Tidak dilakukan secara otomatis. Integrasi database relasional dan constraint kunci asing (restrict delete) diuji secara manual. |

## 3. What is covered

| Component | Coverage approach |
|---|---|
| **`MovementService`** | Diuji unit penuh di [internal/service/movement_test.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/internal/service/movement_test.go): <br>- `TestStockIn_Success`: Sukses menambah kuantitas stok produk.<br>- `TestStockOut_Success`: Sukses mengurangi kuantitas stok produk.<br>- `TestStockOut_InsufficientStock`: Memverifikasi validasi stok tidak boleh negatif mengembalikan error.<br>- `TestStockIn_TransactionRollbackOnDBError`: Memverifikasi jika update stok produk gagal, mutasi log tidak dicatat (rollback simulasi). |
| **`MovementHandler`** | Diuji unit di [internal/handler/movement_test.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/internal/handler/movement_test.go): <br>- `TestStockIn_HandlerSuccess`: Request JSON valid menghasilkan 200 OK.<br>- `TestStockOut_HandlerInsufficient`: Pengurangan stok berlebih menghasilkan 400 Bad Request. |

## 4. What is explicitly NOT covered, and why

- **Pengujian Database Riil Otomatis (Real DB Integration Test):** Kami tidak membuat otomatisasi pengujian relasi foreign key database relasional asli (seperti restrict delete pada kategori/supplier). Hal ini ditunda karena setup rollback database transaksi asli di sela-sela testing dinilai belum sebanding dengan kematangan infrastruktur CI/CD lokal saat ini. Pembuktian integritas data relational diuji secara manual.

## 5. Test data strategy

Kami menggunakan data simulasi di memori menggunakan struct map sederhana (`map[uint]*entity.Product`) di berkas `movement_test.go`. Data produk diinisialisasi secara lokal di awal fungsi test dan otomatis dibersihkan saat eksekusi test suite berakhir.

## 6. Running tests

Jalankan pengujian unit dengan detail laporan lengkap:
```bash
go test -v ./...
```

Periksa cakupan pengujian unit:
```bash
go test -cover ./...
```

**Hasil Pengujian Saat Ini:**
- `internal/service`: **100% PASS** (Mencakup logika bisnis transaksional dan validasi stok).
- `internal/handler`: **100% PASS** (Mencakup endpoint Gin binding dan responses).

## 7. CI integration

Pengujian integrasi CI/CD otomatis belum diintegrasikan di proyek ini dan direncanakan rilis pada Phase 3/4.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen strategi pengujian transaksi dan handler untuk sistem inventaris |
