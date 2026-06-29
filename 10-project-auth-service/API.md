# API Documentation: Authentication Service

---

## 1. Authentication Endpoints

### 1. Register User
- **HTTP Method:** `POST`
- **Path:** `/auth/register`
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
    "message": "registration successful. check logs for email verification link",
    "id": 1,
    "email": "user@email.com"
  }
  ```

### 2. Verify Email
Mengaktifkan akun menggunakan token verifikasi email.
- **HTTP Method:** `GET`
- **Path:** `/auth/verify-email?token=<verification_token>`
- **Response (200 OK):**
  ```json
  {
    "message": "email verified successfully. you can now login"
  }
  ```

### 3. Login
- **HTTP Method:** `POST`
- **Path:** `/auth/login`
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
    "access_token": "eyJhbGciOiJSUzI1NiIs...",
    "refresh_token": "5aeb322b4604f5f4b9..."
  }
  ```

### 4. Refresh Token (RTR)
Mengajukan pasangan token baru menggunakan refresh token aktif.
- **HTTP Method:** `POST`
- **Path:** `/auth/refresh`
- **Request Body:**
  ```json
  {
    "refresh_token": "5aeb322b4604f5f4b9..."
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "access_token": "eyJhbGciOiJSUzI1NiIs...",
    "refresh_token": "de1265e32a7b5d6fc6..."
  }
  ```

### 5. Logout
Mencabut keaktifan refresh token.
- **HTTP Method:** `POST`
- **Path:** `/auth/logout`
- **Request Body:**
  ```json
  {
    "refresh_token": "de1265e32a7b5d6fc6..."
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "message": "successfully logged out"
  }
  ```

### 6. Token Introspection
Dipanggil oleh microservice downstream secara online.
- **HTTP Method:** `POST`
- **Path:** `/auth/introspect`
- **Request Body:**
  ```json
  {
    "token": "eyJhbGciOiJSUzI1NiIs..."
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "active": true,
    "user_id": 1,
    "email": "user@email.com",
    "role": "customer",
    "permissions": [
      "wallet:read"
    ]
  }
  ```

---

## 2. RBAC Management Endpoints (Admin only)

### 1. Create Role
- **HTTP Method:** `POST`
- **Path:** `/auth/rbac/roles`
- **Request Body:**
  ```json
  {
    "name": "admin",
    "description": "System Administrator"
  }
  ```

### 2. Create Permission
- **HTTP Method:** `POST`
- **Path:** `/auth/rbac/permissions`
- **Request Body:**
  ```json
  {
    "name": "wallet:write",
    "description": "Modify wallet balances"
  }
  ```

### 3. Assign Role to User
- **HTTP Method:** `POST`
- **Path:** `/auth/rbac/users/:id/roles`
- **Request Body:**
  ```json
  {
    "role_id": 1
  }
  ```

### 4. Assign Permission to Role
- **HTTP Method:** `POST`
- **Path:** `/auth/rbac/roles/:id/permissions`
- **Request Body:**
  ```json
  {
    "permission_id": 1
  }
  ```
