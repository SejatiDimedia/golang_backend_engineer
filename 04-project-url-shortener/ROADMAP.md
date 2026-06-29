# Roadmap: URL Shortener Service

**Status:** `Planning`

Dokumen ini mengurutkan alur pengembangan fitur dalam proyek URL Shortener Service untuk memastikan pengerjaan berjalan teratur dan terdokumentasi.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Inisialisasi Go module, struktur Clean Architecture, pembacaan konfigurasi env, koneksi GORM ke PostgreSQL, dan endpoint health check. | — | `Planned` |
| 2 — Core Domain & Migration | Definisi skema GORM untuk URL entity, auto-migration database, dan utilitas pembuatan short code berbasis Base64 URL-safe dari timestamp. | Phase 1 | `Planned` |
| 3 — Shorten URL Features | REST API endpoint untuk membuat short URL (acak) dan custom alias (FR-1, FR-2). | Phase 2 | `Planned` |
| 4 — Redirect & Statistics | REST API endpoint untuk pengalihan URL (`GET /r/:short_code`) dengan penambahan hitungan klik dan validasi kedaluwarsa, serta endpoint statistik (FR-3, FR-4, FR-5, FR-6). | Phase 3 | `Planned` |
| 5 — Hardening | Penulisan unit tests untuk logika bisnis (Service layer) dan Handler, pemolesan validasi input, serta formatting standard error response. | All above | `Planned` |
| 6 — Deployment | Dockerization (Dockerfile & `docker-compose.yml`) serta penyelesaian 12 dokumen proyek standar. | Phase 5 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| Health Check | FR-7 | Wajib ada pertama kali untuk memvalidasi bahwa sasis aplikasi, koneksi database, dan routing dasar berfungsi dengan benar. |
| URL Shortening (Random) | FR-1 | Fondasi utama untuk mengisi data target URL ke database. Harus ada sebelum kita bisa menguji fitur pengalihan (*redirect*). |
| Custom Alias Validation | FR-2 | Merupakan variasi dari pembuatan URL shortener. Menguji logika validasi keunikan alias di database. |
| URL Redirection | FR-3 | Logika utama konsumsi short URL. Bergantung pada data yang dibuat di FR-1/FR-2. |
| Expiration Check | FR-4 | Bagian dari alur pengalihan. Pengalihan harus menolak akses jika tautan sudah kedaluwarsa. |
| Click Counter | FR-5 | Logika pencatatan aktivitas yang dipicu tepat sebelum proses pengalihan sukses dilakukan. |
| Statistics Endpoint | FR-6 | Endpoint untuk melihat metadata URL dan jumlah klik. Bergantung pada data yang ada di database. |

## 3. Concepts this project is exercising

- **Clean Architecture dasar:** Memisahkan tanggung jawab Handler (HTTP), Service (Business Logic), dan Repository (Database/GORM).
- **Interface Segregation:** Menggunakan interface kecil untuk abstraksi repository guna memudahkan mocking saat pengujian.
- **Go Context Propagation:** Meneruskan `context.Context` dari HTTP request ke layer service dan database I/O.
- **Go Error Handling:** Mengembalikan error secara eksplisit dan membungkusnya (*wrapping*) dengan informasi kontekstual yang jelas.
- **Docker & Compose:** Menjalankan aplikasi backend Go dan PostgreSQL secara terisolasi.

## 4. Known risks / unknowns at planning time

- **Concurrency pada Click Counter:** Penambahan jumlah klik (`click_count`) secara simultan dapat menyebabkan kondisi balapan (*race condition*) jika diakses secara concurrent tinggi. Kita perlu menggunakan query update atomic (misal `gorm.Expr("click_count + 1")`) alih-alih mengambil data ke memori, menambahkannya, lalu menyimpannya kembali.
- **Panjang Short Code & Tabrakan:** Encoding Base64 dari timestamp mikrodetik menghasilkan string sekitar 8-10 karakter. Kami perlu memotongnya jika ingin lebih pendek, namun harus meminimalkan risiko tabrakan jika ada request masuk di saat yang sama.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek URL Shortener |
