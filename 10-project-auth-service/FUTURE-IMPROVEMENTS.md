# Future Improvements: Authentication Service

Rencana perbaikan dan peningkatan fitur untuk rilis Authentication Service berikutnya.

---

## 1. Migration Retrofit Plan for Downstream Services
- **Masalah Saat Ini:** Downstream services (Project 3: Booking System, Project 4: Digital Wallet, Project 6: Notification Service) masih menggunakan otentikasi JWT ad-hoc lokal.
- **Rencana Retrofit:**
  - Ambil file `certs/public.key` dari Auth Service dan pasang di folder konfigurasi downstream masing-masing.
  - Tulis ulang middleware otentikasi `AuthMiddleware` di downstream service agar memvalidasi token JWT secara offline menggunakan public key tersebut (menguji kecocokan signature RS256).
  - Ganti validasi role lokal di downstream service dengan membaca claim `role` dan `permissions` yang disematkan langsung di dalam claim JWT RS256.
  - Ini akan menghilangkan database user dan utility register/login ad-hoc di downstream services sepenuhnya!

## 2. OAuth2 / OpenID Connect (OIDC) Compliance
- **Masalah Saat Ini:** Auth flow login/register masih menggunakan format kustom non-standard OAuth2.
- **Rencana Solusi:** Kembangkan endpoint otentikasi agar kompatibel dengan spesifikasi RFC 6749 (Authorization Code Flow dengan PKCE untuk mobile client, dan Client Credentials Flow untuk server-to-server).

## 3. JWT Blacklisting with Redis
- **Masalah Saat Ini:** Saat user melakukan logout atau ganti password, Access Token JWT berumur pendek (15 menit) yang sudah dipegang client tetap valid secara offline hingga masa aktifnya habis.
- **Rencana Solusi:** Daftarkan JWT id (`jti`) token yang di-logout ke Redis blacklist cache dengan waktu kadaluarsa sesuai sisa umur token. Downstream middleware wajib memeriksa Redis blacklist sebelum menerima token.
