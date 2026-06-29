# Testing Strategy: URL Shortener Service

---

## 1. Scope of testing for this project

Untuk proyek pertama (Beginner), fokus pengujian kami adalah pada **Unit Testing** untuk memvalidasi:
1. Logika bisnis inti di Service Layer (validasi URL, pembuatan short code unik, dan pengecekan kedaluwarsa URL).
2. Perilaku HTTP Routing dan serialization di Handler Layer menggunakan HTTP mock recorders.

Strategi testing sengaja difokuskan pada unit testing terisolasi (menggunakan mock memori untuk database repository) untuk memperkuat pemahaman fundamental unit testing Go tanpa kerumitan interaksi database nyata.

## 2. Test types in use

| Type | Used? | Tooling | Scope |
|---|---|---|---|
| **Unit tests (Service)** | Yes | Standard `testing` library | Menguji logika bisnis di `internal/service`, terisolasi dari database asli menggunakan mock repository memori sederhana. |
| **Unit tests (Handler)** | Yes | Standard `testing` + Gin `httptest` | Menguji response body, HTTP status codes, dan parsing request di `internal/handler` dengan merekam request simulasi. |
| **Integration tests** | No | — | Tidak digunakan di Project 1. Interaksi database riil divalidasi secara manual melalui health check dan pengujian client (Postman/cURL) terhadap DB PostgreSQL Docker lokal. |
| **Load/performance tests** | No | — | Di luar ruang lingkup proyek pertama ini. |

## 3. What is covered

| Component | Coverage approach |
|---|---|
| **Handlers** | Divalidasi menggunakan `httptest.ResponseRecorder` untuk memastikan parameter input tidak valid menghasilkan HTTP 400, alias konflik menghasilkan HTTP 409, short code tidak ditemukan menghasilkan HTTP 404, dan redirect sukses menghasilkan HTTP 302 dengan header `Location` yang benar. |
| **Services** | Menguji `Shorten`, `GetAndRecordClick`, dan `GetStats`. Memastikan validasi URL menolak tautan tanpa `http://` atau `https://`, memastikan alias kustom dicek keunikannya, dan memastikan URL kedaluwarsa diblokir dari proses pengalihan. |
| **Repositories** | Tidak diuji unit secara langsung karena hanya memuat kueri dasar GORM. Operasi DB repository diuji secara tidak langsung lewat manual validation API. |

## 4. What is explicitly NOT covered, and why

- **Pengujian Database Integrasi (Integration DB Test):** Kami tidak membuat pengujian otomatis terhadap database PostgreSQL asli menggunakan *test database instances*. Hal ini ditunda karena setup infrastruktur pengujian database relasional secara otomatis (migration & database teardown per test) dinilai terlalu rumit untuk Project 1 dan akan diperkenalkan secara resmi pada proyek berikutnya yang memiliki transaksi database lebih kompleks.
- **Race Condition Concurrency Test:** Kami tidak menulis automated test konkurensi tinggi untuk click counter. Masalah konkurensi diatasi secara teoretis melalui query SQL update atomic di repository layer.

## 5. Test data strategy

Untuk unit test, kami tidak menggunakan database asli atau berkas sql eksternal. Kami menggunakan `mockURLRepository` di [internal/service/url_test.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/internal/service/url_test.go) yang mengadopsi struktur `map[string]*entity.URL` di memori. Data dipersiapkan (seed) langsung di awal setiap fungsi test (seperti `TestGetAndRecordClick_IncrementClickCount`) dan otomatis terhapus saat siklus fungsi test berakhir (*garbage-collected*).

## 6. Running tests

Jalankan seluruh test suite dengan menampilkan output log detail:
```bash
go test -v ./...
```

Jalankan test suite untuk melihat persentase coverage kode:
```bash
go test -cover ./...
```

**Hasil Pengujian Saat Ini:**
- `internal/service`: **100% test pass** (menguji fungsionalitas core logic)
- `internal/handler`: **100% test pass** (menguji fungsionalitas routing HTTP dan binding)

## 7. CI integration

Pengujian integrasi CI/CD (seperti GitHub Actions) saat ini dideferensiasi. Pengujian dijalankan secara manual di terminal pengembang sebelum kode di-commit. Otomasi CI akan diperkenalkan pada Phase 3/4 setelah kesiapan infrastruktur deployment kami lebih matang.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen strategi pengujian unit untuk service dan handler |
