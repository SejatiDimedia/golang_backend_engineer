# API Documentation: File Management Service

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

## 2. File Management Endpoints

Seluruh endpoint di bawah ini memerlukan header `Authorization: Bearer <token>`.

### 1. Upload File (Multipart Form)
Mengunggah berkas fisik.
- **HTTP Method:** `POST`
- **Path:** `/files/upload`
- **Content-Type:** `multipart/form-data`
- **Form Data Fields:**
  - `file`: (Binary File Payload)
- **Response (201 Created):**
  ```json
  {
    "id": 1,
    "user_id": 1,
    "file_name": "photo.jpg",
    "file_size": 102432,
    "content_type": "image/jpeg",
    "object_key": "user_1/1782713765075580000_photo.jpg",
    "status": "SUCCESS",
    "created_at": "2026-06-29T14:16:05Z"
  }
  ```

### 2. Get Uploaded Files List
Melihat list file sukses terunggah milik user.
- **HTTP Method:** `GET`
- **Path:** `/files`
- **Response (200 OK):**
  ```json
  [
    {
      "id": 1,
      "user_id": 1,
      "file_name": "photo.jpg",
      "file_size": 102432,
      "content_type": "image/jpeg",
      "object_key": "user_1/1782713765075580000_photo.jpg",
      "status": "SUCCESS",
      "created_at": "2026-06-29T14:16:05Z"
    }
  ]
  ```

### 3. Get Presigned Download URL
Mengambil tautan S3 bertanda tangan dengan TTL 15 menit.
- **HTTP Method:** `GET`
- **Path:** `/files/:id/download`
- **Response (200 OK):**
  ```json
  {
    "download_url": "http://localhost:9000/user-files/user_1/1782713765075580000_photo.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Expires=900..."
  }
  ```

### 4. Direct Server Streaming View
Membuka direct rendering berkas di browser (direct pipe stream).
- **HTTP Method:** `GET`
- **Path:** `/files/:id/view`
- **Response (200 OK):**
  - Mengembalikan data biner berkas fisik dengan header `Content-Type` yang sesuai (misal `image/jpeg`).

### 5. Delete File
Menghapus berkas dari DB dan Object Storage.
- **HTTP Method:** `DELETE`
- **Path:** `/files/:id`
- **Response (200 OK):**
  ```json
  {
    "message": "file deleted successfully"
  }
  ```
---
