# Deployment: URL Shortener Service

**Deployment status:** `Local only`

---

## 1. Deployment target

Layanan ini dikonfigurasi untuk berjalan di lingkungan pengembangan lokal menggunakan Docker Compose dan Docker container. Layanan ini belum di-deploy ke cloud production (seperti AWS, GCP, atau Heroku), namun arsitekturnya dirancang agar siap di-containerize.

## 2. Build process

Kami menyediakan berkas [Dockerfile](./Dockerfile) multi-stage untuk meminimalkan ukuran image akhir dan meminimalkan celah keamanan di runtime.

Untuk membuat image Docker secara manual:
```bash
docker build -t url-shortener:latest .
```

**Penjelasan Tahapan Build:**
- **Stage 1 (Builder):** Menggunakan `golang:1.21-alpine` untuk compile aplikasi. Dependensi diunduh terlebih dahulu menggunakan cache layer. Kompilasi dinonaktifkan untuk CGO (`CGO_ENABLED=0`) agar menghasilkan binary statis yang tidak bergantung pada dynamic linker library C host.
- **Stage 2 (Runtime):** Menggunakan image minimalis `alpine:3.18`. Kami menyalin hanya file binary yang sudah dicompile dan file `.env.example` sebagai konfigurasi dasar. Ini menghasilkan ukuran image akhir yang sangat kecil (~20MB).

## 3. Configuration management

- **Lokal:** Konfigurasi dibaca dari file `.env` di direktori kerja aplikasi.
- **Produksi (Teoretis):** Di lingkungan produksi yang sesungguhnya (seperti Kubernetes atau AWS ECS), file `.env` tidak akan dipaketkan ke dalam image. Sebagai gantinya, variabel lingkungan akan dipasang (*injected*) menggunakan Secrets Manager (misal AWS Secrets Manager atau HashiCorp Vault) langsung ke container runtime environment.

## 4. Database migrations in deployment

Migrasi database dilakukan secara otomatis saat boot aplikasi menggunakan fungsi `AutoMigrate` GORM di [main.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/cmd/server/main.go).

*Catatan Produksi:* Pada sistem skala produksi yang sesungguhnya, migrasi otomatis saat boot sangat berisiko karena dapat menyebabkan lock database jika ada beberapa instance kontainer berjalan secara simultan (rolling updates). Proses produksi harus memisahkan langkah migrasi skema ke dalam job pipeline terpisah (CI/CD pre-deploy job) sebelum kontainer aplikasi baru dijalankan.

## 5. Health checks and readiness

Layanan menyediakan endpoint `/health` khusus untuk dimanfaatkan oleh *container orchestrator* (seperti Kubernetes Liveness/Readiness Probe atau Docker Compose Healthcheck).

```yaml
# Contoh integrasi di docker-compose (atau k8s deployment)
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 10s
  timeout: 5s
  retries: 3
```

## 6. Rollback strategy

Jika terjadi kegagalan setelah deployment kontainer baru:
1. Segera lakukan deploy ulang menggunakan image tag versi sebelumnya yang stabil (misal `url-shortener:v1.0.0` diganti ke `url-shortener:v0.9.0`).
2. Hindari melakukan rollback skema database otomatis jika skema database bersifat *backward-compatible*. Perubahan skema merusak harus diperbaiki dengan taktik migrasi maju (*forward migration*).

## 7. What real production would add

Untuk mematangkan sistem ke tingkat produksi sungguhan, kami menunda hal-hal berikut ke dalam [FUTURE-IMPROVEMENTS.md](./FUTURE-IMPROVEMENTS.md):
- **CI/CD Pipeline:** Otomasi build image Docker dan push ke Registry (seperti Docker Hub atau AWS ECR) setelah testing lulus di GitHub Actions.
- **Database Migration Tool:** Menggunakan alat migrasi mandiri (seperti `golang-migrate`) daripada GORM AutoMigrate guna memiliki kontrol versi migrasi skema (up/down).
- **SSL/TLS Termination:** Memasang Reverse Proxy (seperti Nginx atau Caddy) atau API Gateway di depan aplikasi untuk menangani enkripsi HTTPS.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen deployment lokal menggunakan Dockerfile multi-stage |
