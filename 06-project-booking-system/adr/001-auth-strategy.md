# ADR-001: Strategi Otentikasi JWT Ad-Hoc

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Layanan pemesanan (*booking*) coworking space memerlukan sistem otentikasi untuk membatasi hak akses pelanggan dan admin. Kami perlu melindungi endpoint pemesanan dan pembatalan, serta mengidentifikasi siapa pemilik booking terkait.

Di tingkat repository roadmap, Project 7 adalah pembuatan sebuah layanan otentikasi terpusat (Auth Service). Namun, untuk Project 3 ini, kami perlu menentukan apakah akan mengintegrasikan otentikasi secara langsung (*ad-hoc*) di dalam codebase proyek ini atau langsung mematangkan shared service.

## Decision

Kami memutuskan untuk mengimplementasikan **otentikasi JWT secara mandiri (ad-hoc)** langsung di dalam codebase Project 3, menggunakan pustaka `golang-jwt/jwt` dan hashing password dengan `bcrypt`.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Otentikasi JWT Ad-Hoc (Chosen)** | - Sangat mandiri (self-contained), tidak bergantung pada service luar.<br>- Menghemat waktu setup di awal.<br>- Memberikan pelajaran nyata mengenai duplikasi middleware otentikasi sebelum disatukan di Project 7. | - Kode otentikasi terduplikasi jika proyek berikutnya memerlukan auth serupa.<br>- Harus mengelola tabel `users` secara lokal. |
| **B. Shared Auth Service (Premature)** | - Bersih dari awal, tidak ada duplikasi kode otentikasi. | - Melanggar alur roadmap pembelajaran.<br>- Kompleksitas setup network/routing bertambah sebelum pondasi backend kami matang. |

## Reasoning

Mengikuti panduan master roadmap (`01-roadmap.md` §5), penulisan otentikasi JWT secara ad-hoc di Project 3 dan Project 4 adalah langkah pembelajaran yang disengaja. Rasa sakit akibat menduplikasi logika otentikasi, verifikasi token, dan manajemen user di dua layanan terpisah adalah bahan bakar yang memotivasi mengapa kita merancang *Authentication Service* khusus di Project 7.

Oleh karena itu, membuat shared auth service pada tahap ini dinilai terlalu dini (*premature optimization*).

## Consequences

- **Positif:** Kecepatan inisiasi sasis backend tinggi, database mandiri tanpa cross-service calls.
- **Negatif:** Tabel `users` dipasang di database lokal `booking_db`, dan logika JWT parsing ditulis langsung sebagai Gin Middleware lokal.

## Revisit conditions

Keputusan ini akan ditinjau kembali secara resmi pada **Project 7 (Authentication Service)**. Saat layanan otentikasi terpusat selesai dibangun, Project 3 ini akan direnovasi (*retrofit*) untuk membuang auth lokal dan mengonsumsi shared service.
