# Lessons Learned: Booking Management System

**Written:** 2026-06-29

---

## 1. What I got wrong in the initial design, and when I noticed

- **Binding Datetime di Gin:** Awalnya saya tidak menyadari bahwa tag parameter binding Gin `time_format` memerlukan format layout yang tepat untuk RFC3339 dengan timezone offset. Jika format layout yang dipaksakan salah, Gin akan gagal melakukan de-serialization JSON payload tanggal dan mengembalikan error HTTP 400. Saya menyadari hal ini saat menulis handler booking dan segera memperbaikinya dengan tag format `time_format:"2006-01-02T15:04:05Z07:00"` pada struct request.

## 2. What I'd change if I rebuilt this today

- **Token Expiry Management:** Token JWT ad-hoc yang dihasilkan saat ini tidak dapat dibatalkan secara paksa sebelum kedaluwarsa (*no token revocation/blacklist*). Jika password user diganti atau user melakukan logout, token JWT lama tetap valid hingga masa aktifnya habis. Jika membangun ulang hari ini, saya ingin mengevaluasi integrasi penyimpanan blacklist token berbasis Redis untuk memitigasi risiko keamanan ini.

## 3. The concept that most needs reinforcement

- **Penanganan Zona Waktu (Go time.Time):** Meskipun kami secara ketat memanggil `.UTC()` di service layer, perilaku library database GORM yang memetakan struct `time.Time` ke tipe `timestamp without time zone` di PostgreSQL dapat memicu konversi otomatis ke zona lokal database jika driver koneksi tidak dikonfigurasikan dengan parameter `TimeZone=UTC`. Saya perlu menegaskan parameter koneksi database yang aman untuk menghindari drift waktu.

## 4. Which earlier project should be revisited

- **Project 2 (Inventory Management):** Pola isolasi transaksi `TransactionManager` berbasis context yang kami bawa dari Project 2 terbukti sangat adaptif dan mempermudah penulisan unit test untuk logic bisnis transaksional kompleks (seperti overlap check) tanpa harus melakukan setup mock database yang rumit.

## 5. Estimate vs. reality

- **Estimasi Waktu:** Estimasi pengerjaan master roadmap untuk Booking System adalah 3-4 minggu. Namun, karena kerangka kerja dependency injection dan Clean Architecture sudah matang, penulisan seluruh backend logic, JWT helper, role middleware, dan pengujian unit diselesaikan dalam waktu kurang dari 1 hari.
