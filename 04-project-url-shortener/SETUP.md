# Setup: URL Shortener Service

## Prerequisites

- **Go 1.21** atau versi lebih baru.
- **Docker & Docker Compose** (Opsional tetapi sangat direkomendasikan untuk menjalankan database).
- **PostgreSQL 15** (Jika tidak menggunakan Docker, Anda harus menyiapkan server database PostgreSQL lokal secara manual).

## Environment variables

Aplikasi membaca konfigurasi dari variabel lingkungan. Anda dapat mengonfigurasinya melalui file `.env`.

| Variable | Description | Default | Example |
|---|---|---|---|
| `PORT` | Port tempat HTTP server mendengarkan request | `8080` | `8080` |
| `ENV` | Mode aplikasi (`development` atau `production`) | `development` | `development` |
| `DB_HOST` | Host PostgreSQL database server | `localhost` | `localhost` |
| `DB_PORT` | Port PostgreSQL database server | `5432` | `5432` |
| `DB_USER` | Username login PostgreSQL | `postgres` | `postgres` |
| `DB_PASSWORD` | Password login PostgreSQL | `postgres` | `postgres` |
| `DB_NAME` | Nama database yang digunakan | `url_shortener` | `url_shortener` |
| `DB_SSLMODE` | Pengaturan enkripsi SSL database | `disable` | `disable` |

Copy template konfigurasi yang disediakan untuk membuat berkas konfigurasi lokal:
```bash
cp .env.example .env
```

---

## Local setup (With Docker — Recommended)

Cara termudah untuk memulai adalah menjalankan database PostgreSQL menggunakan Docker Compose, sementara backend dijalankan secara lokal untuk mempermudah debugging dan pengembangan.

1. **Jalankan Database PostgreSQL:**
   ```bash
   docker-compose up -d
   ```
   *Perintah ini akan memulai PostgreSQL 15 di port `5432` dengan database `url_shortener`.*

2. **Jalankan Aplikasi Go:**
   ```bash
   go run cmd/server/main.go
   ```
   *Skema database akan otomatis dimigrasi oleh GORM saat server dijalankan.*

---

## Local setup (Fully Manual — Without Docker)

Jika Anda tidak memasang Docker:

1. **Siapkan PostgreSQL:**
   Instal PostgreSQL di sistem Anda, jalankan servicenya, lalu buat database kosong dengan nama `url_shortener`.

2. **Sesuaikan `.env`:**
   Sesuaikan `DB_USER`, `DB_PASSWORD`, dan detail koneksi database di file `.env` dengan kredensial PostgreSQL lokal Anda.

3. **Unduh Dependensi & Jalankan Aplikasi:**
   ```bash
   go mod download
   go run cmd/server/main.go
   ```

---

## Verifying it's running

Kirim request HTTP ke health check endpoint untuk memverifikasi aplikasi berjalan dan sukses terhubung ke database:

```bash
curl http://localhost:8080/health
```

**Expected response:**
```json
{
  "database": "connected",
  "status": "healthy"
}
```

---

## Running tests

Untuk menjalankan seluruh unit test yang mencakup logika bisnis di Service layer dan routing di Handler layer:

```bash
go test -v ./...
```

Untuk detail cakupan pengujian, silakan lihat [TESTING.md](./TESTING.md).

## Troubleshooting

| Issue | Likely cause | Fix |
|---|---|---|
| `Fatal: failed to connect to database...` | Container Docker PostgreSQL belum berjalan atau port `5432` terpakai proses lain. | Pastikan `docker ps` menunjukkan container database berjalan. Hentikan postgres lokal lain jika terjadi tabrakan port. |
| `invalid request body or malformed URL...` | Format JSON request body salah, atau nilai parameter `long_url` tidak diawali `http://` atau `https://`. | Pastikan payload valid dan tautan yang dimasukkan berformat URL lengkap (misal `https://example.com`). |
