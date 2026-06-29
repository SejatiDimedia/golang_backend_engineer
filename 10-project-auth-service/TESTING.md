# Testing Report: Authentication Service

---

## 1. Testing Strategy

Pengujian difokuskan pada keandalan sistem kriptografi asimetris RSA dan ketangguhan deteksi serangan replay pada session management:

1. **RSA Signer & Generator Test (`internal/utils/keys_test.go`):**
   - Menguji pembuatan key pair PEM RSA 2048-bit secara dinamis.
   - Memverifikasi generator tidak menimpa kunci yang sudah ada.
   - Memvalidasi token JWT RS256 yang dibuat dengan private key sukses diverifikasi secara offline menggunakan public key.
2. **RTR & RBAC Transactional Test (`internal/service/auth_test.go`):**
   - Menggunakan database **SQLite in-memory** (`file::memory:?cache=shared`) untuk menguji transaction behavior relasional yang sesungguhnya.
   - Menguji registrasi (dengan default role `customer`), verifikasi email token, dan login aman.
   - **Simulasi Replay Attack:** Melakukan rotasi refresh token, mengirim ulang refresh token lama yang sudah mati, lalu memverifikasi sistem mendeteksi fraud, menolak request (`401 Unauthorized`), dan mencabut seluruh sesi token user terkait secara massal di database.
   - **RBAC Join Query Verification:** Membuat dynamic roles & permissions, memetakannya ke user, dan memeriksa SQL JOIN query memuat permissions yang akurat.

---

## 2. Test Execution Command

Untuk mengeksekusi suite unit testing:
```bash
go test -v ./...
```

### Hasil Test Suites (PASS):
```
=== RUN   TestAuthService_RegisterAndLogin
--- PASS: TestAuthService_RegisterAndLogin (0.33s)
=== RUN   TestAuthService_RefreshTokenRotation_ReplayAttack
--- PASS: TestAuthService_RefreshTokenRotation_ReplayAttack (0.20s)
=== RUN   TestAuthService_RBAC_Query
--- PASS: TestAuthService_RBAC_Query (0.22s)
=== RUN   TestEnsureRSAKeys
--- PASS: TestEnsureRSAKeys (0.03s)
=== RUN   TestTokenManager_GenerateAndValidate
--- PASS: TestTokenManager_GenerateAndValidate (0.03s)
PASS
```

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen laporan pengujian database SQLite in-memory dan replay attack |
