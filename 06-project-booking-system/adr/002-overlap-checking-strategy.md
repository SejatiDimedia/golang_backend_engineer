# ADR-002: Strategi Pencegahan Double-Booking (Overlap Checking)

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Layanan pemesanan coworking space harus mencegah bentrokan pemesanan (double-booking). Yaitu, tidak boleh ada dua pemesanan aktif untuk meja/ruangan kerja yang sama pada rentang waktu (`start_time` s/d `end_time`) yang saling bertumpang tindih (*overlap*).

Kami perlu memilih strategi pencegahan overlap yang aman dari kondisi balapan konkurensi (*concurrent race condition*) ketika dua pengguna memesan slot yang sama secara simultan.

## Decision

Kami memutuskan menggunakan **Pola Transaksional SELECT ... FOR UPDATE (Pessimistic Locking)** di tingkat aplikasi Go dan database relasional untuk melakukan pengecekan tumpang tindih slot secara sinkron.

Mekanisme logika overlap di Go:
Dua rentang waktu $[S_1, E_1]$ dan $[S_2, E_2]$ dinyatakan bentrok (overlap) jika dan hanya jika:
$$S_1 < E_2 \quad \text{dan} \quad S_2 < E_1$$

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Select Lock & Evaluasi di Go (Chosen)** | - Sangat portable (bisa berjalan di SQLite/MySQL/Postgres).<br>- Memudahkan unit testing dengan Mock database.<br>- Mengasah pemahaman perbandingan data waktu (`time.Time`) di Go. | - Sedikit lebih lambat karena evaluasi terjadi di memori aplikasi.<br>- Memerlukan query locking yang disiplin agar terhindar dari race condition. |
| **B. PostgreSQL Exclusion Constraint** | - Efisien dan super cepat karena dikelola langsung oleh kernel database.<br>- Keunikan data dijamin 100% oleh database engine. | - Mengikat ketat kode ke PostgreSQL (tidak bisa diuji menggunakan database SQLite in-memory saat testing).<br>- Syntax DDL PostgreSQL range index rumit dikelola di GORM. |

## Reasoning

Sebagai proyek di tingkatan Intermediate, salah satu tujuan pembelajaran utama adalah memahami penanganan datetime dan merancang algoritma perbandingan waktu di tingkat Go secara mandiri. Menggunakan Exclusion Constraint PostgreSQL (Opsi B) akan menyembunyikan proses ini secara magis di dalam database.

Selain itu, Opsi A membolehkan kami menggunakan database SQLite in-memory yang cepat saat menulis unit test otomatis untuk `BookingService` tanpa menghasilkan error sintaks kueri indeks Postgres.

Untuk mengamankan konkurensi di Opsi A, kami melakukan locking transaksional terhadap baris master data meja/ruangan yang dipesan (`SELECT ... FOR UPDATE` pada baris meja/ruangan terkait), yang memaksa transaksi pemesanan kedua untuk meja/ruangan tersebut mengantri hingga pengecekan transaksi pertama selesai.

## Consequences

- **Positif:** Logika bisnis terdokumentasi jelas di level kode Go, portabilitas database tetap terjaga untuk SQLite testing.
- **Negatif:** Menambah kueri database untuk melakukan *pessimistic lock* pada meja/ruangan sebelum memasukkan data booking baru.

## Revisit conditions

Keputusan ini akan dievaluasi jika volume pemesanan per detik (*write throughput*) meja/ruangan yang sama meningkat sangat tajam, di mana locking antrian pada baris meja/ruangan menjadi bottleneck utama. Di sana kita bisa meninjau opsi distributed lock berbasis Redis (Project 4).
