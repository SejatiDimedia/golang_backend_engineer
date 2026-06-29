# Lessons Learned: File Management Service

**Written:** 2026-06-29

---

## 1. What I got wrong in the initial design, and when I noticed

- **Siklus Transaksi Terdistribusi (PostgreSQL vs MinIO):** Di awal perancangan, saya memikirkan pembungkusan upload file fisik ke dalam transaksi database SQL biasa. Namun saya menyadari bahwa GORM `tx.Rollback()` tidak bisa membatalkan objek yang terlanjur terunggah ke server MinIO eksternal. Jika database PostgreSQL gagal menyimpan data *setelah* upload MinIO sukses, berkas fisik di MinIO menjadi yatim (*orphaned file*).
- **Penyelesaian:** Mengubah arsitektur alur menjadi **Compensating Write** dengan menulis status `PENDING` di DB relasional, melempar objek ke MinIO, dan memperbaruinya menjadi `SUCCESS`. Jika melempar ke MinIO gagal, DB `PENDING` langsung di-delete.

## 2. What I'd change if I rebuilt this today

- **MIME Type Detection by Byte Header:** Saat ini, server mempercayai header `Content-Type` yang dikirim dari form HTTP client. Padahal, penyerang bisa merubah nama berkas `.exe` menjadi `.jpg` dan mengirim header palsu. Jika membangun ulang hari ini, saya akan membaca 512 byte pertama dari berkas di memori backend dan memanggil `http.DetectContentType` untuk memverifikasi tipe berkas asli berdasarkan struktur byte.

## 3. The concept that most needs reinforcement

- **Offloading Bandwidth via Presigned URLs:** Memahami betapa besarnya penghematan CPU dan memori server dengan Presigned URL. Dibandingkan dengan Direct Streaming (di mana data byte harus melewati memori server Go), Presigned URL membiarkan browser client mengunduh langsung dari AWS S3/MinIO. Ini adalah pola arsitektur mutlak untuk sistem berskala jutaan unduhan.

## 4. Which earlier project should be revisited

- **Project 4 (Digital Wallet):** Logika error handling dan struktur response API JSON di File Management ini terasa lebih modular dan terstruktur rapi. Standarisasi format parsing params dan middleware di Project 4 bisa ditingkatkan agar setara kebersihannya dengan Project 5.

## 5. Estimate vs. reality

- **Estimasi Waktu:** Estimasi master roadmap untuk File Management Service adalah 2-3 minggu. Namun, karena pemahaman fundamental adapter storage dan Gin streaming sudah matang, pengerjaan tuntas (termasuk unit test compensating rollback) dapat diselesaikan dalam waktu kurang dari 2 jam.
