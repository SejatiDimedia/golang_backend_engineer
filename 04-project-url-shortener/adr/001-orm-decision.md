# ADR-001: Pilihan ORM untuk URL Shortener (GORM vs sqlx)

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Untuk menghubungkan aplikasi Go dengan PostgreSQL, kami memiliki dua pilihan utama yang umum digunakan di ekosistem Go: GORM (sebagai ORM penuh) atau sqlx (sebagai SQL mapper yang tipis di atas pustaka bawaan `database/sql`). Kami memerlukan pustaka yang mendukung tujuan pembelajaran dasar (REST, CRUD relasional) tanpa memperumit pengaturan database pada proyek pertama.

## Decision

Kami memutuskan untuk menggunakan **GORM** sebagai pustaka interaksi database untuk Project 1 (URL Shortener).

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **GORM (Chosen)** | - Auto-migration terintegrasi.<br>- Sintaks CRUD sangat ringkas.<br>- Mempercepat pengerjaan sasis REST API. | - Menyembunyikan kueri SQL asli.<br>- Memiliki overhead performa kecil karena refleksi. |
| **sqlx** | - Performa sangat cepat.<br>- Transparan (mengharuskan penulisan SQL asli). | - Harus menulis SQL manual untuk semua CRUD dasar.<br>- Membutuhkan alat migrasi terpisah (seperti golang-migrate). |

## Reasoning

Sebagai proyek pertama setelah Phase 0, fokus utama pembelajaran adalah:
1. Memahami alur request-response menggunakan framework Gin.
2. Mempelajari pemisahan struktur folder (Clean Architecture dasar).
3. Penanganan error di tingkat handler, service, dan repository.

Penggunaan sqlx pada tahap ini akan menambah beban kognitif berupa manajemen migrasi manual dan penulisan boilerplate SQL. GORM membantu mengotomatisasi bagian database sehingga kami bisa fokus pada fondasi kode Go itu sendiri.

## Consequences

- **Positif:** Menghemat waktu pembuatan tabel (`AutoMigrate`) dan query CRUD dasar.
- **Negatif:** Menutupi logika database sesungguhnya. Ada kemungkinan kami membiasakan diri dengan keajaiban (*magic*) GORM tanpa memahami operasi SQL di bawahnya.

## Revisit conditions

Keputusan ini akan ditinjau kembali pada:
- **Project 2 (Inventory Management):** Di mana transaksi SQL yang kompleks dan query optimasi performa tinggi mulai dibutuhkan. Di sana kami akan mengevaluasi apakah sqlx lebih cocok.
