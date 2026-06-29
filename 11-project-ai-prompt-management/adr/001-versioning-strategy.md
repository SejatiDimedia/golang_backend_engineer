# ADR-001: Pilihan Strategi Versioning Prompt

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Layanan AI Prompt Management Service mengelola template prompt dinamis yang mengalami iterasi versi berulang oleh tim prompt engineer. Kami harus memutuskan struktur penyimpanan database untuk mencatat riwayat versi prompt (v1, v2, v3, dst.) guna memastikan performa kompilasi prompt optimal saat dipanggil oleh microservice client.

## Decision

Kami memutuskan menggunakan model **Full Snapshots per Versi** di dalam database relasional PostgreSQL.

Rincian implementasi:
1. Setiap kali prompt diubah atau versi baru dirilis, server menyimpan teks prompt secara utuh (*full snapshot*) sebagai baris record baru di tabel `prompt_versions`.
2. Versi aktif (status `ACTIVE`) bersifat *immutable* (tidak dapat diubah) untuk menjamin konsistensi performa aplikasi client. Jika ingin mengubah isi prompt, user wajib membuat draft versi baru (misal v2).

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Full Snapshots (Chosen)** | - Kecepatan baca / load prompt sangat cepat ($O(1)$) karena tidak membutuhkan kalkulasi reconstruct.<br>- Logika kode sederhana dan tangguh terhadap corruption data.<br>- Mempermudah audit trail riwayat versi secara langsung. | - Menggunakan ruang disk database sedikit lebih besar karena menduplikasi string teks berulang. |
| **B. Diffs/Deltas** | - Menghemat penyimpanan database dengan hanya menyimpan baris yang berubah (git-like). | - Mengambil prompt versi tertentu membutuhkan loop rekonstruksi overhead CPU yang menurunkan throughput kompilasi prompt di runtime. |

## Reasoning

Ukuran rata-rata satu berkas prompt instruksi AI/LLM sangat kecil (biasanya di bawah 5 KB). Oleh karena itu, penghematan memori dari penyimpanan diffs (Opsi B) tidak sebanding dengan overhead performa CPU dan kompleksitas kode yang ditimbulkannya saat merekonstruksi teks prompt di setiap HTTP request. Opsi A memberikan performa baca optimal yang sangat kritis untuk middleware microservices terdistribusi.

## Consequences

- **Positif:** Latensi kompilasi prompt sangat rendah, sistem tangguh dan mudah di-audit.
- **Negatif:** Konsumsi storage DB bertambah secara linier seiring pertambahan jumlah versi, namun tetap sangat kecil untuk kapasitas server modern.
