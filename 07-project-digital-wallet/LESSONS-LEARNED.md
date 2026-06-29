# Lessons Learned: Digital Wallet API

**Written:** 2026-06-29

---

## 1. What I got wrong in the initial design, and when I noticed

- **Deadlock pada Locking Paralel:** Di awal perancangan transfer, saya sempat merencanakan mengunci lock key pengirim kemudian mengunci lock key penerima secara langsung. Saya menyadari bahwa jika User A mentransfer ke User B di saat yang sama User B mentransfer ke User A secara paralel, Utas 1 mengunci A (menunggu B) dan Utas 2 mengunci B (menunggu A). Deadlock ini akan membekukan aplikasi Go server. Saya segera memperbaiki rancangan ini dengan memaksa urutan locking terstruktur (mengunci ID terkecil terlebih dahulu).

## 2. What I'd change if I rebuilt this today

- **Idempotency Key Hash Strategy:** Saat ini idempotency key dicocokkan berdasarkan key hash SHA-256 yang disimpan di Redis. Namun, jika ada dua user yang tidak sengaja mengirim request berbeda namun menggunakan nilai idempotency key yang sama (karena generated client-side bug), tabrakan bisa terjadi. Jika membangun ulang hari ini, saya akan menyertakan prefix user_id secara eksplisit di dalam hash kunci Redis: `idempotency:user:<user_id>:<key_hash>` agar keunikan terisolasi per user akun.

## 3. The concept that most needs reinforcement

- **Lua Scripts Atomicity di Redis:** Memahami mengapa skrip Lua di Redis berjalan secara *single-threaded* dan atomic sangat penting. Ini mencegah kondisi balapan (*race condition*) di mana kunci dihapus setelah masa kadaluarsa lock habis dan kunci tersebut ternyata sudah diambil alih oleh proses lain.
- **Tipe Data Decimal/Numeric:** Penggunaan numeric(15,2) di SQL menjamin akurasi saldo. Namun di tingkat Go, GORM memetakannya kembali ke float64 yang memiliki sedikit kelemahan presisi biner jika digunakan untuk operasi pertambahan ribuan kali. Mempelajari library presisi seperti `github.com/shopspring/decimal` di Go untuk memetakan balance numeric relasional adalah langkah penguatan yang sangat bagus di masa depan.

## 4. Which earlier project should be revisited

- **Project 3 (Booking System):** Pola ad-hoc JWT middleware di Project 4 ditulis terduplikasi dari Project 3. Ini menegaskan bahwa duplikasi kode otentikasi antar modul sangat tidak efisien dan memperkuat urgensi pembangunan Auth Service terpusat di Project 7.

## 5. Estimate vs. reality

- **Estimasi Waktu:** Estimasi master roadmap untuk Digital Wallet API adalah 3-4 minggu. Namun, dengan pondasi Clean Architecture dan pemahaman concurrency dari proyek sebelumnya, pengerjaan tuntas (termasuk distributed locks, idempotency, caching, dan unit test paralel) dapat diselesaikan dalam 1 hari.
