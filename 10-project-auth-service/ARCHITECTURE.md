# Architecture: Authentication Service

---

## 1. Asymmetric Token Verification (Offline Validation)

Salah satu pilar utama arsitektur platform ini adalah meminimalkan overhead jaringan verifikasi token. Melalui asymmetric cryptography RS256, downstream services memverifikasi validitas token secara offline:

```mermaid
sequenceIndex
    Client ->> Auth_Service: POST /auth/login (password)
    Auth_Service ->> Auth_Service: Sign JWT using private.key
    Auth_Service -->> Client: Return Access Token (RS256)
    
    Note over Client, Downstream: Client requests to other services
    Client ->> Downstream: GET /wallet/balance (JWT Token)
    Note over Downstream: Downstream loads public.key locally
    Downstream ->> Downstream: Verify JWT Offline (O(1) CPU only)
    Downstream -->> Client: 200 OK Wallet Data
```

## 2. Refresh Token Rotation (RTR) & Replay Attack Prevention

Untuk mengamankan platform dari token refresh yang dicuri, RTR melacak pohon silsilah token (`ParentToken`) dan status pembatalannya (`IsRevoked`). Transaksi database menggunakan baris penguncian `SELECT ... FOR UPDATE` untuk mencegah race condition.

### Skenario Normal:
- User mengirim Refresh Token A.
- Server membatalkan Token A (`is_revoked = true`), menghasilkan Refresh Token B, dan mengembalikannya.

### Skenario Replay Attack (Peretas mencoba mengirim ulang Token A yang sudah hangus):

```mermaid
sequenceIndex
    Hacker ->> Auth_Service: POST /auth/refresh (Token A - already revoked!)
    Note over Auth_Service: DB Transaction: SELECT FOR UPDATE on Token A
    Auth_Service ->> Auth_Service: Detects Token A IsRevoked == true!
    Note over Auth_Service: CRITICAL: Security breach detected
    Auth_Service ->> PostgreSQL: Revoke ALL active sessions for this User ID
    PostgreSQL -->> Auth_Service: Sessions invalidated
    Auth_Service -->> Hacker: 401 Unauthorized (Force Logout Everywhere)
```

---

## 3. Directory structure

```
10-project-auth-service/
├── cmd/
│   └── server/
│       └── main.go         # Bootstrap & graceful shutdown
├── certs/                  # (.gitignore) Folder private & public keys
├── internal/
│   ├── config/             # Config parser loading env variables
│   ├── entity/             # relational GORM models (User, Role, Permission)
│   ├── handler/            # HTTP endpoints (Auth, Introspect, RBAC)
│   ├── middleware/         # Security headers
│   ├── repository/         # DB query operations with transaction locks
│   ├── service/            # Authentication & RTR business logic
│   └── utils/              # RSA generation tools and token managers
```

---

## 4. Key architectural decisions

- **Offline Offline JWT Verification:** Kami menolak database query terpusat untuk otentikasi di setiap request. Downstream service hanya membutuhkan `public.key` untuk memvalidasi user secara mandiri, menghemat miliaran database calls pada skala produksi.
- **RTR Replay Commit Strategy:** Untuk menjamin sesion pembatalan massal user tersimpan permanen di database saat serangan replay terdeteksi, GORM transaction *tidak di-rollback* melainkan di-commit (menyimpan status `is_revoked = true` untuk semua baris), baru kemudian di luar callback transaksi error `401` dikembalikan ke client.
