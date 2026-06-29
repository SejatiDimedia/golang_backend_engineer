# Setup Guide: Notification Service

---

## 1. Prerequisites

Pastikan komputer Anda telah terinstal:
- Go 1.20+
- Docker & Docker Compose
- Alat penguji API seperti `curl` atau Postman

---

## 2. Local Installation

1. **Inisiasi Konfigurasi:**
   Salin berkas konfigurasi template environment:
   ```bash
   cp .env.example .env
   ```

2. **Menjalankan PostgreSQL & Redis Containers:**
   Jalankan kontainer database relasional PostgreSQL dan Redis:
   ```bash
   docker-compose up -d
   ```

3. **Jalankan Aplikasi:**
   Nyalakan server backend lokal Go:
   ```bash
   go run cmd/server/main.go
   ```
   Server akan berjalan di port `8080`, melakukan auto-migration tabel, meluncurkan 5 worker paralel, dan mengaktifkan loop poller scheduler 1 detik.

4. **Jalankan Unit Test:**
   ```bash
   go test -v ./...
   ```
   Unit test otomatis melakukan pengujian `redisQueueManager` (LPUSH/BRPOP & Lua scheduler) dan formula penundaan exponential backoff.

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
*Salin token JWT.*

### 2. Mengirimkan Notifikasi Instan (Email)
Kirim request notifikasi:
```bash
curl -i -X POST http://localhost:8080/notifications \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"type": "email", "target": "user@gmail.com", "content": "Halo, ini email penting!"}'
```
*API merespon secara instan `202 Accepted` dengan ID notifikasi (misal ID 1).*

### 3. Pantau Log Tunda Retry (Exponential Backoff) di Terminal
Di terminal tempat server Go berjalan, perhatikan lognya. 
- Provider dummy diatur memiliki 30% tingkat kegagalan acak.
- Jika percobaan pertama gagal, worker memanggil requeue backoff:
  `[Worker 1] FAILED: Notification ID 1 failed...`
  `[Worker 1] RETRY QUEUED: Notification ID 1 will retry in 4 seconds...`
- Setelah 4 detik ( jatuh tempo delay ), poller memindahkan notifikasi kembali ke antrean utama untuk dieksekusi ulang.

### 4. Periksa Riwayat Status & Audit Logs
Gunakan ID notifikasi untuk memeriksa audit logs-nya:
```bash
curl -i http://localhost:8080/notifications/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```
*Kueri akan mengembalikan status notifikasi terupdate (`SENT` atau `FAILED` jika melebihi 5 kali coba) beserta riwayat error log di setiap attempt.*

### 5. Menguji Notifikasi Terjadwal (Scheduled Notification)
Kirim notifikasi dengan waktu kirim 15 detik di masa depan.
Dapatkan timestamp UTC 15 detik di masa depan (contoh: `2026-06-29T18:00:15Z`).
```bash
curl -i -X POST http://localhost:8080/notifications \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"type": "webhook", "target": "https://httpbin.org/post", "content": "data", "send_at": "2026-06-29T18:00:15Z"}'
```
*Tugas akan tertahan di Redis Sorted Set, dan log terminal akan menunjukkan poller memindahkan tugas tersebut ke antrean utama tepat saat 15 detik berlalu.*
