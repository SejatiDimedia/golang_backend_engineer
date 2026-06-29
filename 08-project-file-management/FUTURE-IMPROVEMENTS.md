# Future Improvements: File Management Service

---

## Deferred by design (in scope for a future project)

| Item | Why deferred | Where it's actually addressed |
|---|---|---|
| **Shared Auth Service Retrofit** | Logika auth (register/login) dan token parser saat ini masih disalin lokal secara ad-hoc (duplikasi kode). | [Project 7 (Auth Service)](../10-project-auth-service/) |
| **Notification Integration** | Notifikasi saat file berhasil diunggah / di-share belum memicu alert asinkron. | [Project 6 (Notification Service)](../09-project-notification-service/) |

## Deferred due to scope (not a future project's job, just not done)

| Item | Why deferred | Would require |
|---|---|---|
| **Chunked Multipart Resume Upload** | File berukuran gigabytes tidak bisa diunggah utuh dalam satu request HTTP biasa karena risiko timeout. | Pembangunan endpoint inisiasi chunk, upload chunk per part (S3 Multipart Upload API), dan API finalisasi assembly. |
| **File Anti-Virus Scanning** | Risiko berkas terinfeksi malware lolos diunggah ke cloud bucket. | Integrasi hook pipeline ClamAV scanner di service layer sesaat sebelum berkas diteruskan ke MinIO. |
| **Shared Link Access Control** | Pengguna ingin membagikan berkasnya ke pengguna lain tanpa memberikan hak akses hapus/tulis. | Tabel ACL (`file_shares`) dan API validasi hak akses sebelum presigned link dihasilkan. |

## Known weaknesses worth revisiting

| Weakness | Risk if unaddressed | Candidate trigger to fix |
|---|---|---|
| **Temporary File Leak on Crash** | Jika server crash di tengah-tengah upload MinIO, record metadata `PENDING` di DB bisa tertinggal tanpa compensating write. | Ketika database relasional memiliki lebih dari 100 baris dengan status `PENDING` yang berusia lebih dari 24 jam. *(Solusi: Buat background cron worker untuk menghapus record pending yang kedaluwarsa).* |
| **MinIO Connection Timeout** | Jika server MinIO offline, startup server Go langsung mati (Fatal). | Saat deployment cloud kubernetes memerlukan health-probe mandiri tanpa mematikan total pod aplikasi. |
