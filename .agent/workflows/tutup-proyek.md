# Workflow: Tutup Proyek (Definition of Done Check)

**Cara pakai:** `/tutup-proyek` setelah merasa fitur-fitur utama proyek sudah selesai dikerjakan.

## Tujuan
Cegah proyek "selesai" secara kode tapi dokumentasinya bolong — sesuai Definition of Done di `@00-learning-plan.md` §8.

## Langkah-langkah

1. Cek checklist berikut terhadap folder proyek yang sedang aktif, laporkan mana yang sudah/belum:
   - [ ] `README.md` terisi (bukan placeholder template)
   - [ ] `PRD.md` mencerminkan apa yang benar-benar dibangun (bukan rencana awal yang sudah berubah tanpa update)
   - [ ] `ROADMAP.md` proyek — semua fase ditandai status terkini
   - [ ] `ARCHITECTURE.md` cocok dengan struktur folder kode yang sebenarnya
   - [ ] `DATABASE.md` cocok dengan skema migrasi yang sebenarnya
   - [ ] `API.md` mencantumkan semua endpoint yang benar-benar ada
   - [ ] `SETUP.md` sudah dicoba ulang dari nol dan benar-benar berhasil
   - [ ] `DEPLOYMENT.md` terisi, minimal status "local only" dengan alasan
   - [ ] `TESTING.md` mencerminkan test yang benar-benar ada, `go test ./...` lulus
   - [ ] Minimal 1 ADR ada di `adr/`
   - [ ] `CHANGELOG.md` versi awal tercatat
   - [ ] `FUTURE-IMPROVEMENTS.md` terisi jujur

2. Bantu pengguna menulis `LESSONS-LEARNED.md` dengan menanyakan 4 pertanyaan dari `@00-learning-plan.md` §9 secara langsung ke pengguna — JANGAN menjawab sendiri atas nama pengguna. Ini refleksi pengguna, bukan ringkasan otomatis dari kamu.

3. Setelah semua lengkap, update tabel status proyek di root `@README.md` jadi "Documented".

4. Tanyakan ke pengguna: apakah ada proyek sebelumnya yang sekarang terasa perlu direvisit berdasarkan apa yang baru dipelajari? Kalau ya, jangan langsung eksekusi — catat sebagai rencana dulu, baru dieksekusi di sesi terpisah dengan persetujuan eksplisit.
