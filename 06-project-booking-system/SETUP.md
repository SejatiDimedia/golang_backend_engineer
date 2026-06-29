# Setup Guide: Booking Management System

---

## 1. Prerequisites

Pastikan komputer Anda telah terinstal:
- Go 1.20 atau versi di atasnya
- Docker & Docker Compose
- Alat penguji API seperti `curl` atau Postman

---

## 2. Local Installation

1. **Inisialisasi Konfigurasi:**
   Salin template variabel environment ke berkas `.env`:
   ```bash
   cp .env.example .env
   ```

2. **Menjalankan Database PostgreSQL:**
   Nyalakan container database relasional PostgreSQL:
   ```bash
   docker-compose up -d
   ```

3. **Inisiasi Database 'booking_db':**
   Jika container sudah menyala namun Anda mendapatkan error *database "booking_db" does not exist*, jalankan perintah pembuatan database manual berikut:
   ```bash
   PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE booking_db;"
   ```

4. **Menjalankan Server:**
   Jalankan server backend lokal:
   ```bash
   go run cmd/server/main.go
   ```
   Server akan berjalan secara default di port `8080` dan melakukan auto-migration tabel `users`, `desks`, dan `bookings`.

---

## 3. Manual Testing Walkthrough (cURL)

Anda dapat menguji siklus pemesanan (dan validasi bentrokan slot) menggunakan rangkaian cURL berikut:

### 1. Registrasi Akun Customer
```bash
curl -i -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "customer@email.com", "password": "password123", "role": "customer"}'
```

### 2. Login Akun Customer
```bash
curl -i -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "customer@email.com", "password": "password123"}'
```
*Salin nilai `"token"` dari JSON respons di atas untuk digunakan pada langkah selanjutnya.*

### 3. Registrasi & Login Akun Admin (Untuk Input Katalog Meja)
```bash
# Register Admin
curl -i -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@email.com", "password": "adminpassword", "role": "admin"}'

# Login Admin
curl -i -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@email.com", "password": "adminpassword"}'
```
*Salin Token Admin Anda.*

### 4. Membuat Meja Baru (Gunakan Token Admin)
```bash
curl -i -X POST http://localhost:8080/desks \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Hot Desk A1", "type": "hot-desk"}'
```
*Meja sukses dibuat dengan ID 1.*

### 5. Melakukan Booking Sukses (Gunakan Token Customer)
Lakukan pemesanan untuk tanggal besok jam 14:00 s/d 16:00 (misalnya):
```bash
curl -i -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer <CUSTOMER_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"desk_id": 1, "start_time": "2026-06-30T14:00:00+07:00", "end_time": "2026-06-30T16:00:00+07:00"}'
```

### 6. Menguji Overlap Bentrokan Waktu (Double-Booking)
Mencoba memesan meja yang sama di jam yang bersinggungan (15:00 s/d 17:00):
```bash
curl -i -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer <CUSTOMER_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"desk_id": 1, "start_time": "2026-06-30T15:00:00+07:00", "end_time": "2026-06-30T17:00:00+07:00"}'
```
*Sistem harus mengembalikan status `400 Bad Request` dengan error `the room/desk is already booked for this time range`.*

### 7. Membatalkan Pemesanan (Cancel Booking)
```bash
curl -i -X POST http://localhost:8080/bookings/1/cancel \
  -H "Authorization: Bearer <CUSTOMER_TOKEN>"
```
*Status booking 1 akan terupdate menjadi CANCELLED.*
