# Setup: Inventory Management API

## Prerequisites

- **Go 1.21** atau versi lebih baru.
- **Docker & Docker Compose** (Sangat disarankan untuk mempermudah instalisasi database relasional).
- **PostgreSQL 15** (Jika menggunakan manual installation).

## Environment variables

Konfigurasi disimpan di file `.env` di root proyek.

| Variable | Description | Default | Example |
|---|---|---|---|
| `PORT` | Port server HTTP mendengarkan request | `8080` | `8080` |
| `ENV` | Mode aplikasi (`development` / `production`) | `development` | `development` |
| `DB_HOST` | Host PostgreSQL database | `localhost` | `localhost` |
| `DB_PORT` | Port PostgreSQL database | `5432` | `5432` |
| `DB_USER` | Username database PostgreSQL | `postgres` | `postgres` |
| `DB_PASSWORD` | Password database PostgreSQL | `postgres` | `postgres` |
| `DB_NAME` | Nama database relasional | `inventory_db` | `inventory_db` |
| `DB_SSLMODE` | Mode enkripsi koneksi SSL database | `disable` | `disable` |

Salin contoh konfigurasi ke berkas `.env` Anda:
```bash
cp .env.example .env
```

---

## Local setup (With Docker — Recommended)

1. **Jalankan Database PostgreSQL:**
   ```bash
   docker-compose up -d
   ```
   *Ini akan meluncurkan PostgreSQL 15 di port 5432.*

2. **Jalankan Aplikasi Go:**
   ```bash
   go run cmd/server/main.go
   ```
   *GORM secara otomatis memigrasi tabel Category, Supplier, Product, dan StockMovement saat boot pertama.*

---

## Local setup (Without Docker)

1. **Jalankan PostgreSQL lokal Anda** dan buat database kosong bernama `inventory_db`.
2. **Perbarui berkas `.env`** dengan username dan password PostgreSQL milik sistem lokal Anda.
3. **Jalankan server aplikasi:**
   ```bash
   go run cmd/server/main.go
   ```

---

## Verifying it's running

Kirim request HTTP ke health check endpoint:
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

Untuk menjalankan unit test Service Layer (termasuk validasi transaksi atomic mutasi stok masuk/keluar, rollback transaksi, dan validasi stok minus):
```bash
go test -v ./...
```
Untuk mengukur coverage kueri unit test:
```bash
go test -cover ./...
```
Detail strategi pengujian dapat dibaca di [TESTING.md](./TESTING.md).

## Troubleshooting

| Issue | Likely cause | Fix |
|---|---|---|
| `cannot delete category... (foreign key constraint)` | Mencoba menghapus kategori yang masih dirujuk oleh salah satu produk. | Hapus produk terkait terlebih dahulu, atau ubah kategori produk ke kategori lain sebelum menghapus kategori asal. |
| `insufficient stock quantity` | Mengirimkan `POST /products/:id/stock-out` dengan nilai `quantity` melebihi total `stock_quantity` saat ini. | Tambahkan stok terlebih dahulu lewat endpoint `/stock-in`, atau periksa kuantitas produk saat ini lewat `GET /products/:id`. |
