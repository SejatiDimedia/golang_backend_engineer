# Lessons Learned: URL Shortener Service

**Written:** 2026-06-29

---

## 1. What I got wrong in the initial design, and when I noticed

- **Perbandingan Tipe Data customAlias:** Di dalam file [service/url.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/internal/service/url.go), saya awalnya membandingkan `customAlias string` dengan `nil` (`customAlias != nil`). Di Go, tipe data string bawaan tidak bisa bernilai `nil` karena nil hanya berlaku untuk pointer, interface, channel, map, slice, dan function. Saya menyadari hal ini saat kompilasi testing dijalankan di terminal dan compiler menolak *mismatched types string and untyped nil*.
- **Import di Uji Unit:** Di file [service/url_test.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/internal/service/url_test.go), saya memanggil pustaka `errors.New` namun lupa menyertakan `"errors"` di blok import utama. Ini menyadarkan saya pentingnya memanfaatkan utilitas formatter/importer otomatis (seperti `goimports` atau IDE autocomplete) saat coding mandiri.

## 2. What I'd change if I rebuilt this today

- **Pengembalian Custom Error:** Daripada hanya membandingkan string error secara manual, saya akan mendesain custom struct error yang membawa HTTP status code secara implisit di level service. Hal ini mempermudah handler layer memetakan error ke respon HTTP tanpa perlu banyak blok `errors.Is` bertumpuk.
- **Penggunaan SQL Mentah:** Meskipun GORM mempermudah `AutoMigrate`, saya merasa kueri SQL relasional di balik layar menjadi terlalu gelap. Jika mendesain ulang hari ini, saya ingin menantang diri menggunakan `sqlx` guna memperdalam pemahaman sintaks SQL DDL dan DML secara langsung.

## 3. The concept that most needs reinforcement

- **Dependency Injection (DI) secara Manual:** Menyatukan Handler, Service, dan Repository di [main.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/04-project-url-shortener/cmd/server/main.go) masih dilakukan secara manual. Walaupun aman untuk proyek skala ini, ini bisa menjadi sangat besar dan sulit dikelola saat jumlah service bertambah. Saya perlu membaca materi mengenai dependency injection framework (seperti Google Wire) atau pola pengaturan DI di Go.
- **Context Cancellation:** Walaupun saya telah meneruskan `c.Request.Context()` dari Gin ke kueri GORM, saya belum mempraktikkan pembatalan kueri secara riil (misal menggunakan simulasi query lambat/timeout). Saya perlu membuat drill khusus di folder `03-exercises/` untuk memahami siklus hidup Context.

## 4. Which earlier project should be revisited

Ini adalah proyek pertama dalam roadmap, sehingga tidak ada proyek sebelumnya yang dapat dikunjungi kembali. Namun, pelajaran di proyek ini (khususnya pembagian layer arsitektur) akan langsung dibawa dan diperluas di **Project 2 (Inventory Management API)**.

## 5. Estimate vs. reality

- **Estimasi Waktu:** Roadmap repo-level mengestimasikan 2–3 minggu untuk Project 1. Secara realitas, karena dasar pemrograman umum sudah dikuasai, pembuatan sasis kode dan test suite diselesaikan jauh lebih cepat (~1-2 hari).
- **Unknowns & Risks:** Concurrency click counter berhasil diantisipasi sejak tahap roadmap dan diimplementasikan secara aman menggunakan ekspresi SQL atomic di repository (`gorm.Expr("click_count + 1")`), menghindari race condition data klik saat diakses simultan.

## 6. What surprised me

- **Kemudahan testing di Gin:** Penggunaan `httptest.NewRecorder()` bersama Gin engine sangat mudah dan bersih. Hal ini memungkinkan simulasi request REST API lengkap tanpa perlu meluncurkan port TCP HTTP server nyata di sistem.
- **Docker Compose Port Conflict:** Saya terkejut mendapati docker daemon tidak berjalan secara default di sistem host yang memicu error saat menjalankan `docker-compose up`. Ini menjadi pengingat penting bahwa integrasi containerization sangat bergantung pada status software daemon eksternal.
