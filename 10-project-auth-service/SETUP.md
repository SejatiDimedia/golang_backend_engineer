# Setup Guide: Authentication Service

---

## 1. Prerequisites

Pastikan komputer Anda telah terinstal:
- Go 1.20+
- Docker & Docker Compose
- Alat penguji API seperti `curl` atau Postman

---

## 2. Local Installation

1. **Inisiasi Konfigurasi:**
   Salin berkas konfigurasi environment:
   ```bash
   cp .env.example .env
   ```

2. **Menjalankan PostgreSQL & Redis Containers:**
   Jalankan kontainer database relasional PostgreSQL dan Redis:
   ```bash
   docker-compose up -d
   ```

3. **Jalankan Aplikasi:**
   Nyalakan server lokal Go:
   ```bash
   go run cmd/server/main.go
   ```
   Server akan berjalan di port `8081` (Port sengaja diset ke 8081 agar tidak bentrok dengan Notification Service di port 8080).
   *RSA keypair `certs/private.key` dan `certs/public.key` otomatis ter-generate di folder root proyek saat server booting.*

4. **Jalankan Unit Test:**
   ```bash
   go test -v ./...
   ```

---

## 3. Manual Testing Walkthrough (cURL)

### 1. Registrasi Akun & Aktivasi Verifikasi Email
```bash
# 1. Register User Baru
curl -i -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user1@email.com", "password": "password123"}'

# 2. Periksa log stdout server Go. Salin token verifikasi dari link:
# http://localhost:8081/auth/verify-email?token=<TOKEN_UUID>

# 3. Lakukan verifikasi email
curl -i http://localhost:8081/auth/verify-email?token=<TOKEN_UUID>
```

### 2. Login & Dapatkan Token
```bash
curl -i -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user1@email.com", "password": "password123"}'
```
*Salin `access_token` JWT RS256 dan `refresh_token` UUID dari respons JSON.*

### 3. Menguji Token Introspect (Online Verification)
```bash
curl -i -X POST http://localhost:8081/auth/introspect \
  -H "Content-Type: application/json" \
  -d '{"token": "<ACCESS_TOKEN>"}'
```
*Kueri mengembalikan status `active: true` beserta data claims, roles (default: `customer`), dan permissions.*

### 4. Menguji Token Rotation (RTR)
```bash
curl -i -X POST http://localhost:8081/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<REFRESH_TOKEN>"}'
```
*Anda mendapatkan `access_token` baru dan `refresh_token` baru. Refresh token lama otomatis di-revoke.*

### 5. Simulasi Replay Attack Safeguard
Kirim kembali request `/auth/refresh` menggunakan refresh token *lama* yang barusan di-refresh di langkah 4.
```bash
curl -i -X POST http://localhost:8081/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<REFRESH_TOKEN_LAMA>"}'
```
*API mengembalikan `401 Unauthorized` dengan pesan error token invalid. Periksa database: seluruh refresh token milik user tersebut otomatis diubah menjadi `is_revoked = true` demi melindungi data dari peretas.*
