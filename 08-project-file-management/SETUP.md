# Setup Guide: File Management Service

---

## 1. Prerequisites

Pastikan komputer Anda telah terinstal:
- Go 1.20+
- Docker & Docker Compose
- Alat penguji API seperti `curl` atau Postman
- PostgreSQL client (`psql`) - Opsional

---

## 2. Local Installation

1. **Inisiasi Konfigurasi:**
   Salin berkas konfigurasi template environment:
   ```bash
   cp .env.example .env
   ```

2. **Menjalankan PostgreSQL & MinIO Containers:**
   Jalankan kontainer database relasional PostgreSQL dan MinIO di background:
   ```bash
   docker-compose up -d
   ```

3. **Dashboard Web Console MinIO:**
   Buka dashboard administrasi MinIO di browser:
   - **URL:** `http://localhost:9001`
   - **Username:** `minioadmin`
   - **Password:** `minioadmin`
   - *Verifikasi bahwa bucket `user-files` otomatis terbuat sesaat setelah server Go diaktifkan.*

4. **Jalankan Aplikasi:**
   Nyalakan server backend lokal Go:
   ```bash
   go run cmd/server/main.go
   ```
   Server akan berjalan di port `8080` dan melakukan auto-migration tabel database relasional.

5. **Jalankan Unit Test:**
   ```bash
   go test -v ./...
   ```
   Unit test otomatis melakukan pengujian otentikasi JWT dan pembatalan compensating write secara mock.

---

## 3. Manual Testing Walkthrough (cURL)

### 1. Registrasi Akun & Login
```bash
# Register User
curl -i -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "usera@email.com", "password": "password123"}'

# Login User
curl -i -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "usera@email.com", "password": "password123"}'
```
*Salin token JWT yang dihasilkan.*

### 2. Buat File Dummy & Upload
Buat file text dummy gambar:
```bash
echo "PNG file mock" > dummy.png
```

Kirim berkas multipart form-data:
```bash
curl -i -X POST http://localhost:8080/files/upload \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -F "file=@dummy.png;type=image/png"
```
*Metadata berkas terekam di PostgreSQL dengan ID 1.*

### 3. Peroleh Presigned Download Link
```bash
curl -i http://localhost:8080/files/1/download \
  -H "Authorization: Bearer <JWT_TOKEN>"
```
*Kueri mengembalikan tautan download langsung dari MinIO. Coba buka tautan tersebut di browser Anda untuk mengunduh berkas.*

### 4. Direct Streaming View
Berguna untuk melayani render gambar di tag `<img src="...">` privat:
```bash
curl -i http://localhost:8080/files/1/view \
  -H "Authorization: Bearer <JWT_TOKEN>"
```
*Mengembalikan isi byte biner file secara streaming.*

### 5. Hapus Berkas
```bash
curl -i -X DELETE http://localhost:8080/files/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```
*Berkas fisik terhapus dari MinIO dan metadatanya dihapus di PostgreSQL.*
