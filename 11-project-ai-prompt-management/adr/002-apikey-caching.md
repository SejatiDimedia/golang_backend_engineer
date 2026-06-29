# ADR-002: Pilihan Strategi Otentikasi API Key

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Aplikasi eksternal dan downstream microservices mengonsumsi prompt terkompilasi menggunakan otentikasi API Key (seperti `sk_live_...`) yang ditaruh pada request header. Karena otentikasi ini berjalan pada *setiap* HTTP API call, kami harus memilih strategi verifikasi API Key yang efisien tanpa membebani database relasional PostgreSQL.

## Decision

Kami memutuskan mengimplementasikan **Redis Caching untuk Otentikasi API Key**.

Langkah verifikasi:
1. API Key yang disimpan di database PostgreSQL berbentuk hash satu arah (SHA-256) demi keamanan.
2. Saat API request masuk, key mentah dari header di-hash dengan SHA-256.
3. Server mencari key hash tersebut di Redis cache (`apikey:<hash>`).
4. Jika cache hit: Ambil metadata workspace ID secara instan dan izinkan request.
5. Jika cache miss: Cari di PostgreSQL. Jika ditemukan dan valid, simpan hasilnya ke Redis cache dengan key `apikey:<hash>` dan TTL 1 jam.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Redis Caching (Chosen)** | - Latensi otentikasi sangat rendah ($<2\text{ms}$).<br>- Mengurangi beban query PostgreSQL secara signifikan. | - Harus mengelola invalidasi cache (menghapus key di Redis) saat admin men-delete/revoke API Key. |
| **B. Direct DB Query** | - State otentikasi selalu konsisten secara real-time. | - Lookup ke indeks PostgreSQL di setiap HTTP request dapat memperlambat throughput server-to-server. |

## Reasoning

Performa kompilasi prompt dinamis di runtime sangat bergantung pada efisiensi middleware otentikasi. Dengan meng-cache API Key yang aktif di Redis (Opsi A), kami memastikan overhead otentikasi berada di bawah 2ms, memenuhi target non-functional requirements throughput tinggi.

## Consequences

- **Positif:** Latensi HTTP request sangat rendah, platform tangguh dan skalabel.
- **Negatif:** Adanya sedikit kompleksitas kode untuk menghapus Redis cache saat API Key di-revoke secara manual oleh admin.
