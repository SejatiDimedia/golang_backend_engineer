# ADR-001: Pilihan Strategi Penyimpanan Berkas Fisik

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Aplikasi File Management memerlukan media penyimpanan fisik berkas yang diunggah oleh pengguna. Kami harus memilih arsitektur penyimpanan fisik yang efisien, aman, dan mempermudah skalabilitas server backend di masa depan.

## Decision

Kami memutuskan menggunakan **MinIO Object Storage** (S3-compatible API) untuk penyimpanan berkas fisik, dan melacak metadatanya di database relasional PostgreSQL.

Setiap berkas yang masuk:
1. Validasi metadata (tipe konten, ukuran) diproses di tingkat aplikasi Go.
2. Berkas fisik diunggah ke bucket privat MinIO menggunakan pustaka `minio-go`.
3. Informasi path/key unik MinIO beserta metadata lainnya disimpan di PostgreSQL.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. MinIO Object Storage (Chosen)** | - Membuat server backend bersifat *stateless* (bebas dijalankan multi-instance di bawah load balancer).<br>- Mendukung protokol standar S3 API (membuat migrasi ke AWS S3 atau Google Cloud Storage di cloud produksi menjadi sangat mudah tanpa merombak kode).<br>- Menawarkan fitur keamanan presigned URL bawaan. | - Memerlukan setup service tambahan (kontainer MinIO). |
| **B. Local Server Disk Storage** | - Setup sangat sederhana, cukup menulis berkas ke direktori `/uploads` lokal server. | - Menjadikan server bersifat *stateful* (jika server berskala 2 instance paralel, instance B tidak bisa membaca berkas yang diunggah ke instance A).<br>- Risiko kehilangan data tinggi jika instance server mati/di-rebuild. |

## Reasoning

Menyimpan berkas di disk lokal server (Opsi B) adalah kebiasaan buruk yang menyulitkan skalabilitas (*horizontal scaling*). Object Storage (Opsi A) adalah standar industri modern untuk menangani berkas tidak terstruktur. MinIO menyediakan server object storage kompatibel S3 yang dapat kami jalankan secara lokal menggunakan Docker Compose, memberikan pengalaman pengembangan mirip cloud produksi (AWS S3) tanpa biaya tambahan.

## Consequences

- **Positif:** Server stateless, migrasi cloud lancar, storage terisolasi rapi.
- **Negatif:** Menambah dependencies SDK `github.com/minio/minio-go/v7` pada aplikasi Go.
