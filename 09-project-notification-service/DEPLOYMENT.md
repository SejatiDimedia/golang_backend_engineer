# Deployment Guide: Notification Service

---

## 1. Local containerization (Docker Compose)

Layanan dikemas dengan multi-stage Dockerfile yang ringan berbasis Alpine:

```dockerfile
# Stage 1: Build
FROM golang:1.20-alpine AS builder
...
# Stage 2: Runtime
FROM alpine:3.18
...
```

Untuk menyalakan sistem utuh di lingkungan lokal menggunakan Docker Compose:
```bash
docker-compose up --build -d
```

Ini akan menginisiasi 3 kontainer aktif:
1. `notification_db` (PostgreSQL 15) di port `5432`
2. `notification_redis` (Redis 7) di port `6379`
3. Server utama Go di port `8080` (jika dimasukkan ke compose file, opsional).

---

## 2. Production infrastructure recommendation

Jika layanan ini dideploy ke server production cloud (seperti AWS, GCP, atau Kubernetes):

1. **Redis Management:**
   - Gunakan managed Redis service (misal AWS ElastiCache / GCP Memorystore) dengan fitur Multi-AZ enabled untuk mencegah kehilangan tugas antrean ketika satu node mati.
   - Set konfigurasi `maxmemory-policy` ke `noeviction` di Redis config. Hal ini krusial agar Redis tidak menghapus antrean pesan secara acak saat memori penuh.
2. **PostgreSQL Audit DB:**
   - Gunakan RDS PostgreSQL dengan replikasi write/read split. Indeks audit logs berukuran besar, sehingga pemisahan lalu lintas baca (`GET /notifications/:id`) ke read-replica akan mempercepat performa engine.
3. **Horizontal Pod Autoscaling (HPA):**
   - Layanan ini stateless. Anda bisa melipatgandakan replika pod server Go (misal menjadi 5 pods).
   - Berkat **atomic Lua script** pada poller scheduler, multi-pod server tidak akan memperebutkan atau menduplikasi pengiriman scheduled/retry task.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi panduan docker deployment dan best practices cloud architecture |
