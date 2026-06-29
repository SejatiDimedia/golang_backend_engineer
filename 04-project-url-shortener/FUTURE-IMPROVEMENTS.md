# Future Improvements: URL Shortener Service

Dokumen ini mendokumentasikan batas cakupan teknis yang kami batasi secara sengaja di Project 1 demi tujuan pembelajaran yang efisien, serta memetakan rencana perbaikan di masa depan.

---

## Deferred by design (in scope for a future project)

| Item | Why deferred | Where it's actually addressed |
|---|---|---|
| **Autentikasi Pengguna (JWT Auth)** | URL Shortener saat ini bersifat publik. Mengembangkan sistem autentikasi di sini akan mengalihkan fokus dari dasar routing REST API. | [Project 3 (Booking System)](../06-project-booking-system/) dan [Project 7 (Auth Service)](../10-project-auth-service/) |
| **Caching dengan Redis** | Untuk menghindari kompleksitas integrasi Redis di proyek pertama, kami menyimpan data langsung ke PostgreSQL. | [Project 4 (Digital Wallet API)](../07-project-digital-wallet/) |
| **Background Processing (Queues)** | Pemrosesan asinkron untuk mencatat statistik klik kunjungan ditiadakan agar kami fokus pada database transaction yang aman di PostgreSQL. | [Project 6 (Notification Service)](../09-project-notification-service/) |

## Deferred due to scope (not a future project's job, just not done)

| Item | Why deferred | Would require |
|---|---|---|
| **Alat Migrasi Database Mandiri** | Penggunaan `AutoMigrate` GORM sangat cepat saat memulai, namun rentan dalam melacak sejarah rollback skema. | Pengenalan alat migrasi seperti `golang-migrate` dan pembuatan skrip `.sql` up/down secara terpisah. |
| **Validasi Keberadaan Tautan (Live Check)** | Saat ini kami memvalidasi sintaks URL, tetapi tidak mengecek apakah server target tersebut benar-benar aktif/ada. | Integrasi client HTTP Go (`http.Client` dengan timeout pendek) untuk melakukan ping HTTP `HEAD` ke URL target sebelum membolehkan shorten. |

## Known weaknesses worth revisiting

| Weakness | Risk if unaddressed | Candidate trigger to fix |
|---|---|---|
| **Beban Tulis Klik Tinggi (Write Hotspot)** | Setiap klik redirect memicu kueri tulis (`UPDATE click_count`). Pada traffic tinggi, ini memicu lock database PostgreSQL dan memperlambat response redirect. | Jika layanan dipublikasikan ke publik atau saat kita melatih integrasi Redis di Project 4. |
| **Tidak ada Rate Limiting** | Endpoint `POST /shorten` rentan diserang spammer yang bisa membanjiri database dengan kueri kustom alias secara brutal (*brute force*). | Sebelum meng-host layanan ini di luar jaringan localhost. |

## Ideas considered and explicitly rejected

| Idea | Why rejected |
|---|---|
| **Menyimpan short URL di file memori lokal** | *Rejected* karena kami ingin belajar menggunakan koneksi database PostgreSQL secara riil (GORM) sejak proyek pertama, bukan sekadar menggunakan penyimpanan mock di memori RAM yang akan hilang setiap server dimulai ulang. |

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen perbaikan masa depan untuk proyek URL Shortener |
