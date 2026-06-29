# Future Improvements: AI Prompt Management API

Rencana peningkatan fitur AI Prompt Management API untuk rilis mendatang.

---

## 1. SDK Client Library (Go / JS / Python)
- **Masalah Saat Ini:** Downstream services harus memanggil HTTP REST API manually ke `/compile` untuk mendapatkan prompt, yang menambah overhead latensi parsing payload HTTP.
- **Rencana Solusi:** Sediakan SDK Client bawaan yang mendukung local caching untuk prompt ACTIVE. SDK secara berkala mengunduh snapshot prompt via polling atau WebSocket, lalu melakukan kompilasi regex secara lokal di sisi client untuk latensi mutlak $0\text{ms}$.

## 2. API Key Rotation Policy
- **Masalah Saat Ini:** API Key yang dibuat tidak memiliki rotasi otomatis dan kadaluarsa statis 1 tahun.
- **Rencana Solusi:** Kembangkan fitur rotasi API Key aman dengan masa tenggang (*grace period*) di mana key lama dan key baru aktif bersamaan selama 24 jam untuk menghindari service disruption.

## 3. Dynamic Variables Schema Validation
- **Masalah Saat Ini:** Compiler tidak memvalidasi tipe parameter variabel yang dikirim client.
- **Rencana Solusi:** Integrasikan skema validasi tipe data (seperti JSON Schema) per prompt version. Contoh: variabel `age` wajib berupa integer, `email` wajib berformat email valid sebelum compiler mengganti placeholders.
