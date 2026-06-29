# Database Design: Authentication Service

---

## 1. Entity-Relationship Diagram (ERD)

```mermaid
erDiagram
    USERS {
        uint id PK
        varchar email UK "Not Null"
        varchar password_hash "Not Null"
        boolean is_verified "Default: false, Not Null"
        timestamp created_at
        timestamp updated_at
    }
    ROLES {
        uint id PK
        varchar name UK "Not Null"
        varchar description
        timestamp created_at
        timestamp updated_at
    }
    PERMISSIONS {
        uint id PK
        varchar name UK "Not Null"
        varchar description
        timestamp created_at
        timestamp updated_at
    }
    REFRESH_TOKENS {
        uint id PK
        varchar token UK "Not Null"
        uint user_id FK "On Delete CASCADE"
        timestamp expires_at "Not Null"
        boolean is_revoked "Default: false, Not Null"
        varchar parent_token
        timestamp created_at
        timestamp updated_at
    }
    USER_ROLES {
        uint user_id FK
        uint role_id FK
    }
    ROLE_PERMISSIONS {
        uint role_id FK
        uint permission_id FK
    }
    VERIFICATION_TOKENS {
        uint id PK
        uint user_id FK "On Delete CASCADE"
        varchar token UK "Not Null"
        timestamp expires_at "Not Null"
        timestamp created_at
    }
    RESET_TOKENS {
        uint id PK
        uint user_id FK "On Delete CASCADE"
        varchar token UK "Not Null"
        timestamp expires_at "Not Null"
        timestamp created_at
    }

    USERS ||--o{ REFRESH_TOKENS : "owns"
    USERS ||--o{ VERIFICATION_TOKENS : "has"
    USERS ||--o{ RESET_TOKENS : "has"
    USERS }o--o{ USER_ROLES : "maps"
    ROLES }o--o{ USER_ROLES : "maps"
    ROLES }o--o{ ROLE_PERMISSIONS : "maps"
    PERMISSIONS }o--o{ ROLE_PERMISSIONS : "maps"
```

## 2. Join Tables Many-to-Many

1. **`user_roles`**: Menghubungkan user ke role. Setiap user diasumsikan memiliki minimal satu role (default: `customer`).
2. **`role_permissions`**: Menghubungkan role ke izin permission secara granular (contoh: `admin` memiliki `wallet:read`, `wallet:write`, sedangkan `customer` hanya memiliki `wallet:read`).

---

## 3. Database Indexes

Untuk memitigasi brute-force attack dan mempercepat validasi token:

```sql
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
```

**Justifikasi Indeks:**
Pengecekan eksistensi refresh token (`RotateRefreshToken`) dan kueri login (`GetUserByEmail`) wajib didukung indeks unik untuk menekan latensi kueri di bawah 1ms.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi ERD skema relasional dynamic RBAC dan session token indexes |
