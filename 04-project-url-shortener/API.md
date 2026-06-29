# API Specification: URL Shortener Service

**Base URL:** `http://localhost:8080`
**Auth:** `None` (Seluruh endpoint bersifat publik)

---

## Conventions

- Semua response request pembentukan (`POST`) dan statistik (`GET /stats`) berupa JSON.
- Response pengalihan (`GET /r/:short_code`) menghasilkan HTTP Redirect.
- Struktur error response:
  ```json
  {
    "error": "Pesan deskripsi kesalahan"
  }
  ```

## Endpoints

### `GET /health`

**Description:** Memeriksa kesehatan layanan backend dan koneksi database relasional.

**Auth required:** No

**Response — 200 OK:**
```json
{
  "database": "connected",
  "status": "healthy"
}
```

**Response — 500 Internal Server Error:**
```json
{
  "database": "ping failed",
  "error": "sql: database is closed",
  "status": "unhealthy"
}
```

---

### `POST /shorten`

**Description:** Membuat short URL baru dengan opsi custom alias dan waktu kedaluwarsa.

**Auth required:** No

**Request:**
```json
{
  "long_url": "https://example.com/very/long/url/path",
  "custom_alias": "promo-juni",
  "expires_in_seconds": 3600
}
```
*Catatan: `custom_alias` dan `expires_in_seconds` bersifat opsional.*

**Response — 201 Created:**
```json
{
  "expires_at": "2026-06-29T13:30:00Z",
  "long_url": "https://example.com/very/long/url/path",
  "short_code": "promo-juni"
}
```

**Response — error cases:**

| Status | Condition | Example Response |
|---|---|---|
| 400 Bad Request | Payload JSON tidak valid atau format URL salah | `{"error": "invalid request body or malformed URL: Key: 'ShortenRequest.LongURL' Error:Field validation for 'LongURL' failed on the 'url' tag"}` |
| 409 Conflict | Kustom alias yang diajukan sudah digunakan | `{"error": "custom alias already exists"}` |
| 500 Internal Server Error | Kesalahan database saat menyimpan data | `{"error": "failed to shorten URL: database connection error"}` |

---

### `GET /r/:short_code`

**Description:** Mengalihkan browser/client ke URL target asli.

**Auth required:** No

**Response — 302 Found:**
*Mengembalikan HTTP status `302 Found` dengan header `Location` mengarah ke target URL asli.*

**Response — error cases:**

| Status | Condition | Example Response |
|---|---|---|
| 404 Not Found | Short code tidak ditemukan di database | `{"error": "url not found"}` |
| 410 Gone | Short code kedaluwarsa | `{"error": "url has expired"}` |

---

### `GET /stats/:short_code`

**Description:** Melihat metadata dan total statistik klik kunjungan dari short code tertentu.

**Auth required:** No

**Response — 200 OK:**
```json
{
  "id": 1,
  "short_code": "promo-juni",
  "long_url": "https://example.com/very/long/url/path",
  "click_count": 42,
  "created_at": "2026-06-29T12:30:00Z",
  "expires_at": "2026-06-29T13:30:00Z"
}
```

**Response — error cases:**

| Status | Condition | Example Response |
|---|---|---|
| 404 Not Found | Short code tidak ditemukan | `{"error": "url not found"}` |

---

## Rate limiting

Layanan versi awal ini belum memiliki pembatasan laju request (*rate limiting*). Proteksi terhadap brute force kustom alias didelegasikan ke rencana pengembangan masa depan di [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md).

## Versioning strategy

Jika diperlukan versi baru dengan perubahan besar (*breaking changes*), kami akan menggunakan routing berbasis prefix URL (misalnya `/api/v2`). Versi saat ini (`v1`) diekspos langsung pada root path (`/`) untuk mempermudah akses manual di web browser pada endpoint pengalihan (`/r/:short_code`).

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi spesifikasi REST API untuk v1 (health, shorten, redirect, stats) |
