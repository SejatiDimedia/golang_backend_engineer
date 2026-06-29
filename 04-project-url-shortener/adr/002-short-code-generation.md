# ADR-002: Strategi Pembuatan Short Code URL

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Layanan URL Shortener memerlukan kode unik yang pendek (short code) untuk mewakili URL target. Kode ini harus cukup pendek agar ramah pengguna, unik untuk menghindari tabrakan data, dan efisien saat dibuat pada volume traffic tinggi.

## Decision

Kami memutuskan untuk menggunakan **URL-safe Base64 encoding** yang dihasilkan dari **current timestamp** di sisi aplikasi Go sebagai strategi bawaan pembuatan short code.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **Base64 URL-Safe dari Timestamp (Chosen)** | - Proses pembuatan cepat di memori (independen dari DB).<br>- Risiko tabrakan sangat rendah untuk single-instance deployment. | - Panjang kode bisa sedikit bervariasi bergantung pada presisi timestamp (millisecond vs nanosecond).<br>- Data bersifat semi-sekuensial sehingga bisa dianalisis polanya. |
| **PostgreSQL Auto-Increment ID + Base62** | - Menghasilkan kode terpendek di awal proyek (misal ID 1 = "1").<br>- Keunikan dijamin oleh database. | - Harus melakukan *insert* ke DB terlebih dahulu untuk mendapatkan ID sebelum encoding.<br>- Mudah ditebak (sekuensial). |
| **UUID v4 (dipotong menjadi 6-8 karakter)** | - Acak penuh dan sulit ditebak. | - Risiko tabrakan (*hash collision*) meningkat drastis ketika string UUID dipotong menjadi hanya 8 karakter. |

## Reasoning

Menggunakan timestamp sebagai basis pengodean Base64 memungkinkan pembuatan kode secara *offline* di tingkat aplikasi (application-level generation). Hal ini menghilangkan ketergantungan round-trip ke database hanya untuk mendapatkan ID baris sebelum data disimpan. 

Untuk memastikan kode tersebut aman digunakan di dalam URL (karena Base64 standar menggunakan karakter `+` dan `/` yang bermasalah di browser), kami menggunakan varian **URL-safe Base64 Encoding** (`base64.URLEncoding` di Go) yang menggantikan karakter tersebut dengan `-` dan `_` serta membuang padding `=`.

## Consequences

- **Positif:** Kecepatan pembuatan tinggi, tidak membebani database dengan kueri tambahan untuk alokasi ID sekuensial.
- **Negatif:** Timestamp dengan resolusi rendah (misal detik) dapat menyebabkan tabrakan jika ada dua request masuk pada detik yang sama. Oleh karena itu, kami harus menggunakan resolusi mikrodetik (*microsecond*) atau nanodetik (*nanosecond*) plus salt acak tambahan jika perlu untuk mencegah tabrakan *concurrency*.

## Revisit conditions

Keputusan ini akan ditinjau kembali jika:
- Terjadi tabrakan kode pendek yang signifikan akibat peningkatan concurrency.
- Pengguna meminta agar short code harus sepenuhnya acak dan tidak dapat diprediksi (karena timestamp secara tidak langsung menunjukkan waktu pembuatan).
