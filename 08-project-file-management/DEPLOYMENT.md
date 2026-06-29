# Deployment: File Management Service

---

## 1. Dockerization

Aplikasi ini menggunakan Dockerfile multi-stage untuk menjamin image runtime bersih, kecil, dan aman.

### Dockerfile
Dockerfile mengompilasi kode server ke binary executable static (`CGO_ENABLED=0 GOOS=linux`) di bawah golang alpine image, lalu memindahkannya ke alpine minimal runtime image (`alpine:3.18`).

### Build Image Manual
```bash
docker build -t file-management:latest .
```

---

## 2. Multi-Container Orchestration

Kami menggunakan `docker-compose.yml` untuk menjalankan orchestration backend secara terpadu:
- **postgres:** Relational Database PostgreSQL 15 untuk melacak metadata berkas.
- **minio:** Object Storage untuk menyimpan berkas fisik secara aman.

### Docker Compose Run
Jalankan stack kontainer di background:
```bash
docker-compose up -d --build
```

Matikan kontainer beserta volumes:
```bash
docker-compose down -v
```

---

## 3. Production Hardening Checklist

- [ ] Ubah environment variable `ENV` ke `production`. Mode debug Gin otomatis dinonaktifkan (`gin.SetMode(gin.ReleaseMode)`).
- [ ] Ganti `JWT_SECRET` dengan string acak (min. 32 karakter) yang diamankan menggunakan Vault/Secret Manager.
- [ ] Ganti kata sandi default PostgreSQL `DB_PASSWORD` dan MinIO root keys (`MINIO_ROOT_USER`/`MINIO_ROOT_PASSWORD`).
- [ ] Aktifkan SSL/TLS pada database postgres (`DB_SSLMODE=require`).
- [ ] Aktifkan SSL/TLS untuk MinIO endpoint (`MINIO_USE_SSL=true`) agar tanda tangan Signature Presigned URL dialihkan menggunakan enkripsi HTTPS.
- [ ] Konfigurasikan MinIO lifecycle policy (misal: auto-expire atau auto-cleanup untuk membersihkan folder temp lama di Object Storage).
