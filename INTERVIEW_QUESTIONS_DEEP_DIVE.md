# Deep Dive: 40+ Technical Interview Q&A

Dokumen ini berisi kumpulan pertanyaan wawancara teknis mendalam (*deep-dive*) untuk posisi Senior/Mid Go Backend Engineer, yang dirancang khusus berdasarkan arsitektur 8 proyek portofolio yang telah Anda selesaikan.

---

## Bagian 1: Go Internals & Concurrency

### Q-1: Apa perbedaan antara Goroutine dan Thread sistem operasi (OS)? Mengapa Goroutine sangat ringan?
- **Jawaban:**
  Goroutine dikelola oleh Go Runtime Scheduler (model M:N scheduler), bukan langsung oleh OS Kernel.
  - **Memory Footprint:** Goroutine dimulai dengan ukuran stack yang sangat kecil (sekitar 2 KB) yang dapat tumbuh dan menyusut secara dinamis. Sedangkan OS Thread memiliki stack berukuran tetap (biasanya 1 MB - 2 MB).
  - **Creation & Context Switch Cost:** Context switch pada OS thread memerlukan perpindahan ke mode kernel, menyimpan register CPU, dan memakan waktu sekitar 1-2 mikrodetik. Context switch Goroutine dilakukan di mode user space oleh Go runtime, hanya menyimpan beberapa register, dan hanya memakan waktu sekitar 100-200 nanodetik.

### Q-2: Apa perbedaan antara unbuffered channel dan buffered channel di Go? Kapan Anda harus menggunakan buffered channel?
- **Jawaban:**
  - **Unbuffered Channel:** Operasi kirim (`ch <- value`) akan memblokir (*block*) Goroutine pengirim hingga ada Goroutine lain yang menerima (`<-ch`) dari channel tersebut secara bersamaan. Ini menjamin sinkronisasi instan (*strong synchronization*).
  - **Buffered Channel:** Memiliki kapasitas penyimpanan internal. Operasi kirim hanya akan memblokir pengirim jika buffer sudah penuh. Operasi penerimaan hanya akan memblokir penerima jika buffer kosong.
  - **Kapan Digunakan:** Buffered channel sangat ideal untuk membatasi throughput (*rate limiting*), antrean tugas asinkron (seperti analytics logging daemon di Project 8), atau pola worker pool untuk menghindari bottleneck pada Goroutine pengirim.

### Q-3: Mengapa kita harus mempassing `context.Context` ke seluruh call chain I/O di Go? (Project 3, 4, 6, 7, 8)
- **Jawaban:**
  `context.Context` digunakan untuk mengirimkan sinyal pembatalan (*cancellation*), batas waktu (*timeout*), dan metadata scope request ke seluruh Goroutine yang terlibat dalam proses I/O.
  - Jika klien menutup koneksi HTTP sebelum proses selesai, context akan dibatalkan (`ctx.Done()`).
  - GORM / Driver Database dan Redis yang peka terhadap context akan mendeteksi pembatalan ini dan langsung menghentikan query yang sedang berjalan di DB Server, mencegah penumpukan query yatim (*orphaned queries*) dan menghemat koneksi database.

### Q-4: Apa itu "Escape Analysis" di Go dan bagaimana pengaruhnya terhadap performa alokasi memori?
- **Jawaban:**
  Escape Analysis adalah algoritma kompilator Go untuk menentukan apakah suatu variabel dialokasikan di **Stack** atau di **Heap** saat kompilasi.
  - **Stack Allocation:** Sangat cepat karena otomatis dibersihkan saat fungsi selesai dieksekusi (*zero GC overhead*).
  - **Heap Allocation:** Variabel "lolos" (*escape*) ke heap jika nilainya direferensikan ke luar fungsi (misal mengembalikan pointer lokal). Variabel di heap harus dibersihkan oleh Garbage Collector (GC) yang dapat memicu latensi aplikasi.
  - **Optimalisasi:** Menghindari return pointer untuk struct kecil jika tidak diperlukan, guna meminimalisasi alokasi heap dan meringankan beban GC.

---

## Bagian 2: Database & Transactions (GORM/PostgreSQL/SQLite)

### Q-5: Apa itu tingkat isolasi transaksi database (Transaction Isolation Levels) dan mengapa row-locking `SELECT ... FOR UPDATE` penting? (Project 7 & 8)
- **Jawaban:**
  Isolasi transaksi menentukan seberapa aman data dari interupsi transaksi lain yang berjalan konkuren (menghindari Dirty Read, Non-repeatable Read, dan Phantom Read).
  - Saat melakukan pemutaran token (RTR) di Project 7, jika kita menggunakan SELECT biasa, dua request konkuren dari token yang sama bisa membaca status `is_revoked = false` secara bersamaan (race condition).
  - Dengan menambahkan **`FOR UPDATE`** (Row-Level Lock), PostgreSQL akan menahan dan memblokir transaksi kedua untuk membaca baris tersebut hingga transaksi pertama selesai melakukan UPDATE dan COMMIT. Ini menjamin perlindungan replay attack 100% aman.

### Q-6: Mengapa SQLite in-memory direkomendasikan untuk pengujian unit (unit testing) relasional dibandingkan mock manual (sqlmock)? (Project 7 & 8)
- **Jawaban:**
  - **Mock Manual (sqlmock):** Memaksa kita menulis ekspektasi query string secara kaku. Jika struktur query SQL berubah sedikit (seperti spasi atau kolom baru), test akan gagal walau logikanya benar.
  - **SQLite In-Memory:** Adalah database relasional nyata yang berjalan di memori RAM selama testing. Ini memungkinkan kita memvalidasi query JOIN SQL yang kompleks, foreign key constraints, index unique behavior, dan relasi many-to-many secara fungsional 100% akurat tanpa overhead jaringan fisik.

### Q-7: Bagaimana Anda memodelkan relasi dinamis RBAC (Role-Based Access Control) secara native di database relasional? (Project 7)
- **Jawaban:**
  Saya memodelkan RBAC menggunakan 5 tabel relasional dengan skema many-to-many:
  1. `users`: Menyimpan data profil user.
  2. `roles`: Menyimpan nama role (seperti `admin`, `customer`).
  3. `permissions`: Menyimpan hak akses granular (seperti `wallet:write`, `prompt:compile`).
  4. `user_roles`: Join table yang menghubungkan `user_id` ke `role_id`.
  5. `role_permissions`: Join table yang menghubungkan `role_id` ke `permission_id`.
  Dengan skema ini, kita cukup menjalankan satu kueri SQL JOIN untuk menarik seluruh permissions yang dimiliki user secara efisien saat login.

---

## Bagian 3: System Design & Distributed Systems

### Q-8: Jelaskan bagaimana mekanisme Idempotency Key bekerja untuk mencegah transaksi ganda pada API transfer saldo? (Project 4)
- **Jawaban:**
  Idempotency Key menjamin request ganda yang tidak sengaja terkirim oleh client (misal karena timeout jaringan) hanya dieksekusi satu kali oleh server:
  1. Klien mengirim request HTTP disertai header `X-Idempotency-Key` (UUID unik).
  2. Middleware memeriksa ke Redis apakah key tersebut sudah pernah terdaftar.
  3. Jika **belum ada**: Jalankan transaksi transfer saldo di database. Setelah sukses, simpan salinan body response sukses ke Redis dengan key idempotency tersebut dan TTL 1 jam. Kembalikan response ke client.
  4. Jika **sudah ada**: Langsung kembalikan salinan response yang disimpan di Redis ke client tanpa menyentuh database atau memproses transaksi transfer saldo ulang.

### Q-9: Apa keuntungan menggunakan Redis Cache-Aside Pattern untuk otentikasi API Key? Bagaimana cara kerja invalidasi cache-nya? (Project 8)
- **Jawaban:**
  - **Keuntungan:** Membaca data API Key langsung dari database PostgreSQL pada setiap request compiler prompt akan memicu overhead I/O disk. Cache-aside memindahkan lookup ke Redis in-memory dengan latensi $<2\text{ms}$.
  - **Alur Kerja:**
    1. Klien mengirim API Key, server menghitung SHA-256 hashnya.
    2. Cari hash tersebut di Redis (`apikey:<hash>`). Jika ada (*Cache Hit*), loloskan.
    3. Jika tidak ada (*Cache Miss*), cari di Postgres. Jika ditemukan, simpan ke Redis cache dengan TTL 1 jam (*Write-through*).
  - **Invalidasi Cache:** Jika admin melakukan pencabutan (*revoke*) API Key di database, server wajib langsung mengeksekusi `RDB.Del` pada key Redis tersebut agar hak akses client dicabut seketika secara real-time.

### Q-10: Bagaimana Anda merancang algoritma Exponential Backoff untuk pengiriman ulang notifikasi yang gagal? (Project 6)
- **Jawaban:**
  Exponential Backoff menunda percobaan pengiriman berikutnya secara eksponensial untuk memberi waktu server target pulih dari down, sekaligus menghindari serangan DDoS tidak sengaja (*thundering herd problem*):
  - Rumus penundaan: \(\text{delay} = 2^{\text{attempt}} \times 2\) detik.
  - Upaya 1: Tunda 4 detik.
  - Upaya 2: Tunda 8 detik.
  - Upaya 3: Tunda 16 detik.
  - Setelah batas maksimal percobaan (misal 5 kali) terlampaui, pesan dipindahkan ke **Dead-Letter Queue (DLQ)** di database agar dapat ditinjau manual oleh admin.

---

## Bagian 4: Cryptography & Security

### Q-11: Apa perbedaan antara enkripsi simetris (HMAC-SHA256) dan asimetris (RSA RS256) dalam penandatanganan token JWT?
- **Jawaban:**
  - **Simetris (HS256):** Menggunakan satu kunci rahasia yang sama untuk proses *signing* (membuat token) dan *verifying* (validasi token). Semua pihak yang ingin memvalidasi token harus tahu kunci rahasia tersebut.
  - **Asimetris (RS256):** Menggunakan sepasang kunci (Private/Public Key). Auth Service menandatangani JWT menggunakan **Private Key** (yang harus disimpan super aman). Downstream services memverifikasi validitas JWT secara offline menggunakan **Public Key** yang disebarkan bebas.

### Q-12: Mengapa kita tidak boleh menyimpan API Key mentah di database? Bagaimana cara mengamankannya? (Project 8)
- **Jawaban:**
  Jika database PostgreSQL diretas atau bocor, seluruh API Key mentah yang tersimpan dapat disalahgunakan oleh peretas untuk mengakses resource client.
  - **Pengamanan:** API Key yang dibuat berformat `prompt_live_<random_bytes>`. Kita hanya menyimpan nilai **SHA-256 hash** dari key tersebut di database.
  - Saat request masuk, kita melakukan hashing pada API Key kiriman client, lalu mencocokkan hash tersebut dengan data di database. Karena SHA-256 adalah enkripsi satu arah, peretas tidak bisa merekonstruksi API Key asli dari data hash database yang bocor.

### Q-13: Apa itu Replay Attack pada sesi otentikasi JWT dan bagaimana cara mencegahnya menggunakan Refresh Token Rotation? (Project 7)
- **Jawaban:**
  Replay Attack terjadi ketika peretas berhasil mencuri token refresh lama dan mencoba mengirimkannya berulang kali ke server untuk mendapatkan access token baru secara ilegal.
  - **Pencegahan:** Dengan **Refresh Token Rotation (RTR)**, setiap token refresh lama yang digunakan untuk mendapatkan token baru akan otomatis ditandai statusnya menjadi `is_revoked = true`.
  - Jika token berstatus `is_revoked` dikirimkan kembali ke server, sistem mendeteksinya sebagai penyerangan, memblokir request, dan langsung menonaktifkan seluruh silsilah refresh token aktif milik user tersebut di database untuk force logout massal.

---

## Bagian 5: Platform Engineering & Clean Architecture

### Q-14: Mengapa implementasi worker queue asinkron menggunakan Redis disarankan menggunakan Lua Script? (Project 6)
- **Jawaban:**
  Redis memproses perintah satu per satu (*single-threaded*). **Lua Script** dieksekusi secara atomik di dalam Redis engine.
  - Saat poller scheduler mencari tugas delayed queue yang jatuh tempo (`ZRANGEBYSCORE`) lalu memindahkannya ke antrean aktif (`LPUSH`), jika proses ini dieksekusi lewat perintah terpisah dari Go, node server lain bisa menyela di tengah-tengah proses tersebut (*race condition*).
  - Dengan Lua Script, seluruh alur baca-pindah tersebut dikunci dan dieksekusi sebagai satu operasi atomik tunggal, mencegah duplikasi tugas pada arsitektur server terdistribusi.

### Q-15: Apa itu "Clean Architecture" dan apa keuntungan membagi folder program menjadi handler, service, dan repository? (Project 2-8)
- **Jawaban:**
  Clean Architecture adalah pemisahan kode berdasarkan tanggung jawab (*separation of concerns*) untuk menjaga domain bisnis tetap independen dari infrastruktur luar:
  - **Handler:** Mengurus HTTP framework, parsing JSON request, dan menulis HTTP response.
  - **Service:** Pusat logika bisnis platform (tidak peduli database apa yang digunakan).
  - **Repository:** Mengurus akses query data langsung ke database relasional/NoSQL.
  - **Keuntungan:** Kode menjadi sangat mudah di-unit test (kita bisa menguji Service layer hanya dengan meng-inject mock repository), fleksibel (bisa mengganti database PostgreSQL ke MongoDB hanya dengan mengubah repository layer tanpa menyentuh core business logic di service).

### Q-16: Bagaimana cara menangani deadlock ketika dua request konkuren melakukan transaksi transfer balance saldo berlawanan arah secara bersamaan? (Project 4)
- **Jawaban:**
  Deadlock terjadi ketika Transaksi 1 mengunci Akun A dan menunggu kunci Akun B dibebaskan, sementara Transaksi 2 mengunci Akun B dan menunggu kunci Akun A dibebaskan. Keduanya terkunci selamanya.
  - **Solusi:** Selalu akuisisi lock dengan aturan urutan yang konsisten (*consistent locking order*). Contohnya, kita selalu membandingkan ID Akun pengirim dan penerima, lalu selalu mengunci ID yang lebih kecil terlebih dahulu sebelum mengunci ID yang lebih besar, tidak peduli siapa pengirim atau penerimanya. Ini menjamin alur penguncian searah dan mencegah lingkaran deadlock.

### Q-17: Mengapa kita harus menerapkan Graceful Shutdown pada server produksi Go? Bagaimana caranya? (Project 7 & 8)
- **Jawaban:**
  Shutdown yang kasar (*hard kill*) akan langsung memutus koneksi HTTP klien yang sedang berjalan di tengah proses penulisan database, memicu korupsi data atau transaksi gantung.
  - **Cara Kerja:**
    1. Tangkap sinyal interupsi OS (`SIGINT`, `SIGTERM`) menggunakan channel `os.Signal`.
    2. Begitu sinyal diterima, panggil `srv.Shutdown(ctx)` dengan batas timeout (misal 5 detik).
    3. Server akan berhenti menerima request baru, menyelesaikan seluruh HTTP request yang sedang berjalan, membiarkan background worker database menyelesaikan antrean sisa, memutus koneksi Redis/DB secara bersih, lalu exit dengan aman.

### Q-18: Apa itu "Thundering Herd Problem" di sistem caching terdistribusi dan bagaimana cara Anda memitigasinya?
- **Jawaban:**
  Terjadi ketika data cache yang sangat sering diakses (*hot key*) kedaluwarsa secara tiba-tiba. Ratusan request konkuren yang masuk bersamaan akan mengalami *cache miss* dan semuanya langsung menyerang database PostgreSQL secara serentak, membuat database crash karena beban overload.
  - **Mitigasi:** Menggunakan distributed locking (seperti Redis Lock) agar hanya satu request pertama yang diizinkan melakukan query ke database dan menulis ulang cache, sementara request lainnya menunggu atau membaca cache stale sejenak.

### Q-19: Bagaimana estimasi token length dihitung di Go Compiler Engine Anda secara efisien tanpa memanggil API pihak ketiga? (Project 8)
- **Jawaban:**
  Memanggil API tokenizer eksternal (seperti tiktoken) di setiap request kompilasi prompt akan merusak target latensi. 
  - **Pola Estimasi:** Saya memecah teks terkompilasi menjadi kata-kata menggunakan `strings.Fields(compiled)`.
  - Berdasarkan standar industri LLM, 1 kata rata-rata mewakili sekitar 1.33 token. Saya mengalikan jumlah kata dengan `1.33` untuk mendapatkan estimasi token instan ($O(N)$ komplekstitas CPU lokal) dengan tingkat akurasi yang memadai untuk logging analitik dasar.

---

## Bagian 6: REST API Design & HTTP Protocols

### Q-20: Kapan Anda menggunakan status HTTP `201 Created` vs `202 Accepted`? Berikan contohnya dari proyek Anda! (Project 6 & 8)
- **Jawaban:**
  - **`201 Created`:** Digunakan ketika resource sukses dibuat secara sinkron dan instan sebelum respons HTTP dikembalikan. Contoh: membuat database Workspace (`POST /workspaces`) di Project 8.
  - **`202 Accepted`:** Digunakan untuk operasi asinkron di mana server menerima request dan memasukkannya ke antrean, tetapi pemrosesan sesungguhnya belum selesai. Contoh: mengirim notifikasi (`POST /notifications`) di Project 6, di mana server hanya mengembalikan status sukses antrean, lalu background worker pool yang memproses pengiriman aslinya nanti.

### Q-21: Bagaimana Anda menangani API Rate Limiting secara terdistribusi di platform microservices?
- **Jawaban:**
  Saya menggunakan database **Redis** dengan algoritma **Token Bucket** atau **Leaky Bucket**.
  - Setiap client ID memetakan ke satu key di Redis yang mencatat sisa token limit dan timestamp pengisian ulang (*refresh rate*).
  - Menggunakan command Redis secara atomik (atau script Lua) untuk mengurangi jumlah token di setiap request. Jika jumlah token mencapai 0, server langsung mengembalikan status `429 Too Many Requests` secara instan tanpa membebani database PostgreSQL.

### Q-22: Apa perbedaan antara HTTP POST, PUT, dan PATCH? Kapan Anda memilih salah satunya?
- **Jawaban:**
  - **`POST`:** Digunakan untuk membuat resource baru (*non-idempotent*). Setiap request POST baru akan membuat baris baru di database (contoh: registrasi user).
  - **`PUT`:** Digunakan untuk mengganti resource secara keseluruhan secara *idempotent*. Mengirim data PUT yang sama berulang kali tidak mengubah state database setelah pemanggilan pertama (contoh: mengaktifkan prompt version).
  - **`PATCH`:** Digunakan untuk melakukan modifikasi parsial (hanya beberapa kolom tertentu) pada suatu resource.

---

## Bagian 7: Advanced Data Structures & Memory Management

### Q-23: Bagaimana Go Map bekerja di bawah kap (*under the hood*)? Apakah map aman dari race condition?
- **Jawaban:**
  Go map diimplementasikan sebagai **Hash Table** yang tersusun atas deretan kotak memori bernama **Buckets** (masing-masing bucket menampung hingga 8 key-value pair). Pointer hash digunakan untuk menentukan lokasi bucket tujuan.
  - **Thread-Safety:** Go map **tidak aman** dari akses konkuren (*not thread-safe*). Jika satu Goroutine menulis ke map sementara Goroutine lain membaca secara konkuren, runtime Go akan langsung panik (*fatal error: concurrent map writes*).
  - **Mitigasi:** Menggunakan pelindung mutex (`sync.RWMutex` seperti di Project 2) atau beralih ke `sync.Map` untuk kasus read-heavy konkuren.

### Q-24: Mengapa `sync.Pool` sangat berguna di Go? Berikan contoh kasus penggunaannya!
- **Jawaban:**
  `sync.Pool` adalah kumpulan objek temporer yang dapat digunakan kembali secara terpisah untuk meminimalisasi siklus alokasi memori baru di heap, sehingga mengurangi beban kerja Garbage Collector (GC).
  - **Contoh Kasus:** Saat mengompilasi teks prompt berulang kali atau melakukan marshaling JSON, kita membutuhkan buffer byte (`bytes.Buffer`). Dibanding membuat buffer baru di setiap request (`var buf bytes.Buffer` yang memicu heap allocation), kita meminjam buffer dari `sync.Pool`, menggunakannya, mereset isinya, lalu mengembalikannya kembali ke pool.

### Q-25: Bagaimana cara mendeteksi kebocoran memori (Memory Leak) di aplikasi Go yang berjalan sebagai background daemon worker? (Project 6 & 8)
- **Jawaban:**
  Memory Leak di Go biasanya terjadi karena Goroutine yang terblokir selamanya (tidak pernah exit), menyisakan referensi variabel di heap secara permanen.
  - **Deteksi:** Menggunakan tool bawaan **pprof** (Go Profiler) untuk menangkap heapsnapshot saat server berjalan.
  - **Analisis:** Mengaktifkan route `/debug/pprof/heap` dan memantau pertumbuhan jumlah Goroutine aktif serta ukuran heap memori seiring waktu menggunakan visual grafis pprof. Jika jumlah Goroutine terus bertambah tanpa pernah turun, dipastikan ada Goroutine leak.

---

## Bagian 8: Database & SQL Advanced Query Optimization

### Q-26: Apa perbedaan antara indeks B-Tree dan Hash Index di PostgreSQL? Kapan Anda menggunakannya?
- **Jawaban:**
  - **B-Tree Index (Default):** Mengorganisasikan data dalam struktur pohon seimbang. Sangat optimal untuk pencarian perbandingan (`=`), range query (`<`, `>`, `BETWEEN`), dan pengurutan (`ORDER BY`).
  - **Hash Index:** Hanya mendukung operator perbandingan sama dengan (`=`). Sangat cepat untuk pencarian nilai tunggal yang eksak (seperti query `key_hash` API Key di Project 8), namun tidak bisa digunakan untuk range query atau sort order.

### Q-27: Apa perbedaan antara Row Lock dan Table Lock? Kapan Deadlock di level tabel terjadi?
- **Jawaban:**
  - **Row Lock (e.g. `SELECT FOR UPDATE`):** Hanya mengunci baris data spesifik yang di-query. Transaksi lain masih bebas membaca/menulis baris lain pada tabel yang sama.
  - **Table Lock (e.g. `LOCK TABLE`):** Mengunci seluruh tabel. Transaksi lain diblokir total untuk melakukan perubahan pada tabel tersebut.
  - **Deadlock Level Tabel:** Terjadi jika Transaksi 1 menahan Row Lock di Tabel A lalu mencoba mengunci Tabel B, sementara Transaksi 2 menahan Row Lock di Tabel B lalu mencoba mengunci Tabel A secara konkuren.

### Q-28: Bagaimana cara mengoptimalkan query JOIN lambat yang melibatkan jutaan data relasional RBAC di PostgreSQL?
- **Jawaban:**
  1. **Indeks Kolom Relasi:** Pastikan kolom Foreign Key (`user_id`, `role_id`, `permission_id`) dilindungi oleh indeks B-Tree.
  2. **Query Plan Analysis (`EXPLAIN ANALYZE`):** Jalankan query di console dengan prefiks `EXPLAIN ANALYZE` untuk mendeteksi apakah Postgres melakukan *Sequential Scan* (lambat) atau *Index Scan* (cepat).
  3. **Denormalisasi Caching:** Jika beban join sangat tinggi, denormalisasikan permissions user langsung ke dalam payload JWT (seperti Project 7) atau simpan di Redis cache agar tidak perlu men-JOIN 5 tabel relasional di setiap HTTP request.

---

## Bagian 9: Security, Encryption & Token Standards

### Q-29: Mengapa bcrypt dirancang lambat untuk hashing password? Mengapa kita tidak menggunakan SHA-256 saja? (Project 7)
- **Jawaban:**
  SHA-256 dirancang sangat cepat oleh CPU (dapat memproses jutaan hash per detik). Jika database bocor, peretas bisa melakukan brute-force attack menebak password user dengan sangat mudah menggunakan GPU modern.
  - **bcrypt** memiliki parameter **Work Factor (Cost)** yang sengaja memperlambat komputasi hash (membutuhkan sekitar 100-300ms untuk satu kali hash). Kelambatan ini membuat serangan brute-force secara massal menjadi terlalu mahal secara komputasi dan tidak layak dilakukan oleh peretas, sementara overhead 100-300ms bagi user asli saat login masih sangat wajar.

### Q-30: Bagaimana cara mengimplementasikan JWT Blacklist terdistribusi menggunakan Redis secara efisien?
- **Jawaban:**
  Saat user logout sebelum masa aktif Access Token habis, kita harus membatalkan token tersebut secara global:
  1. Ambil JWT ID (`jti`) unik dan sisa waktu aktif token (*Expiration Time*).
  2. Simpan `jti` tersebut ke Redis dengan key `blacklist:<jti>` dan atur TTL (Time to Live) Redis persis sama dengan sisa umur token tersebut.
  3. Setiap request yang masuk disaring oleh middleware: jika `jti` token terdaftar di Redis blacklist, langsung tolak dengan status `401 Unauthorized`. Setelah sisa umur token habis, Redis otomatis menghapus key tersebut untuk menghemat memori.

### Q-31: Bagaimana cara menangani konsistensi data otorisasi (Permissions) di cache Redis jika admin mengubah permission user di database? (Project 8)
- **Jawaban:**
  Ini adalah tantangan konsistensi data *cache invalidation*. 
  - Saat admin mengubah role/permission user di DB, server harus memicu event invalidasi dengan memanggil `RDB.Del` pada key cache terkait (seperti cache API Key `apikey:<hash>` atau cache introspect token).
  - Request berikutnya otomatis memicu *cache miss*, membaca data terbaru dari database PostgreSQL, dan menulis ulang data segar tersebut ke cache Redis.

---

## Bagian 10: System Design Patterns & Infrastructure

### Q-32: Apa itu Circuit Breaker Pattern dalam arsitektur microservices? Kapan Anda menggunakannya?
- **Jawaban:**
  Circuit Breaker mencegah kegagalan satu service merembet ke service lain (*cascading failure*).
  - Jika **Notification Service** (Project 6) mengalami down, **Booking Service** (Project 3) yang memanggilnya akan mengalami penumpukan request antrean tunggu (*timeout queue*) hingga kehabisan RAM.
  - **Cara Kerja:** Circuit Breaker memantau kegagalan. Jika rasio error melebihi batas (misal 50% request gagal), circuit berubah menjadi **OPEN** (langsung mengembalikan error lokal tanpa memanggil Notification Service fisik). Setelah jeda waktu tertentu, circuit menjadi **HALF-OPEN** untuk menguji apakah Notification Service sudah pulih sebelum kembali ke status **CLOSED** (normal).

### Q-33: Bagaimana Anda mengamankan File Upload Service dari serangan malicious file upload (seperti upload virus)? (Project 5)
- **Jawaban:**
  1. **Validation MIME Type:** Jangan percaya extension file (`.jpg`). Baca 512 byte pertama file menggunakan fungsi `http.DetectContentType(buffer)` untuk memvalidasi tipe berkas sesungguhnya secara biner.
  2. **Restrict File Size:** Batasi ukuran maksimum upload (misal 5 MB) di level web server Nginx / Gin Middleware.
  3. **Malware Scanning:** Lewatkan file ke server scanning Antivirus (seperti ClamAV) sebelum menyimpannya ke MinIO Object Storage.
  4. **Randomized Filenames:** Selalu ganti nama file asli dengan UUID acak di cloud storage untuk mencegah peretas menebak lokasi file dan memicu eksekusi kode jarak jauh (*Remote Code Execution*).

### Q-34: Bagaimana Anda mendesain server Go agar bersifat "Stateless" untuk mempermudah Horizontal Scaling?
- **Jawaban:**
  Server Go stateless artinya server tidak menyimpan data sesi (*session state*) di memori RAM lokalnya.
  - **Langkah Desain:**
    1. Seluruh data session dipindahkan ke database PostgreSQL terpusat atau Redis cluster terdistribusi.
    2. Autentikasi menggunakan token JWT RS256 sehingga validasi data user bisa dilakukan mandiri tanpa menyimpan session state lokal di memori pod.
    3. File upload disimpan ke Object Storage terpusat (seperti MinIO/S3), bukan di local disk server.
  - Dengan desain stateless, kita bebas menyalakan 100 replika server Go di belakang load balancer tanpa perlu khawatir klien terlempar ke replika mana pun.
