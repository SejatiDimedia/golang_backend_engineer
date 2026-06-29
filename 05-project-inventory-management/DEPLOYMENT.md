# Deployment: Inventory Management API

**Deployment status:** `Local only`

---

## 1. Deployment target

Layanan Inventory Management API saat ini dikonfigurasi untuk dijalankan pada infrastruktur pengembangan lokal menggunakan Docker Compose. Aplikasi ini siap di-package dan dijalankan pada runtime platform containerized seperti Google Cloud Run, AWS ECS, atau Kubernetes.

## 2. Build process

Kami menyediakan [Dockerfile](./Dockerfile) berbasis multi-stage build:

1. **Stage 1 (Builder):** Menggunakan base image `golang:1.21-alpine`. Berkas dependensi `go.mod` dan `go.sum` disalin terlebih dahulu agar di-cache. Binary dikompilasi statis dengan parameter `CGO_ENABLED=0 GOOS=linux`.
2. **Stage 2 (Runtime):** Menggunakan base image ultra-ringkas `alpine:3.18`. Menyalin hanya berkas binary `inventory-api` dan file `.env.example`. 

**Command pembuatan Image:**
```bash
docker build -t inventory-api:latest .
```

## 3. Configuration management

- **Lokal:** Dibaca dari file `.env` di direktori kerja kontainer.
- **Produksi (Teoretis):** Dibaca langsung dari variabel lingkungan (Env) yang dimasukkan oleh platform container orchestrator menggunakan Secrets Vault (seperti AWS Secrets Manager) tanpa menyisipkan berkas `.env` fisik ke dalam image kontainer.

## 4. Database migrations in deployment

Migrasi skema database dijalankan secara otomatis saat kontainer aplikasi dijalankan lewat `db.AutoMigrate` pada [main.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/cmd/server/main.go).

*Analisis Risiko Produksi:* Di produksi riil dengan beban tinggi, inisiasi `AutoMigrate` otomatis saat startup kontainer berisiko tinggi. Jika beberapa instance kontainer baru berjalan paralel secara simultan (misal autoscaling), kontainer-kontainer tersebut dapat berbut melakukan perubahan skema DDL yang mengakibatkan kegagalan lock database relasional. Untuk produksi, proses migrasi DDL harus dilepas dari startup server dan dikelola dalam CI/CD deploy job terpisah menggunakan script migration versioning (misal `golang-migrate`).

## 5. Health checks and readiness

Layanan menyediakan REST API endpoint `/health` untuk memantau status kontainer. Probe eksternal (seperti Kubernetes Liveness Probe) dapat mengonsumsi endpoint ini secara berkala:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 15
```

## 6. Rollback strategy

Apabila deployment kontainer baru mengalami malafungsi:
1. Kembalikan versi image kontainer ke versi tag stabil sebelumnya (seperti `inventory-api:v1.1.0` ke `inventory-api:v1.0.9`).
2. Jangan melakukan rollback skema DDL secara otomatis (seperti menghapus kolom/tabel baru) jika skema database bersifat *backward-compatible*, untuk menghindari kehilangan data transaksi yang sudah terjadi di DB selama rentang rilis gagal.

## 7. What real production would add

Penyempurnaan infrastruktur ke tingkat produksi dirinci di [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md):
- **CI/CD Pipeline:** Otomasi kompilasi image, running tests, dan deployment kontainer via GitHub Actions.
- **Database Migration Tooling:** Memisahkan proses inisiasi skema database dari startup server menggunakan `golang-migrate`.
- **API Gateway & HTTPS:** Menaruh proxy server (seperti Nginx atau cloud load balancer) untuk SSL termination di depan HTTP server Gin.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen deployment kontainer menggunakan Dockerfile multi-stage |
