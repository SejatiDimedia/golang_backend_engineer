# API Documentation: Notification Service

---

## 1. Authentication Endpoints

### 1. Register Account
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

## 2. Notification Endpoints

Seluruh endpoint di bawah ini memerlukan header `Authorization: Bearer <token>`.

### 1. Create Notification (Instant / Scheduled)
Membuat permintaan pengiriman notifikasi baru. Jika parameter `send_at` tidak disertakan, notifikasi langsung diproses instan.
- **HTTP Method:** `POST`
- **Path:** `/notifications`
- **Request Body:**
  ```json
  {
    "type": "email",
    "target": "recipient@email.com",
    "content": "Halo, ini adalah email konfirmasi pemesanan coworking space.",
    "send_at": "2026-06-29T18:00:00Z" 
  }
  ```
- **Response (202 Accepted):**
  ```json
  {
    "message": "notification queued successfully",
    "notification_id": 1,
    "status": "PENDING",
    "send_at": "2026-06-29T18:00:00Z"
  }
  ```

### 2. Get Notification Status & Audit Logs
Melihat status terkini pengiriman notifikasi dan riwayat audit logs kegagalannya.
- **HTTP Method:** `GET`
- **Path:** `/notifications/:id`
- **Response (200 OK):**
  ```json
  {
    "id": 1,
    "type": "email",
    "target": "recipient@email.com",
    "content": "Halo, ini adalah email konfirmasi pemesanan coworking space.",
    "status": "SENT",
    "max_retries": 5,
    "attempt_count": 2,
    "send_at": "2026-06-29T18:00:00Z",
    "created_at": "2026-06-29T17:59:50Z",
    "updated_at": "2026-06-29T18:00:08Z",
    "logs": [
      {
        "id": 1,
        "notification_id": 1,
        "attempt": 1,
        "status": "FAILED",
        "error_message": "provider connection failure: server returned 503 service unavailable",
        "created_at": "2026-06-29T18:00:00Z"
      },
      {
        "id": 2,
        "notification_id": 1,
        "attempt": 2,
        "status": "SENT",
        "created_at": "2026-06-29T18:00:08Z"
      }
    ]
  }
  ```
---
