# Deployment Guide: Authentication Service

---

## 1. Local containerization (Docker Compose)

Layanan dikemas menggunakan multi-stage Dockerfile yang ringan berbasis Alpine:

```dockerfile
# Stage 1: Build
FROM golang:1.20-alpine AS builder
...
# Stage 2: Runtime
FROM alpine:3.18
...
```

Untuk membangun dan menyalakan kontainer lokal:
```bash
docker-compose up --build -d
```

Ini akan menginisiasi:
1. `auth_db_container` (PostgreSQL 15) di port `5432`
2. `auth_redis_container` (Redis 7) di port `6379`
3. Server Auth Service Go di port `8081` (jika dimasukkan ke compose file, opsional).

---

## 2. Production key management

1. **RSA Key Pair Security:**
   - **PENTING:** Jangan pernah menaruh file `private.key` di repositori Git. Folder `certs/` sudah disaring lewat `.gitignore`.
   - Di cloud production (seperti AWS/GCP), RSA Private Key harus dimuat secara aman dari managed secret services (AWS Secrets Manager / GCP Secret Manager / HashiCorp Vault) dan disuntikkan ke runtime kontainer via environment variables / file mounts, alih-alih di-generate dinamis secara lokal saat boot.
2. **Offline Verification Distribution:**
   - Downstream services hanya membutuhkan berkas `public.key` untuk memvalidasi token JWT secara mandiri.
   - Distribusikan `public.key` ke downstream pod Kubernetes sebagai ConfigMap secara terpusat untuk mempermudah rotasi kunci publik di masa depan.
3. **Database Transaction Limits:**
   - Karena alur RTR (`RotateRefreshToken`) menggunakan query lock `SELECT ... FOR UPDATE` di PostgreSQL, pastikan koneksi DB pool dikonfigurasikan dengan batas timeout yang wajar guna menghindari penumpukan transaksi terkunci (*deadlocks*) saat beban traffic tinggi.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi panduan docker deployment dan pengamanan RSA keypair |
