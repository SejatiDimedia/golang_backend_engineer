# Database Design: URL Shortener Service

**Status:** `Implemented`
**Engine:** PostgreSQL 15

---

## 1. Entity-relationship overview

Proyek ini sangat sederhana dan hanya memiliki satu tabel entitas utama (`urls`). Tidak ada hubungan relasi (*foreign keys*) dengan tabel lain pada versi awal ini karena autentikasi user sengaja ditiadakan (*non-goal*).

```
┌───────────────────────────────────────────┐
│                   urls                    │
├───────────────────────────────────────────┤
│ id (PK)                                   │
│ long_url                                  │
│ short_code (Unique Index)                 │
│ click_count                               │
│ created_at                                │
│ expires_at                                │
└───────────────────────────────────────────┘
```

## 2. Schema

### Table: `urls`

Menyimpan semua pemetaan short code ke target URL asli beserta statistik klik dan masa kedaluwarsa.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | BIGINT (auto-increment) | PRIMARY KEY | Dihasilkan secara otomatis oleh PostgreSQL. |
| long_url | TEXT | NOT NULL | Menyimpan URL tujuan asli yang panjang. |
| short_code | VARCHAR(50) | UNIQUE INDEX, NOT NULL | Kode pendek yang dihasilkan dari base64 timestamp atau kustom alias. |
| click_count | BIGINT | NOT NULL, DEFAULT 0 | Total hitungan kunjungan/klik yang berhasil. |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | Waktu pembuatan URL pendek. |
| expires_at | TIMESTAMP WITH TIME ZONE | NULLABLE | Waktu kedaluwarsa opsional. Jika NULL, URL tidak pernah kedaluwarsa. |

## 3. Relationships

Tidak ada relasi antar-tabel di versi awal.

## 4. Indexes

| Table | Index | Columns | Type | Reason |
|---|---|---|---|---|
| urls | `idx_urls_short_code` | `short_code` | UNIQUE | Pencarian baris URL paling sering dilakukan menggunakan `short_code` (saat redirect & statistik). Indeks unik wajib untuk menjamin pencarian instan (O(1)) dan mencegah duplikasi kode. |

## 5. Transactions and consistency

Meskipun sistem ini beroperasi secara sederhana, ada satu operasi penulisan yang sensitif terhadap konkurensi: **Pembaruan hitungan klik** (`click_count`).

| Operation | Transaction boundary | Isolation concern |
|---|---|---|
| Increment click (`click_count`) | Pembaruan satu baris data (single row update). | **Race Condition:** Jika dua request pengalihan masuk bersamaan untuk kode yang sama, pengambilan data ke memori (read) lalu penulisan ulang (write) dapat menyebabkan klik hilang (*lost update*). <br><br>**Solusi:** Kami menggunakan query update SQL atomic secara langsung di database: `UPDATE urls SET click_count = click_count + 1 WHERE short_code = ?` tanpa membacanya ke memori aplikasi terlebih dahulu. |

## 6. Migrations strategy

Pada proyek pertama ini, kami menggunakan fitur **AutoMigration GORM** (`db.AutoMigrate(&entity.URL{})`) yang dideklarasikan di dalam file [main.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/cmd/server/main.go) ketika program pertama kali dijalankan. Strategi ini mempercepat siklus awal pengembangan.

## 7. Sample queries

### Create Shortened URL
```sql
INSERT INTO urls (long_url, short_code, click_count, created_at, expires_at)
VALUES ('https://google.com', 'abcd123', 0, NOW(), NULL);
```

### Fetch Target URL (Redirect & Check Expiry)
```sql
SELECT long_url, expires_at 
FROM urls 
WHERE short_code = 'abcd123' 
LIMIT 1;
```

### Atomic Click Counter Update
```sql
UPDATE urls 
SET click_count = click_count + 1 
WHERE short_code = 'abcd123';
```

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi desain database PostgreSQL 15 menggunakan tabel single `urls` |
