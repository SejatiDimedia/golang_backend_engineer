# Deployment: Digital Wallet API

---

## 1. Dockerization

Aplikasi ini menggunakan Dockerfile multi-stage untuk menjamin image runtime bersih, kecil, dan aman.

### Dockerfile
Dockerfile mengompilasi kode server ke binary executable static (`CGO_ENABLED=0 GOOS=linux`) di bawah golang alpine image, lalu memindahkannya ke alpine minimal runtime image (`alpine:3.18`).

### Build Image Manual
```bash
docker build -t digital-wallet:latest .
```

---

## 2. Multi-Container Orchestration

Kami menggunakan `docker-compose.yml` untuk menjalankan multi-container orchestration backend:
- **postgres:** Relational Database PostgreSQL 15 untuk penyimpanan transaksi ACID.
- **redis:** In-memory Cache & Distributed Lock manager.

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
- [ ] Ganti kata sandi default PostgreSQL `DB_PASSWORD` dan Redis `REDIS_PASSWORD` (pastikan Redis tidak terekspos tanpa password ke internet).
- [ ] Aktifkan SSL/TLS pada database postgres (`DB_SSLMODE=require`).
- [ ] Batasi resource limit CPU dan Memori pada kontainer Go app di dalam docker compose untuk mencegah serangan DDoS OOM.
