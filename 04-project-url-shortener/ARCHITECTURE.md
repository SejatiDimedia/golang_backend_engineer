# Architecture: URL Shortener Service

**Status:** `Implemented`
**Last updated:** 2026-06-29

---

## 1. Architectural style

Layanan ini mengadopsi **Clean Architecture** sederhana yang dibagi menjadi tiga layer utama: HTTP Handler -> Service -> Repository. Gaya ini dipilih untuk mendisiplinkan pemisahan tanggung jawab (Separation of Concerns), memastikan dependensi mengarah ke dalam, serta mempermudah pengujian unit dengan melakukan mocking terhadap interface repository.

## 2. System diagram

```mermaid
graph TD
    Client[Client / Browser]
    
    subgraph Go Application (Gin HTTP Server)
        Router[HTTP Router & Middleware]
        Handler[Handler Layer]
        Service[Service Layer]
        RepoInterface[Repository Interface]
        RepoImpl[Repository Implementation]
    end
    
    subgraph Database
        DB[(PostgreSQL)]
    end
    
    Client -->|HTTP Request| Router
    Router -->|Route mapping| Handler
    Handler -->|Invoke business method| Service
    Service -->|Database boundary| RepoInterface
    RepoInterface -.->|Structural satisfaction| RepoImpl
    RepoImpl -->|GORM DB operation| DB
```

## 3. Folder structure

```
04-project-url-shortener/
├── cmd/
│   └── server/
│       └── main.go         # Entrypoint aplikasi, setup router & DB migration
├── internal/
│   ├── config/
│   │   └── config.go       # Struct dan parsing variabel environment
│   ├── entity/
│   │   └── url.go          # URL entity model & logic internal (IsExpired)
│   ├── handler/
│   │   ├── health.go       # Health check HTTP handler
│   │   ├── url.go          # HTTP handler (Shorten, Redirect, Stats)
│   │   └── url_test.go     # Unit test untuk HTTP handler
│   ├── repository/
│   │   └── url.go          # Interface & implementasi database (GORM)
│   └── service/
│       ├── url.go          # Bisnis logic utama & generator short code
│       └── url_test.go     # Unit test untuk Service layer
├── .env                    # Variabel environment lokal (gitignored)
├── .env.example            # Template variabel environment
├── docker-compose.yml      # DB PostgreSQL lokal untuk dev
└── Dockerfile              # Docker multi-stage build file
```

## 4. Component responsibilities

| Component | Responsibility | Does NOT do |
|---|---|---|
| **Handler** | Menangani masalah protokol HTTP: parsing JSON payload, memvalidasi parameter request, menentukan HTTP Status Code, dan membentuk response JSON/Redirect. | Mengandung logika bisnis inti (seperti pembentukan short code), melakukan interaksi database secara langsung. |
| **Service** | Mengandung logika bisnis inti: memvalidasi target URL, membuat short code unik, mengecek kedaluwarsa URL target, dan memicu penambahan klik kunjungan. | Mengurus status HTTP response, mem-parsing format input payload HTTP, menulis kueri SQL atau memanggil langsung instance GORM. |
| **Repository** | Menyediakan antarmuka data presisten: melakukan insert data URL, mengambil URL berdasarkan kode, dan melakukan update klik counter secara atomic. | Mengecek apakah URL sudah kedaluwarsa (logika bisnis), mem-parsing URL target. |

## 5. Data flow — URL Redirection

Walkthrough alur request ketika pengguna mengakses short URL (`GET /r/:short_code`):

1. Request HTTP tiba di `Redirect` handler di [handler/url.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/internal/handler/url.go).
2. Handler mengambil parameter `short_code` dari route parameter Gin.
3. Handler memanggil `GetAndRecordClick` pada `URLService` dengan context request.
4. Di [service/url.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/internal/service/url.go), service memanggil repository untuk mencari URL berdasarkan `short_code`.
5. Database (PostgreSQL) merespon dengan data URL, lalu dipetakan ke struct entity.
6. Service memverifikasi apakah URL telah kedaluwarsa lewat struct method `IsExpired()`.
   - Jika kedaluwarsa, kembalikan `ErrURLExpired`.
7. Service memanggil database secara sinkron via repository `IncrementClick` untuk menambah hitungan klik secara atomic (`click_count + 1`).
8. Setelah berhasil, Service mengembalikan entitas URL ke Handler.
9. Handler memproses kembalian:
   - Jika sukses, handler memberikan respon **HTTP 302 (Found)** dengan header `Location` terarah ke URL target panjang.
   - Jika terjadi error (misal expired), handler merespon dengan status HTTP yang tepat (misal **410 Gone**).

## 6. Cross-cutting concerns

| Concern | Approach |
|---|---|
| **Logging** | Menggunakan default Logger bawaan Gin untuk melacak request masuk, dan standard library `log` untuk logging inisiasi startup aplikasi. |
| **Error handling** | Error dikembalikan sebagai *values* dari repository -> service -> handler. Handler memetakan jenis error spesifik (seperti `ErrURLNotFound` menjadi HTTP 404, `ErrURLExpired` menjadi HTTP 410). |
| **Configuration** | Menggunakan library standard `os` untuk membaca environment variables dari file `.env` (melalui `godotenv`). |
| **Context propagation** | Parameter `context.Context` selalu di-thread di sepanjang rantai pemanggilan dari Handler ke GORM database query untuk memungkinkan pembatalan request (request cancellation) secara aman. |

## 7. Dependencies on other projects in this repo

Proyek ini tidak memiliki dependensi eksternal pada proyek lain di repository ini karena merupakan proyek pondasi pertama (Project 1). Di masa depan, proyek ini dapat direnovasi untuk menggunakan **Auth Service (Project 7)** guna mengamankan endpoint `/shorten` dan `/stats`.

## 8. Known architectural limitations

- **Click Counter Concurrency:** Logging klik dilakukan secara sinkron tepat sebelum redirect. Pada traffic tinggi, ini membebani PostgreSQL karena write locking. Solusi masa depan: simpan jumlah klik di Redis/Queue, lalu update DB secara berkala (batching).
- **Generator Short Code Semi-Sekuensial:** Base64 dari timestamp mikrodetik tidak sepenuhnya acak. Struktur detail limitasi ini dibahas di [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md).

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi arsitektur proyek URL Shortener |
