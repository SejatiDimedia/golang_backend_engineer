# Deployment: Booking Management System

---

## 1. Dockerization

Aplikasi ini menggunakan Dockerfile multi-stage untuk menjamin image runtime bersih, kecil, dan aman.

### Dockerfile Breakdown
- **Build Stage:** Menggunakan golang alpine image versi ringan (`golang:1.20-alpine`) untuk mengunduh modul dependencies dan mengompilasi kode server ke binary executable static (`CGO_ENABLED=0 GOOS=linux`).
- **Run Stage:** Menggunakan alpine minimal image (`alpine:3.18`). Hanya berkas binary kompilasi `server` dan berkas environment `.env.example` (disalin sebagai default `.env`) yang dimasukkan ke image runtime akhir.

### Build Image Manual
```bash
docker build -t booking-system:latest .
```

---

## 2. Multi-Container Orchestration

Kami menyertakan berkas `docker-compose.yml` untuk mempermudah orkestrasi kontainer server backend dan database relasional PostgreSQL 15 di server staging/produksi lokal.

### Docker Compose Run
Jalankan seluruh stack aplikasi di background:
```bash
docker-compose up -d --build
```

Matikan seluruh kontainer beserta volumes datanya:
```bash
docker-compose down -v
```

---

## 3. Production Hardening Checklist

Sebelum mempublikasikan Booking System ke cloud (seperti GCP, AWS, DigitalOcean):
- [ ] Ubah environment variable `ENV` ke `production`. Mode debug Gin otomatis dinonaktifkan (`gin.SetMode(gin.ReleaseMode)`).
- [ ] Ganti `JWT_SECRET` dengan string acak (min. 32 karakter) yang diamankan menggunakan Vault/Secret Manager. Jangan biarkan bernilai default.
- [ ] Ganti kata sandi PostgreSQL `DB_PASSWORD` dan hilangkan port `5432:5432` dari publik (cukup biarkan PostgreSQL terekspos di jaringan internal docker compose).
- [ ] Aktifkan SSL/TLS pada database postgres (`DB_SSLMODE=require`).
- [ ] Batasi resource limit CPU dan Memori pada kontainer Go app di dalam docker compose untuk mencegah serangan DDoS OOM.
