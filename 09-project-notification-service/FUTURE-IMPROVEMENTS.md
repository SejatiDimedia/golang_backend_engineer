# Future Improvements: Notification Service

Rencana perbaikan dan peningkatan fitur untuk rilis Notification Service berikutnya.

---

## 1. Webhook Signature Verification
- **Masalah Saat Ini:** Pengiriman webhook ke client dilakukan dengan HTTP POST payload mentah tanpa verifikasi integritas data.
- **Rencana Solusi:** Implementasikan HMAC-SHA256 signature verification. Server notifikasi menghasilkan signature berdasarkan secret key bersama dan menyisipkannya di header `X-Notification-Signature` agar client dapat memverifikasi keaslian pengirim.

## 2. Shared Single-Sign-On (SSO) JWT Integration
- **Masalah Saat Ini:** Auth token JWT masih dibuat secara ad-hoc lokal per project (berulang dari Project 5).
- **Rencana Solusi:** Sesuai roadmap, integrasikan verifikasi token dengan JWT signer pusat / OAuth2 provider dari UserService (Project 2/3) agar layanan notifikasi tidak perlu menyimpan register/login database users secara mandiri.

## 3. Real SMTP / SendGrid Adapters
- **Masalah Saat Ini:** Pengiriman email fisik masih menggunakan mock log simulation.
- **Rencana Solusi:** Tambahkan SMTP adapter nyata (memanfaatkan library standard `net/smtp`) atau REST API Integration ke provider eksternal berbayar seperti SendGrid/Mailgun dengan fallback dynamic routing.

## 4. Web Dashboard UI
- **Masalah Saat Ini:** Admin harus melakukan query sql manual ke PostgreSQL / API call untuk melihat status DLQ (notifikasi yang gagal permanen).
- **Rencana Solusi:** Sediakan dasbor web sederhana menggunakan React/Tailwind untuk menampilkan visualisasi antrean Redis, jumlah worker aktif, dan tabel daftar notifikasi `FAILED` beserta tombol trigger `Re-enqueue` manual.
