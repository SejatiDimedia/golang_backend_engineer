# Deployment Guide: AI Prompt Management API

---

## 1. Multi-Stage Docker Build

Layanan dipaketkan menggunakan multi-stage Dockerfile berbasis Alpine:
```bash
docker-compose up --build -d
```

Compose file akan menyiapkan:
1. `prompt_db_container` (PostgreSQL 15) di port `5432`
2. `prompt_redis_container` (Redis 7) di port `6379`

---

## 2. Public Key Secret Injection

Karena layanan ini menggunakan verifikasi token JWT secara offline, berkas kunci publik (`public.key`) dari Auth Service harus dimasukkan secara aman:
1. **Kubernetes Environment:**
   - Simpan `public.key` milik Auth Service sebagai ConfigMap di Kubernetes cluster.
   - Mount ConfigMap tersebut sebagai file ke dalam pod AI Prompt Management API pada folder path `/app/certs/public.key`.
2. **Cloud Volume Mounts:**
   - Gunakan ECS/EKS volume mounts atau managed secret configurations untuk memisahkan file kunci dari container image statis. Jangan pernah memasukkan file `.key` langsung ke dalam Docker image saat proses building.
