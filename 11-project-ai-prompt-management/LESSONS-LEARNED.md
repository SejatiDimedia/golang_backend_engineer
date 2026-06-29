# Lessons Learned: AI Prompt Management API

Retrospektif pembelajaran dari arsitektur platform AI Prompt Management API.

---

## 1. Keandalan Otentikasi Ganda (Dual-Auth Middleware)
- **Temuan:** Menggabungkan admin dashboard web (melalui JWT RS256) dan server-to-server microservices (melalui API Key) dalam satu API gateway dapat memperumit routing middleware.
- **Pembelajaran:** Pemisahan route group secara eksplisit di Gin (`/api/v1` vs `/api/v1/client`) adalah desain terbaik. Ini mempermudah pemetaan middleware keamanan yang berbeda (offline JWT verifier vs Redis key cache) tanpa mencemari context request masing-masing klien.

## 2. Pemanfaatan Redis Cache-Aside untuk API Key Lookup
- **Temuan:** Di awal perancangan, otentikasi API Key mentah SHA-256 yang langsung menanyakan database relasional PostgreSQL di setiap request kompilasi prompt memicu overhead latensi.
- **Pembelajaran:** Meng-cache hash API Key yang valid di Redis dengan TTL 1 jam (`apikey:<hash> -> workspace_id`) memotong latensi otentikasi di bawah 2ms. Sinkronisasi cache di-manage secara rapi dengan menghapus Redis key (`RDB.Del`) ketika API Key di-revoke manual oleh admin di database.

## 3. Desain Asinkronous Buffered Channels untuk Logging Analytics
- **Temuan:** Logging data statistik pemakaian compiler (latensi, hit, token estimate) langsung ke PostgreSQL secara sinkron membebani throughput utama respons API.
- **Pembelajaran:** Menggunakan buffered channel ukuran 1000 dan background worker daemon Goroutine mengisolasi operasi I/O penulisan log database ke latar belakang. Klien menerima respons kompilasi prompt instan, sementara data analitik diproses secara non-blocking di latar belakang secara aman.
