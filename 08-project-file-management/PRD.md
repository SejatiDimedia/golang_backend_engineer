# PRD: File Management Service

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Aplikasi web modern sering kali memerlukan fitur unggah berkas (seperti avatar pengguna, dokumen PDF, atau gambar lampiran). Menyimpan berkas fisik langsung di penyimpanan disk lokal server (*local server storage*) tidaklah scalable karena menyulitkan proses deploy multi-instance (load balanced servers). Sebaliknya, menyimpan berkas berukuran besar di database relasional dalam format BLOB memicu overhead kueri yang sangat lambat. Sistem memerlukan layanan terdedikasi yang memisahkan penyimpanan fisik berkas ke Object Storage (S3-compatible) sambil melacak metadatanya di database relasional secara aman.

## 2. Goals

- Memungkinkan pengguna mengunggah berkas menggunakan format HTTP Multipart Form Data.
- Mengintegrasikan penyimpanan berkas fisik ke **MinIO Object Storage** menggunakan S3 API standard.
- Melacak metadata berkas (ID, nama, ukuran, tipe MIME, dan object key MinIO) di database relasional PostgreSQL.
- Mengamankan hak akses berkas menggunakan **S3 Presigned URL** yang kedaluwarsa secara dinamis.
- Menyediakan endpoint streaming berkas langsung dari server backend untuk unduhan langsung.
- Membatasi ukuran unggahan berkas (maksimal 10 MB) dan menyaring tipe MIME tertentu (gambar & dokumen).

## 3. Non-goals

- **Client-Side Direct Upload:** Berkas wajib diunggah melalui backend server Go terlebih dahulu (bukan langsung dari browser client ke MinIO menggunakan presigned post url) agar metadata tervalidasi dengan aman.
- **File Versioning:** Modifikasi berkas dengan nama yang sama akan disimpan sebagai entri berkas baru (tidak ada pelacakan riwayat versi di object storage).

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Regular User | Mengunggah dokumen/gambar pribadi, melihat daftar berkas yang diunggah, dan mengunduh berkasnya secara aman. | Beberapa kali sehari |
| System Admin | Memantau total penggunaan storage fisik dan membersihkan berkas-berkas sampah yang yatim (*orphaned files*). | Mingguan |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | Pengguna dapat melakukan registrasi dan masuk (login JWT ad-hoc) sebelum mengunggah berkas. | Must |
| FR-2 | **Upload File:** Pengguna dapat mengunggah berkas gambar (JPEG/PNG) dan dokumen (PDF) dengan ukuran maksimal 10 MB melalui multipart form-data. | Must |
| FR-3 | **Metadata Tracking:** Server mencatat metadata berkas di database PostgreSQL, termasuk nama asli, ukuran berkas, tipe MIME, dan kunci path objek (`object_key`). | Must |
| FR-4 | **S3 Storage Integration:** Berkas fisik sukses terunggah ke bucket privat MinIO yang ditentukan. | Must |
| FR-5 | **Presigned URL Download:** Pengguna dapat meminta link unduh berkas yang diamankan menggunakan S3 Presigned URL (tautan kedaluwarsa dalam 15 menit). | Must |
| FR-6 | **Streaming Direct Download:** Server menyediakan alternatif unduhan streaming langsung (`GET /files/:id/view`) di mana server mem-pipe data byte stream dari MinIO ke respon HTTP. | Must |
| FR-7 | **Delete File:** Pengguna dapat menghapus berkasnya, yang otomatis menghapus metadata di PostgreSQL dan objek fisiknya di MinIO. | Must |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Storage Security | Bucket MinIO dikonfigurasikan secara privat. Berkas fisik tidak boleh diakses langsung tanpa otentikasi JWT atau presigned token yang sah. |
| Performance | Proses pengunduhan berkas direct stream wajib dikirim secara chunked streaming menggunakan `io.Copy` agar tidak memakan konsumsi memori RAM server backend. |
| Portability | Konfigurasi local stack menggunakan Docker Compose yang mencakup: Go server, PostgreSQL, dan MinIO console. |
| Robustness | Jika proses upload ke MinIO gagal, baris metadata berkas di PostgreSQL otomatis di-rollback atau dihapus untuk menghindari data yatim. |

## 7. Constraints

- **Teknologi:** Go, PostgreSQL, GORM, Gin, MinIO Go SDK (`github.com/minio/minio-go/v7`), Docker.
- **Ukuran File Maksimal:** 10 MB.
- **Allowed MIME Types:** `image/jpeg`, `image/png`, `application/pdf`.

## 8. Success criteria

- Berkas sukses terunggah dan tersimpan rapi di bucket privat MinIO.
- S3 Presigned URL sukses dihasilkan dan dapat diunduh langsung lewat browser sebelum masa kedaluwarsa habis (serta diblokir oleh MinIO setelah kedaluwarsa).
- Pengunggahan berkas >10 MB atau tipe konten terlarang (seperti `.exe` atau `.sh`) sukses ditolak di tingkat Gin validator.

## 9. Open questions

- **MinIO Bucket Auto-Creation:** Memilih **Rekomendasi (Ya)**. Server Go saat booting (`main.go`) otomatis mendeteksi apakah bucket target (contoh: `user-files`) sudah terdaftar di MinIO. Jika belum, server akan membuatnya secara otomatis menggunakan MinIO Go SDK.
- **Presigned URL Expiry Time:** Memilih **Opsi A (Konstan)**. Waktu kedaluwarsa tautan presigned diatur konstan 15 menit (melalui konfigurasi variabel environment `.env`). Hal ini menjaga interface API tetap sederhana dan terprediksi secara keamanan.

---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
