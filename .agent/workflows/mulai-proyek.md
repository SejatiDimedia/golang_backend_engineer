# Workflow: Mulai Proyek Baru

**Cara pakai:** `/mulai-proyek` lalu sebutkan nomor proyek (contoh: "mulai proyek 1, URL Shortener")

## Tujuan
Memastikan proyek baru dimulai dengan urutan yang benar: konteks dulu, lalu PRD, lalu arsitektur, baru kode. Jangan pernah langsung menulis kode di langkah pertama.

## Langkah-langkah

1. **Baca konteks roadmap.**
   - Buka `@01-roadmap.md`, cari section proyek yang dimaksud (§3).
   - Buka `@00-learning-plan.md` untuk cek apakah Phase 0 (fundamentals) sudah dianggap selesai oleh pengguna. Kalau belum pernah dikonfirmasi, tanyakan ke pengguna.
   - Cek folder proyek tersebut (`0N-project-xxx/`) — kalau sudah ada isi sebagian, baca dulu sebelum menyarankan apa pun, jangan menimpa.

2. **Cek apakah PRD sudah ada.**
   - Kalau `PRD.md` di folder proyek itu belum ada atau kosong: STOP. Jangan lanjut ke kode.
   - Tawarkan ke pengguna: "Saya bantu isi PRD dulu dari template `@99-templates/PRD.template.md`?" Diskusikan goals, non-goals, functional requirements bersama pengguna — ini bukan kamu menulis sendirian, ini kolaborasi.

3. **Scaffold dokumen lain.**
   - Setelah PRD disepakati, copy seluruh template dari `99-templates/` ke folder proyek (lihat instruksi di `@99-templates/README.md`).
   - Isi `ROADMAP.md` (project-local) bersama pengguna: urutan fitur mana dibangun duluan.
   - Isi `ARCHITECTURE.md` minimal bagian struktur folder dan komponen, SEBELUM baris kode pertama ditulis.

4. **Baru mulai coding.**
   - Mulai dari fondasi: scaffold folder sesuai `ARCHITECTURE.md`, setup config, koneksi database, health check endpoint — sesuai urutan Phase 1 di `ROADMAP.template.md`.
   - Setiap keputusan teknis nontrivial (pilih library, pola desain) → buat ADR baru di `adr/` SEBELUM atau SEGERA SETELAH keputusan itu diambil, jangan ditunda sampai akhir proyek.

5. **Selalu update dokumen seiring kode berjalan.**
   - Endpoint baru → update `API.md` di saat yang sama.
   - Tabel/kolom baru → update `DATABASE.md` di saat yang sama.
   - Jangan biarkan dokumentasi jadi utang yang menumpuk di akhir.

## Yang TIDAK boleh dilakukan
- Jangan generate seluruh codebase proyek dalam satu langkah tanpa PRD/ARCHITECTURE disepakati dulu.
- Jangan menyamaratakan semua proyek — cek tingkat kesulitan proyek ini di `01-roadmap.md` §3 dan sesuaikan kedalaman testing/dokumentasi (lihat `00-learning-plan.md` §8).
