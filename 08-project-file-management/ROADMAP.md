# Roadmap: File Management Service

**Status:** `Planning`

Dokumen ini memetakan urutan pembangunan fitur pada File Management Service untuk melatih integrasi MinIO Object Storage, streaming upload/download, dan penyajian S3 Presigned URLs.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Go module, docker-compose PostgreSQL 15 + MinIO, config load, dan utilitas JWT helper. | — | `Planned` |
| 2 — MinIO Bootstrapper | Inisialisasi MinIO Go SDK client, validasi koneksi, dan pembuatan bucket target `user-files` secara otomatis saat server booting. | Phase 1 | `Planned` |
| 3 — User & File Metadata Domain | Entitas `User` & `File` (PostgreSQL), login/register JWT ad-hoc, dan repository metadata berkas. | Phase 2 | `Planned` |
| 4 — Multipart Upload (Core) | REST API `POST /files/upload` menerima multipart form-data, validasi size max 10MB & MIME tipe, upload fisik ke MinIO, dan simpan metadata DB. | Phase 3 | `Planned` |
| 5 — Presigned URL Download | REST API `GET /files/:id/download` terproteksi JWT, menghasilkan URL bertanda tangan S3 dengan masa aktif 15 menit. | Phase 4 | `Planned` |
| 6 — Direct Server Streaming | REST API `GET /files/:id/view` membaca byte stream MinIO privat dan menyalurkannya langsung ke socket HTTP response menggunakan `io.Copy`. | Phase 5 | `Planned` |
| 7 — File Deletion | REST API `DELETE /files/:id` menghapus metadata di PostgreSQL dan memicu penghapusan objek fisik di MinIO bucket. | Phase 6 | `Planned` |
| 8 — Hardening & Testing | Unit test validasi file payload (size & format), dan integrasi mock MinIO client test suite. | Phase 7 | `Planned` |
| 9 — Deployment & Docs | Dockerfile multi-stage, perbaruan root README, dan pemenuhan 12 dokumen standar rilis. | Phase 8 | `Planned` |

## 2. Feature breakdown

| Feature | PRD ref | Build order reason |
|---|---|---|
| MinIO Initialization | — | Konektor inti. Harus terhubung sukses sebelum endpoint upload ditulis. |
| User Authentication | FR-1 | Membatasi hak akses upload/download hanya bagi user terdaftar. |
| Multipart Upload | FR-2 | Menerima berkas dari HTTP client, menyaring ukuran (<10MB), dan melemparnya ke Object Storage. |
| Metadata Tracking | FR-3 | Merekam data nama asli berkas & MIME tipe ke database relasional PostgreSQL. |
| Presigned Download | FR-5 | Membuat URL unduhan privat aman yang dikelola MinIO. |
| Streaming View | FR-6 | Memfasilitasi unduhan direct pipe stream untuk keperluan tag image HTML. |
| Delete File | FR-7 | Melakukan pembersihan data master di database dan berkas fisik di storage. |
| Health Check | — | Menampilkan keaktifan server, PostgreSQL, dan MinIO client. |

## 3. Concepts this project is exercising

- **Object Storage vs Block Storage:** Memahami keuntungan memisahkan file fisik ke S3 storage terpisah.
- **Multipart/Form-Data Parsing:** Menerima berkas stream di Gin menggunakan `c.FormFile("file")`.
- **S3 Presigned Signature:** Membuat url temporer bertanda tangan HMAC.
- **I/O Streaming (io.Copy):** Menyajikan berkas besar dengan konsumsi memori backend $O(1)$ RAM lewat pipe stream.
- **Rollback Data Terdistribusi:** Menghapus metadata jika upload fisik ke storage gagal (*compensating transaction*).

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi roadmap proyek File Management Service |
