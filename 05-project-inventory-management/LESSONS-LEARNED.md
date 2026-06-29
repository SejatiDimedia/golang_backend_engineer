# Lessons Learned: Inventory Management API

**Written:** 2026-06-29

---

## 1. What I got wrong in the initial design, and when I noticed

- **Inisiasi Route Path Parameter:** Di dalam file [main.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/cmd/server/main.go), saya hampir mendaftarkan route `GET /products/:id` sebelum `GET /products/export`. Di framework Gin (dan router http sejenis), jika route path parameter diletakkan di atas route statis, Gin akan salah mencocokkan request `/products/export` dan menganggap kata `"export"` sebagai `:id` produk, menghasilkan error parse ID. Saya menyadari limitasi ini saat menyusun peta route dan segera menaruh route `/products/export` di atas `/products/:id`.

## 2. What I'd change if I rebuilt this today

- **Penggunaan SQL Transactions yang Eksplisit:** Meskipun pola wrapper transaksi context manager yang kami rancang di [tx_manager.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/internal/repository/tx_manager.go) sangat bersih dan memisahkan detail DB dari service, implementasi ini tetap membawa overhead kognitif berupa dependency terselubung (context key extraction). Jika membangun ulang hari ini, saya ingin mengevaluasi penerapan Clean Architecture dengan Unit of Work (UoW) pattern untuk melihat mana yang lebih mudah diuji dan dikelola.

## 3. The concept that most needs reinforcement

- **Pessimistic vs Optimistic Locking:** Pada Stock Out, kami menggunakan pessimistic locking (`FOR UPDATE`). Walaupun aman, ini memicu lock database yang berat. Saya perlu memperdalam konsep optimistic locking (menggunakan kolom versioning/timestamp: `UPDATE products SET stock = stock - qty, version = version + 1 WHERE id = ? AND version = ?`) untuk membandingkan karakteristik performa dan skalabilitas keduanya di bawah beban kueri konkuren.
- **File Streaming di Gin:** Proses ekspor CSV menggunakan `csv.NewWriter(c.Writer)` menulis secara sinkron. Jika data sangat besar, kueri database mengambil ratusan ribu baris sekaligus ke memori dapat memicu Out Of Memory (OOM). Saya perlu mempelajari pemrosesan data chunk (cursor fetching) agar memori tetap hemat saat melakukan streaming data besar.

## 4. Which earlier project should be revisited

- **Project 1 (URL Shortener):** Tidak ada revisi mendesak untuk Project 1 saat ini. Namun, keberhasilan pola dependency injection dan isolasi error handler di Project 2 menegaskan bahwa struktur folder Clean Architecture yang kami pakai sudah sangat matang dan siap dipertahankan di proyek mendatang.

## 5. Estimate vs. reality

- **Estimasi Waktu:** Estimasi master roadmap untuk Project 2 adalah 2-3 minggu. Namun, karena kami telah menyelesaikan sasis dasar di Project 1 dan pemahaman Go-specific idioms sudah lebih baik, pembuatan database relasional 4 tabel, transaction manager, CSV streaming, dan unit test dapat diselesaikan dalam 1-2 hari.
- **Unknowns & Risks:** Kebocoran transaksi (transaction leak) berhasil dicegah berkat helper `db.Transaction` bawaan GORM yang secara internal otomatis memicu rollback jika callback mengembalikan error dan menutup transaksi dengan benar.

## 6. What surprised me

- **Kemudahan Preloading di GORM:** Sintaks `Preload` GORM untuk memuat Category dan Supplier sangat memudahkan penulisan query. Saya terkejut betapa sedikitnya boilerplate code yang diperlukan untuk menyajikan data relasional kompleks ke format JSON REST API jika dibandingkan dengan penulisan query JOIN manual di sqlx.
