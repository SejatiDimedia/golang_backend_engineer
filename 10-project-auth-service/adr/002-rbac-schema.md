# ADR-002: Pilihan Model Kebijakan Otorisasi (RBAC)

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Layanan Auth Service harus mendukung otorisasi Role-Based Access Control (RBAC) terperinci. Kami harus memutuskan apakah akan memodelkan kebijakan otorisasi secara native menggunakan tabel relasional database di PostgreSQL atau menggunakan pustaka kebijakan seperti `Casbin`.

## Decision

Kami memutuskan memodelkan RBAC secara **Native PostgreSQL Relational Database Tables** menggunakan GORM.

Kami akan membuat 5 tabel berikut:
1. `users` (id, email, password_hash, is_verified, etc.)
2. `roles` (id, name, description)
3. `permissions` (id, name, description, e.g. `wallet:read`, `booking:create`)
4. `role_permissions` (role_id, permission_id) - Tabel join relasi many-to-many
5. `user_roles` (user_id, role_id) - Tabel join relasi many-to-many

Pemeriksaan hak akses dilakukan dengan kueri JOIN SQL sederhana: memverifikasi apakah user memiliki role yang diasosiasikan dengan permission yang diminta.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Postgres Native (Chosen)** | - Sangat intuitif, mudah dirancang dan didokumentasikan di ERD.<br>- Menggunakan SQL JOIN standar yang berkinerja tinggi.<br>- Tidak menambahkan overhead runtime parsing aturan eksternal. | - Mengubah kebijakan roles/permissions secara dinamis membutuhkan penulisan baris database (INSERT/DELETE). |
| **B. Pustaka Casbin** | - Mendukung aturan otorisasi kompleks (seperti ABAC, RBAC dengan domain).<br>- Kebijakan dapat ditulis dalam berkas teks model terpisah. | - Memiliki learning curve tinggi karena menggunakan sintaks DSL (Domain Specific Language) khusus.<br>- Overhead parsing aturan dapat menurunkan performa verifikasi otorisasi. |

## Reasoning

Notification Service dan Wallet Service membutuhkan otorisasi RBAC tingkat dasar hingga menengah yang dinamis namun dapat dipahami dengan mudah. Dengan menggunakan PostgreSQL Relational Tables (Opsi A), kami memaksimalkan tujuan pembelajaran pemodelan relasi database dan ORM di Go. Struktur relasi tabel juga sangat transparan bagi tim operator sistem untuk memodifikasi hak akses pengguna lewat kueri database.

## Consequences

- **Positif:** Struktur database bersih, terdokumentasi di ERD, dan performa kueri optimal dengan database indexes.
- **Negatif:** Harus menulis kode mapping join model GORM di layer repository.
