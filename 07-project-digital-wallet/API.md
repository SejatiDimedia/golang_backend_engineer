# API Documentation: Digital Wallet API

---

## 1. Authentication Endpoints

### 1. Register Account
Mendaftarkan pengguna baru. Sistem otomatis membuatkan wallet kosong dengan nomor unik (contoh: `W-10001`).
- **HTTP Method:** `POST`
- **Path:** `/register`
- **Request Body:**
  ```json
  {
    "email": "user@email.com",
    "password": "secretpassword"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 1,
    "email": "user@email.com"
  }
  ```

### 2. Login
Masuk akun untuk mendapatkan JWT token.
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

## 2. Wallet & Financial Endpoints

Seluruh endpoint di bawah ini memerlukan header `Authorization: Bearer <token>`.

### 1. Get Balance
Mendapatkan saldo dompet (ter-cache).
- **HTTP Method:** `GET`
- **Path:** `/wallet/balance`
- **Response (200 OK):**
  ```json
  {
    "balance": 100000.00
  }
  ```

### 2. Top-up Saldo (Idempotency Protected)
Menambah saldo dompet digital. Memerlukan header `X-Idempotency-Key` untuk proteksi request ganda.
- **HTTP Method:** `POST`
- **Path:** `/wallet/top-up`
- **Headers:**
  - `X-Idempotency-Key`: `unique-string-1`
- **Request Body:**
  ```json
  {
    "amount": 50000.00,
    "description": "Top-up via transfer bank"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": 1,
    "destination_wallet_id": 1,
    "amount": 50000.00,
    "type": "top-up",
    "description": "Top-up via transfer bank",
    "created_at": "2026-06-29T14:00:00Z"
  }
  ```

### 3. Withdraw Saldo (Idempotency Protected)
Menarik saldo dompet. Memerlukan header `X-Idempotency-Key`.
- **HTTP Method:** `POST`
- **Path:** `/wallet/withdraw`
- **Headers:**
  - `X-Idempotency-Key`: `unique-string-2`
- **Request Body:**
  ```json
  {
    "amount": 20000.00,
    "description": "Tarik tunai ATM"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": 2,
    "source_wallet_id": 1,
    "amount": 20000.00,
    "type": "withdraw",
    "description": "Tarik tunai ATM",
    "created_at": "2026-06-29T14:15:00Z"
  }
  ```

### 4. Transfer Saldo (Idempotency & Lock Protected)
Mengirim uang ke dompet pengguna lain berdasarkan nomor rekening wallet tujuan. Memerlukan header `X-Idempotency-Key`.
- **HTTP Method:** `POST`
- **Path:** `/wallet/transfer`
- **Headers:**
  - `X-Idempotency-Key`: `unique-string-3`
- **Request Body:**
  ```json
  {
    "destination_wallet_number": "W-10002",
    "amount": 30000.00,
    "description": "Bayar patungan makan siang"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": 3,
    "source_wallet_id": 1,
    "destination_wallet_id": 2,
    "amount": 30000.00,
    "type": "transfer",
    "description": "Bayar patungan makan siang",
    "created_at": "2026-06-29T14:20:00Z"
  }
  ```

### 5. Get Transaction History
Mendapatkan daftar seluruh riwayat mutasi debit/kredit ledger.
- **HTTP Method:** `GET`
- **Path:** `/wallet/transactions`
- **Response (200 OK):**
  ```json
  [
    {
      "id": 3,
      "source_wallet_id": 1,
      "destination_wallet_id": 2,
      "amount": 30000.00,
      "type": "transfer",
      "description": "Bayar patungan makan siang",
      "created_at": "2026-06-29T14:20:00Z"
    },
    {
      "id": 2,
      "source_wallet_id": 1,
      "amount": 20000.00,
      "type": "withdraw",
      "description": "Tarik tunai ATM",
      "created_at": "2026-06-29T14:15:00Z"
    }
  ]
  ```
