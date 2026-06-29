# Setup Guide: Digital Wallet API

---

## 1. Prerequisites

Pastikan komputer Anda telah terinstal:
- Go 1.20+
- Docker & Docker Compose
- Alat penguji API seperti `curl` atau Postman
- PostgreSQL client (`psql`) - Opsional untuk reset database

---

## 2. Local Installation

1. **Inisialisasi Konfigurasi:**
   Salin berkas konfigurasi template environment:
   ```bash
   cp .env.example .env
   ```

2. **Menjalankan PostgreSQL & Redis Containers:**
   Jalankan kontainer database relasional PostgreSQL dan Redis:
   ```bash
   docker-compose up -d
   ```

3. **Inisiasi Database 'wallet_db':**
   Jika database belum terbentuk otomatis, buat secara manual:
   ```bash
   PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE wallet_db;"
   ```

4. **Jalankan Aplikasi:**
   Nyalakan server backend lokal Go:
   ```bash
   go run cmd/server/main.go
   ```
   Server akan berjalan di port `8080` dan melakukan auto-migration tabel database relasional.

5. **Jalankan Unit Test & Concurrency Test:**
   ```bash
   go test -v ./...
   ```
   Unit test otomatis melakukan simulasi 10 request transfer paralel konkuren untuk menguji keamanan deadlock dan balapan kondisi (*race condition*).

---

## 3. Manual Testing Walkthrough (cURL)

### 1. Registrasi Akun User A & User B
```bash
# Register User A (Sender)
curl -i -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "usera@email.com", "password": "password123"}'
# Nomor Wallet: W-10001 (User ID 1)

# Register User B (Receiver)
curl -i -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "userb@email.com", "password": "password123"}'
# Nomor Wallet: W-10002 (User ID 2)
```

### 2. Login User A & Dapatkan JWT Token
```bash
curl -i -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "usera@email.com", "password": "password123"}'
```
*Salin token JWT User A.*

### 3. Top-up Saldo User A Rp 100.000 (Idempotency Key Protected)
Kirim request pertama:
```bash
curl -i -X POST http://localhost:8080/wallet/top-up \
  -H "Authorization: Bearer <USER_A_TOKEN>" \
  -H "X-Idempotency-Key: topup-a-100k" \
  -H "Content-Type: application/json" \
  -d '{"amount": 100000.00, "description": "Gaji bulanan"}'
```
*Gaji sukses ditambahkan Rp 100.000.*

Kirim request kedua dengan **X-Idempotency-Key yang sama** untuk menguji kegagalan duplikasi:
```bash
curl -i -X POST http://localhost:8080/wallet/top-up \
  -H "Authorization: Bearer <USER_A_TOKEN>" \
  -H "X-Idempotency-Key: topup-a-100k" \
  -H "Content-Type: application/json" \
  -d '{"amount": 100000.00, "description": "Gaji bulanan"}'
```
*Sistem harus mengembalikan respons yang ter-cache dengan header `X-Cache-Lookup: HIT` tanpa menambah saldo User A menjadi Rp 200.000.*

### 4. Periksa Saldo (Balance Cache Check)
```bash
curl -i http://localhost:8080/wallet/balance \
  -H "Authorization: Bearer <USER_A_TOKEN>"
```
*Saldo User A harus bernilai tepat `100000.00`.*

### 5. Transfer Saldo Rp 40.000 ke User B (`W-10002`)
```bash
curl -i -X POST http://localhost:8080/wallet/transfer \
  -H "Authorization: Bearer <USER_A_TOKEN>" \
  -H "X-Idempotency-Key: transfer-a-to-b-1" \
  -H "Content-Type: application/json" \
  -d '{"destination_wallet_number": "W-10002", "amount": 40000.00, "description": "Bayar jajan"}'
```
*Transfer sukses. Saldo User A berkurang menjadi Rp 60.000. Cache saldo User A & B ter-invalidate otomatis.*
   
### 6. Cek Riwayat Mutasi Rekening Ledger
```bash
curl -i http://localhost:8080/wallet/transactions \
  -H "Authorization: Bearer <USER_A_TOKEN>"
```
*Menampilkan record Top-up (+100.000) dan Transfer (-40.000).*
