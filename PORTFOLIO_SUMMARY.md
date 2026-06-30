# Portfolio Summary & Technical Interview Guide

Dokumen ini merangkum portofolio arsitektur canggih dari 8 proyek Go Backend yang telah diselesaikan pada ekosistem platform backend Anda. Setiap proyek dirancang untuk memenuhi standar production-grade, mengatasi tantangan konkurensi, keamanan, dan skalabilitas sistem terdistribusi.

---

## 1. Ringkasan Proyek & Highlight Arsitektur

### Proyek 1 — Fundamental Go
- **Konsep:** Pointer, Structs, Slices, Map, Control Flow.
- **Key Highlight:** Pemahaman dasar alokasi memori heap vs stack di Go, idiom penulisan Go, error handling sebagai return value, dan interface terkecil.

### Proyek 2 — REST API: In-Memory
- **Konsep:** Clean Architecture (Handler → Service → Repository), dependency injection manual, Gin routing, in-memory concurrency protection (`sync.RWMutex`).
- **Key Highlight:** Menangani race condition saat operasi baca/tulis data secara konkuren menggunakan mutex lock.

### Proyek 3 — REST API: PostgreSQL & GORM (Booking System)
- **Konsep:** Relational database migrations, dynamic filters, transaction validation.
- **Key Highlight:** Penerapan transaksi database GORM atomik untuk booking ketersediaan seat guna menghindari double-booking.

### Proyek 4 — Digital Wallet API (High Concurrency & Ledger)
- **Konsep:** Distributed Lock Manager (Redis `SET NX PX` + atomik Lua Script unlock), Double-Entry Ledger database schema, Idempotency key middleware, Balance Caching (Cache-aside pattern).
- **Key Highlight:** Mengatasi race condition transfer balance antar dompet secara konkuren menggunakan Distributed Lock terurut untuk menghindari deadlock. Menghindari pembebanan database relasional via cache-aside.

### Proyek 5 — File Management Service (Streaming & Storage)
- **Konsep:** Multipart upload streaming, MinIO Object Storage, data chunking, metadata synchronization.
- **Key Highlight:** Melakukan streaming data biner langsung ke object storage tanpa memuat seluruh berkas ke RAM server (latensi & RAM rendah).

### Proyek 6 — Asynchronous Notification Service (Message Queue)
- **Konsep:** Hand-rolled Redis Message Queue (`LPUSH`/`BRPOP`), Atomic Lua Scheduler, Worker Pool Concurrency Daemon, Exponential Backoff retry strategy, Dead-Letter Queue (DLQ).
- **Key Highlight:** Memproses tugas pengiriman secara asinkron non-blocking di latar belakang. Polling scheduler terdistribusi aman dari race condition multi-node disangga skrip Lua atomik di Redis.

### Proyek 7 — Centralized Authentication Service (IdP & Security)
- **Konsep:** Cryptographic RSA 2048-bit PEM keypair auto-generator, Asymmetric RS256 JWT, Refresh Token Rotation (RTR) dengan PostgreSQL row locking (`SELECT ... FOR UPDATE`), Replay Attack Prevention, Native relational dynamic RBAC tables.
- **Key Highlight:** RTR mendeteksi replay attack dan otomatis mencabut seluruh sesi token aktif user terkait (mass revocation) dengan commit transaction logic. Offline verification JWT RS256 downstream.

### Proyek 8 — AI Prompt Management API (Multi-Tenant & Compiler)
- **Konsep:** Dual-auth routing (JWT RS256 offline verifier vs Stripe-like API Key `prompt_live_...`), SHA-256 key hashing, Redis cache-aside API Key validation, Prompt Version full snapshots, Regex template compiler, Async buffered channel logging worker.
- **Key Highlight:** Kompilasi prompt template regex asinkron log analytics menggunakan buffered Go channel (size 1000) dan background Goroutine worker daemon.

---

## 2. Kisi-Kisi Pertanyaan & Jawaban Wawancara Teknis (Interview Guide)

### Q-1: Mengapa Anda memilih algoritma asimetris RS256 untuk JWT daripada HS256 di arsitektur microservices Anda? (Project 7 & 8)
- **Jawaban:** 
  "HS256 menggunakan kunci simetris (*Shared Secret*), artinya Auth Service dan setiap downstream service harus memegang kunci rahasia yang sama untuk memvalidasi token. Jika satu downstream service disusupi peretas, seluruh kunci keamanan platform bocor. 
  Sedangkan **RS256** menggunakan kriptografi asimetris (*Private/Public Key Pair*). Hanya Auth Service yang memegang *Private Key* untuk menandatangani token, sementara downstream services hanya membutuhkan *Public Key* untuk memverifikasi token secara lokal. Hal ini memungkinkan **Offline Verification** dengan keamanan tinggi karena downstream services tidak bisa menerbitkan token palsu meskipun public key-nya tersebar."

### Q-2: Bagaimana Anda menangani race condition pada saldo dompet digital saat 100 request transfer masuk secara konkuren ke user yang sama? (Project 4)
- **Jawaban:**
  "Saya mengimplementasikan **Distributed Lock Manager (DLM)** menggunakan Redis. Sebelum mutasi saldo dilakukan, server harus mengamankan lock untuk pengirim dan penerima di Redis menggunakan command `SET key value NX PX` (atomik). Untuk menghindari **Deadlock**, lock selalu diakuisisi dengan urutan ID terkecil terlebih dahulu (misal mengunci ID 3 baru mengunci ID 7). Jika lock gagal didapatkan dalam waktu timeout, transaksi dibatalkan. Setelah mutasi selesai, pelepasan lock dilakukan secara atomik menggunakan **Lua Script** agar aman dari race condition penghapusan lock milik request lain."

### Q-3: Apa itu Refresh Token Rotation (RTR) dan bagaimana sistem Anda mendeteksi jika token tersebut dicuri peretas? (Project 7)
- **Jawaban:**
  "RTR adalah mekanisme keamanan di mana setiap kali client menggunakan refresh token untuk memperbarui access token, server akan menerbitkan pasangan token baru dan me-revoke (membatalkan) token lama.
  Di database, saya melacak silsilah token (`ParentToken`) dan flag status `is_revoked`. Jika peretas mencuri token lama yang sudah mati dan mencoba mengirimkannya ulang (serangan replay), kueri database menggunakan row-locking `SELECT ... FOR UPDATE` mendeteksinya. Begitu terdeteksi `IsRevoked == true` pada token yang dikirim, sistem secara otomatis membatalkan seluruh sesi aktif (`is_revoked = true`) milik user tersebut secara massal. Ini memaksa peretas dan user asli logout bersamaan agar peretas kehilangan akses."

### Q-4: Bagaimana Anda mendesain antrean pesan (Message Queue) asinkron tanpa menggunakan library eksternal berat seperti RabbitMQ/Kafka? (Project 6)
- **Jawaban:**
  "Saya membangun Message Queue asinkron menggunakan struktur data **Redis List** dan **Sorted Set**. Tugas baru dimasukkan ke queue via `LPUSH` dan diproses oleh background worker pool menggunakan `BRPOP` (blocking pop) untuk menghindari polling CPU yang sia-sia.
  Untuk tugas terjadwal (delayed tasks), pesan ditaruh di Sorted Set Redis (`ZADD`) dengan score berupa timestamp eksekusi. Poller scheduler di latar belakang mengevaluasi pesan yang jatuh tempo menggunakan **Lua Script** atomik agar tidak terjadi duplikasi polling jika server dijalankan dalam konfigurasi multi-node (horizontal scaling)."

### Q-5: Bagaimana Anda memastikan logging analitik tidak memperlambat latensi utama respons API Prompt Compiler Anda? (Project 8)
- **Jawaban:**
  "Saya menggunakan pola **Asynchronous Analytics Logging**. Alih-alih menulis log secara sinkron ke database PostgreSQL saat request masuk (yang menambah overhead jaringan), server mengirim log data ke **Buffered Go Channel** berukuran 1000 secara non-blocking:
  ```go
  select {
  case s.ch <- logEntry:
  default:
      log.Println("Buffer full, dropping log")
  }
  ```
  Di latar belakang, sebuah daemon worker Goroutine berjalan membaca channel tersebut dan menuliskan data secara asinkron ke database. Hal ini memastikan performa compiler utama tetap sangat cepat (di bawah 15ms) meskipun database PostgreSQL sedang sibuk."
