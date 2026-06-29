# ADR-002: Pilihan Strategi Akses Unduh Berkas

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Berkas yang disimpan di dalam MinIO bucket diatur secara privat untuk mencegah akses ilegal dari luar. Kami perlu menentukan bagaimana cara menyajikan tautan unduh/baca berkas yang aman kepada pengguna terotentikasi.

## Decision

Kami memutuskan menggunakan **S3 Presigned URL** sebagai metode utama untuk penyajian berkas, dan menyertakan endpoint **Direct Server Streaming** sebagai metode sekunder/fallback.

- **Presigned URL:** API `/files/:id/download` menghasilkan URL bertanda tangan S3 dengan masa aktif 15 menit. Klien mengunduh berkas langsung dari server MinIO menggunakan tautan tersebut.
- **Direct Server Streaming:** API `/files/:id/view` membaca byte stream dari MinIO via SDK dan mem-pipe datanya secara streaming ke Gin HTTP response writer menggunakan `io.Copy`.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. S3 Presigned URL (Chosen - Primary)** | - **Bandwidth & CPU Saver:** Server Go tidak terbebani proses transfer byte berkas berukuran besar (bandwidth di-offload ke MinIO).<br>- Keamanan tinggi karena tautan otomatis kedaluwarsa. | - Klien harus melakukan request kedua (mengakses url presigned) setelah memperoleh respons API. |
| **B. Direct Server Streaming (Chosen - Fallback)** | - Sangat bersih bagi klien, berkas langsung dibaca lewat satu panggilan API `/files/:id/view` terproteksi JWT. | - Membebani bandwidth dan memori RAM server Go jika ribuan user mendownload berkas besar secara paralel. |

## Reasoning

Untuk aplikasi skala besar, menyalurkan byte berkas besar (seperti PDF 10MB) melalui server aplikasi Go (Opsi B) akan memicu beban I/O (*bandwidth choking*) pada web server kita. 

Menggunakan Presigned URL (Opsi A) adalah pola arsitektur terbaik karena server Go hanya bertanggung jawab melakukan pemeriksaan otentikasi JWT (sangat cepat), memanggil SDK untuk membuat string tanda tangan presigned (operasi in-memory < 1ms), dan langsung mengembalikan URL tersebut ke klien. Beban download berkas fisik sepenuhnya ditransfer ke MinIO engine.

Kami tetap menyediakan Opsi B khusus untuk rendering media langsung (seperti menampilkan tag `<img src="...">` privat) yang tidak mendukung alur request ganda presigned URL.

## Consequences

- **Positif:** Menghemat performa server backend, aman, fleksibel.
- **Negatif:** Klien perlu menangani parsing URL presigned untuk memulai proses download.
