# API Documentation: Booking Management System

---

## 1. Authentication Endpoints

### 1. Register Account
Mendaftarkan pengguna baru (customer atau admin).
- **HTTP Method:** `POST`
- **Path:** `/register`
- **Request Body:**
  ```json
  {
    "email": "user@email.com",
    "password": "secretpassword",
    "role": "customer"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 1,
    "email": "user@email.com",
    "role": "customer"
  }
  ```

### 2. Login (JWT Generation)
Masuk akun untuk memperoleh token otentikasi.
- **HTTP Method:** `POST`
- **Path:** `/login`
- **Request Body:**
  ```json
  {
    "email": "user@email.com",
    "password": "secretpassword"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

---

## 2. Desk (Aset) Endpoints

Seluruh endpoint di bawah ini memerlukan header `Authorization: Bearer <token>`.

### 1. Get Active Desks
Melihat daftar meja/ruangan coworking yang aktif (bisa diakses customer).
- **HTTP Method:** `GET`
- **Path:** `/desks`
- **Response (200 OK):**
  ```json
  [
    {
      "id": 1,
      "name": "Meja Kerja A1",
      "type": "hot-desk",
      "is_active": true,
      "created_at": "2026-06-29T13:00:00Z",
      "updated_at": "2026-06-29T13:00:00Z"
    }
  ]
  ```

### 2. Create Desk (Admin Only)
- **HTTP Method:** `POST`
- **Path:** `/desks`
- **Request Body:**
  ```json
  {
    "name": "Meeting Room 1",
    "type": "meeting-room"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 2,
    "name": "Meeting Room 1",
    "type": "meeting-room",
    "is_active": true,
    "created_at": "2026-06-29T13:45:00Z",
    "updated_at": "2026-06-29T13:45:00Z"
  }
  ```

---

## 3. Booking (Pemesanan) Endpoints

Seluruh endpoint di bawah ini memerlukan header `Authorization: Bearer <token>`.

### 1. Create Booking
Membuat pemesanan baru. Rentang waktu wajib dikirim dalam format ISO-8601/RFC3339 dengan info offset timezone.
- **HTTP Method:** `POST`
- **Path:** `/bookings`
- **Request Body:**
  ```json
  {
    "desk_id": 1,
    "start_time": "2026-06-29T14:00:00+07:00",
    "end_time": "2026-06-29T16:00:00+07:00"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 1,
    "user_id": 1,
    "desk_id": 1,
    "start_time": "2026-06-29T07:00:00Z",
    "end_time": "2026-06-29T09:00:00Z",
    "status": "CONFIRMED",
    "created_at": "2026-06-29T13:46:00Z",
    "desk": {
      "id": 1,
      "name": "Meja Kerja A1",
      "type": "hot-desk"
    }
  }
  ```

### 2. List Bookings
Melihat riwayat booking. Jika pemanggil adalah Admin, sistem mengembalikan seluruh riwayat pemesanan di coworking space. Jika pemanggil adalah Customer, hanya menampilkan miliknya saja.
- **HTTP Method:** `GET`
- **Path:** `/bookings`
- **Response (200 OK):**
  ```json
  [
    {
      "id": 1,
      "user_id": 1,
      "desk_id": 1,
      "start_time": "2026-06-29T07:00:00Z",
      "end_time": "2026-06-29T09:00:00Z",
      "status": "CONFIRMED",
      "created_at": "2026-06-29T13:46:00Z"
    }
  ]
  ```

### 3. Cancel Booking
Membatalkan pemesanan aktif. Hanya diizinkan jika waktu mulai terpaut minimal 2 jam dari sekarang (kecuali dilakukan oleh admin yang memiliki hak bypass window check).
- **HTTP Method:** `POST`
- **Path:** `/bookings/:id/cancel`
- **Response (200 OK):**
  ```json
  {
    "message": "booking cancelled successfully"
  }
  ```
