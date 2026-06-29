# Changelog: File Management Service

Format penulisan didasarkan pada [Keep a Changelog](https://keepachangelog.com/).

---

## [1.0.0] — 2026-06-29 — Initial Release

### Added
- **MinIO Object Storage Integration:** Integrasi SDK client MinIO (`minio-go`) untuk meng-offload penyimpanan fisik dari local disk.
- **Auto-Bucket Creation:** Pendeteksi otomatis bucket target `user-files` saat server Go booting.
- **Compensating Writes Rollback:** Logika transaksional yang otomatis menghapus metadata di PostgreSQL jika PutObject ke MinIO gagal.
- **Multipart Form Upload Handler:** Endpoint `POST /files/upload` yang melayani upload streaming direct memory (stateless RAM).
- **S3 Presigned URLs:** Endpoint `GET /files/:id/download` menghasilkan tautan bertanda tangan S3 dengan masa aktif 15 menit.
- **Direct Server Streaming:** Endpoint `GET /files/:id/view` mem-pipe byte stream MinIO langsung ke klien menggunakan `io.Copy`.
- **JWT Auth Middleware:** Middleware otentikasi local ad-hoc.
- **File Deletion:** Sinkronisasi pembersihan baris di DB relasional dan objek fisik di MinIO.
- **Unit Testing:** Unit test komprehensif menguji token JWT, file validation limits, dan compensating rollback.
- **Dockerization:** Berkas Dockerfile multi-stage dan docker-compose.yml (PostgreSQL + MinIO).
