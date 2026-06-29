# Testing Strategy: File Management Service

---

## 1. Scope of testing for this project

Pada Project 5, fokus testing ditujukan untuk memvalidasi **keamanan payload file metadata** (size limit & allowed MIME types) dan keandalan skema transaksional **compensating write rollback** jika upload storage fisik mengalami gangguan.
Kami menguji:
- Penolakan berkas >10MB dan jenis berkas terlarang (seperti `.exe`).
- Pembersihan record database relasional secara otomatis ketika client MinIO mengembalikan error saat melempar berkas.
- Parser token JWT di Gin middleware.

## 2. Test types in use

| Type | Used? | Tooling | Scope |
|---|---|---|---|
| **Unit tests (Service)** | Yes | Go standard `testing` | Menguji logika bisnis `FileService` (validasi ukuran, tipe MIME, dan pembatalan compensating write) menggunakan mock data repository & storage. |
| **Unit tests (Middleware)** | Yes | Gin HTTP recorder `httptest` | Menguji parser middleware JWT. |
| **Integration tests** | No | — | Integrasi database relasional nyata PostgreSQL dan MinIO diuji secara manual. |

## 3. What is covered

### 1. `FileService` Tests (`internal/service/file_test.go`)
- `TestFileService_Upload_Success`: Menguji alur upload berkas valid hingga status berubah sukses.
- `TestFileService_Upload_ValidationLimits`: Menguji penolakan berkas melebihi 10MB atau berekstensi terlarang.
- `TestFileService_Upload_CompensatingRollback`: Mensimulasikan kegagalan upload ke MinIO (disk full), dan memverifikasi baris metadata di PostgreSQL dibersihkan kembali secara otomatis.

### 2. `AuthMiddleware` Tests (`internal/middleware/auth_test.go`)
- `TestAuthMiddleware_MissingHeader`: Menolak request tanpa authorization header.
- `TestAuthMiddleware_InvalidToken`: Menolak token malformed/tidak sah.
- `TestAuthMiddleware_Success`: Mengizinkan request dengan token JWT valid dan menyimpan claims user ID di context.

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

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen strategi pengujian multipart validator dan compensating write |
