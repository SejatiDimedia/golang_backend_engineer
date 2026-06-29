# Future Improvements: Booking Management System

---

## Deferred by design (in scope for a future project)

| Item | Why deferred | Where it's actually addressed |
|---|---|---|
| **Shared Auth Service Retrofit** | Logika registrasi, login, dan validasi middleware JWT saat ini ditulis lokal secara ad-hoc (duplikasi kode). | [Project 7 (Auth Service)](../10-project-auth-service/) |
| **Notification Integration** | Konfirmasi booking saat ini hanya dicatat di standard stdout log (stub). | [Project 6 (Notification Service)](../09-project-notification-service/) |
| **Distributed Concurrency Lock (Redis)** | Penggunaan pessimistic lock `FOR UPDATE` di SQL membebani database jika traffic konkuren sangat tinggi. | [Project 4 (Digital Wallet)](../07-project-digital-wallet/) |

## Deferred due to scope (not a future project's job, just not done)

| Item | Why deferred | Would require |
|---|---|---|
| **Booking Reschedule** | Fitur mengubah jadwal booking yang sudah dipesan tanpa membatalkannya terlebih dahulu. | Penulisan service logic reschedule transaksional dengan double overlap-check (memeriksa bentrokan dengan mengabaikan ID booking aktif itu sendiri). |
| **Recurring Bookings** | Pengguna (terutama korporat) sering kali ingin memesan ruang rapat yang sama setiap hari Senin jam 09:00 selama 1 bulan penuh. | Perancangan skema relasi booking berulang (*recurring patterns*) dan loop pembuatan baris booking massal secara berkala. |
| **Capacity Management** | Ruangan rapat memiliki kapasitas maksimal (misal: 10 orang). Sistem belum memvalidasi jumlah orang saat memesan. | Penambahan kolom `capacity` di tabel `desks` dan `attendees_count` di request booking. |

## Known weaknesses worth revisiting

| Weakness | Risk if unaddressed | Candidate trigger to fix |
|---|---|---|
| **SQL Row-Locking Bottleneck** | Mengunci baris `desks` saat booking membatasi throughput paralel. Dua user memesan meja berbeda namun berada di ID yang berurutan bisa memicu lock contention tergantung indeks database. | Ketika response time API `/bookings` melebihi 200ms di bawah beban konkuren tinggi. |
| **UTC Conversion Drift** | Pengguna yang mengirimkan tanggal tanpa offset timezone default (seperti `2026-06-29 14:00:00`) akan diparsing oleh Go menggunakan zona waktu lokal server, menghasilkan drift. | Jika staff gudang/admin di lapangan melaporkan jam pemesanan tergeser otomatis. |
