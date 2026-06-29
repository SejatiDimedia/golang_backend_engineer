# PRD: Notification Service

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Dalam aplikasi web modern, pengiriman notifikasi seperti email, SMS, push notification, atau pemanggilan webhook sering kali membutuhkan waktu respon lambat karena bergantung pada layanan pihak ketiga (external providers). Jika pengiriman ini dilakukan secara sinkron dalam utas request HTTP utama, response time API akan melonjak drastis dan rentan terhadap kegagalan jika provider eksternal sedang mengalami downtime. Sistem membutuhkan layanan notifikasi asinkron yang andal menggunakan antrean pesan (message queue), pemrosesan latar belakang (worker pool), dan kebijakan retry otomatis saat terjadi gangguan jaringan temporer tanpa mengorbankan performa API utama.

## 2. Goals

- Menerima request pengiriman notifikasi (Email, Webhook, Push) via REST API secara instan dan memasukkannya ke antrean.
- Memisahkan proses pengiriman fisik dari utas HTTP utama menggunakan antrean asinkron (Redis Queue).
- Menggunakan worker pool asinkron untuk memproses antrean pesan secara paralel dan efisien.
- Mengimplementasikan retry mechanism otomatis dengan strategi **Exponential Backoff** dan limit maksimum percobaan jika provider eksternal gagal merespon.
- Melacak seluruh status pengiriman notifikasi (`PENDING`, `PROCESSING`, `SENT`, `FAILED`) dan audit log percobaan di database PostgreSQL.
- Menyediakan fungsionalitas pengiriman notifikasi terjadwal (*scheduled notifications*) menggunakan scheduler asinkron (Cron / Ticker).

## 3. Non-goals

- **Real SMTP/SMS/Push Gateway integration:** Kita tidak akan menggunakan akun berbayar (seperti SendGrid, Twilio, Firebase) untuk melayani pengiriman nyata. Sebagai gantinya, kita akan menulis adapter dummy/mock yang mencetak log pengiriman dan memiliki opsi simulasi *random failure rate* (misalnya kegagalan jaringan acak 30%) guna memicu dan menguji mekanisme retry backoff asinkron.
- **Rich Template Editor:** Sistem hanya menerima pesan teks mentah atau format JSON siap kirim tanpa mesin parser HTML/template dinamis di versi pertama.

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Client Services (e.g. Booking / Wallet APIs) | Mengirimkan konfirmasi booking atau bukti transaksi transfer secara asinkron ke pengguna. | Ratusan kali per menit |
| System Operator | Memantau antrean notifikasi gagal (dead-letter queue) dan menganalisis audit trail retry log untuk menemukan provider yang bermasalah. | Harian |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | Pengguna/Service dapat mengirimkan request notifikasi instan (Email, Webhook, Push) via REST API `POST /notifications` dengan proteksi token JWT ad-hoc. | Must |
| FR-2 | **Asynchronous Queueing:** Setiap request notifikasi langsung dimasukkan ke dalam antrean Redis dan API langsung merespon HTTP `202 Accepted` tanpa menunggu berkas terkirim secara fisik. | Must |
| FR-3 | **Worker Pool Processing:** Sekumpulan worker latar belakang secara paralel mengambil notifikasi dari Redis queue dan memproses pengiriman ke provider dummy. | Must |
| FR-4 | **Retry with Exponential Backoff:** Jika pengiriman gagal, worker otomatis mengantrekan kembali tugas dengan waktu tunda dinamis yang meningkat secara eksponensial (misal: delay = $2^{\text{attempt}} \times 2$ detik, maks 5 kali). | Must |
| FR-5 | **Metadata & Status Tracking:** Setiap notifikasi dicatat di PostgreSQL dengan status pengiriman yang selalu diperbarui beserta log error setiap percobaan gagal. | Must |
| FR-6 | **Scheduled Notifications:** Layanan menerima parameter `send_at` (timestamp UTC) dan menunda pengiriman hingga waktu yang ditentukan tiba menggunakan scheduler worker. | Must |
| FR-7 | **Dead-Letter Queue (DLQ):** Notifikasi yang telah melampaui batas maksimum retry (5 kali) otomatis dipindahkan ke status `FAILED` dan ditandai sebagai DLQ untuk inspeksi admin. | Should |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Performance | API endpoint `POST /notifications` harus merespon dalam waktu $<10\text{ms}$ karena hanya bertugas menulis payload ke Redis. |
| Security | REST API dilindungi JWT token ad-hoc. Koneksi Redis dan PostgreSQL diamankan dengan parameter konfigurasi environment variables. |
| Reliability | Jaminan **At-Least-Once Delivery**: notifikasi tidak boleh hilang dari antrean sebelum sukses dikirim atau ditandai gagal permanen (DLQ). |
| Concurrency | Worker pool harus dapat dikonfigurasikan jumlah worker paralelnya (`WORKER_CONCURRENCY=5`) untuk mengatur konsumsi memori dan batasan rate-limiting provider eksternal. |

## 7. Constraints

- **Teknologi:** Go, Redis (v7), PostgreSQL (v15), GORM, Gin, Docker & Docker Compose.
- **Redis Client:** Menggunakan `github.com/redis/go-redis/v9`.
- **Database:** PostgreSQL digunakan untuk persistensi data audit jangka panjang, sedangkan Redis digunakan khusus sebagai in-memory message queue berkecepatan tinggi.

## 8. Success criteria

- API `/notifications` mengembalikan status `202 Accepted` instan.
- Worker asinkron sukses memproses antrean pesan secara berurutan.
- Pengiriman yang gagal sukses melakukan retry otomatis dengan delay yang meningkat secara eksponensial.
- Notifikasi terjadwal ditunda secara tepat waktu dan diproses ketika waktu `send_at` terlampaui.
- Seluruh histori pengiriman terdata rapi di PostgreSQL.

## 9. Open questions

- **Redis Queue Strategy:** Memilih **Rekomendasi (Hand-Rolled)**. Kami membangun antrean pesan asinkron sendiri dari awal memanfaatkan primitive Redis:
  - Redis List (`LPUSH` / `BRPOP`) untuk instant message queue.
  - Redis Sorted Set (`ZADD` / `ZRANGEBYSCORE`) untuk melayani antrean terjadwal (*scheduled*) dan retry delay (*exponential backoff*).
  Hal ini bertujuan memaksimalkan pemahaman logika concurrency, worker pool, dan pemrosesan asinkron di Go.
- **Deadlock / Race Conditions pada Retry DB Write:** Menggunakan kueri SQL atomic update langsung per baris (berdasarkan `ID` unik notifikasi) tanpa memblokir seluruh tabel, guna menjamin kinerja PostgreSQL tetap aman di bawah beban worker paralel.


---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
