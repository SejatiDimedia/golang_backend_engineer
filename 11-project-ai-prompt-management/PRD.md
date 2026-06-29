# PRD: AI Prompt Management API

**Status:** `Approved`
**Author:** Antigravity (AI Pair Engineer) & Timur (Learner)
**Last updated:** 2026-06-29

---

## 1. Problem statement

Dalam pengembangan sistem bertenaga Kecerdasan Buatan (AI/LLM), pengelolaan instruksi (prompt engineering) yang ditaruh langsung secara statis di dalam baris kode (*hardcoded*) mempersulit iterasi. Perubahan teks instruksi kecil mengharuskan siklus build, testing, dan redeployment ulang seluruh aplikasi. Selain itu, kolaborasi tim prompt engineer dan tracking performa prompt di berbagai versi LLM sulit diukur. Sistem membutuhkan AI Prompt Management API terpusat untuk mengelola prompt template, melacak versi perubahan (versioning), mempartisi kepemilikan berdasarkan Workspace tim, serta menyediakan API Key aman untuk dikonsumsi secara real-time oleh microservices lain.

## 2. Goals

- **Multi-Tenant Workspace:** Menyediakan isolasi workspace bagi tim untuk berkolaborasi mengelola sekumpulan prompt.
- **Prompt Templating & Compiling:** Menyediakan parser teks prompt berbasis variabel (contoh: `Jawab pertanyaan ini: {{question}}` dengan compilations parameter dinamis).
- **Asymmetric Version Control:** Melacak riwayat versi prompt (v1, v2, v3) menggunakan model full-snapshot per versi untuk audit trail instan.
- **API Key Authentication:** Memungkinkan workspace men-generate API Key khusus (seperti `sk_live_...`) dengan enkripsi satu arah (SHA-256 hash) untuk otentikasi server-to-server.
- **Offline JWT Integration:** Mengonsumsi `public.key` dari Auth Service (Project 7) untuk mengamankan akses Dashboard admin web secara offline.
- **Usage Analytics Tracking:** Mencatat statistik konsumsi prompt (latency, total calls, hit count, dan estimasi token length) ke database relasional.

## 3. Non-goals

- **Direct LLM Gateways (OpenAI/Gemini Integration):** Layanan ini *tidak* bertugas melakukan pemanggilan fisik API ke model OpenAI/Gemini/Claude secara langsung. Fokusnya adalah mengembalikan prompt *terkompilasi* (teks yang variabelnya sudah diisi) kembali ke microservice client agar microservice client tersebut yang memicu pemanggilan LLM.
- **Diff Comparison Viewers:** Tidak membuat visual git diff engine di baris interface, cukup pencatatan metadata versi.

## 4. Target users / personas

| Persona | Need | Frequency of use |
|---|---|---|
| Prompt Engineer / AI Builder | Menulis, menguji template prompt, dan merilis versi baru (misal v2) untuk meningkatkan akurasi LLM. | Harian |
| Downstream App (e.g. Chatbot Service) | Mengonsumsi prompt terkompilasi secara asinkron via API Key saat runtime. | Ratusan kali per menit |

## 5. Functional requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | **Workspace Isolation:** User dapat membuat, melihat, dan mengundang anggota ke dalam Workspace. Seluruh data prompt terisolasi per Workspace (multi-tenancy). | Must |
| FR-2 | **Prompt Template Management:** CRUD prompt template beserta pencatatan variabel terikat. | Must |
| FR-3 | **Prompt Versioning:** Setiap modifikasi prompt template dapat disimpan sebagai versi baru (draft/active). Versi aktif yang dirilis tidak boleh diubah secara destruktif (*immutable active version*). | Must |
| FR-4 | **Prompt Compiler:** Endpoint `POST /prompts/:id/compile` menerima parameter variabel JSON (contoh: `{ "name": "Budi" }`) dan mengembalikan prompt terkompilasi. | Must |
| FR-5 | **API Key Management:** Admin dapat men-generate API key untuk Workspace. Kunci asli hanya ditampilkan satu kali (seperti Stripe API keys) dan disimpan secara hash (SHA-256) di database. | Must |
| FR-6 | **Offline Auth Integration:** Middleware dashboard admin memvalidasi token JWT pengguna secara offline menggunakan berkas `public.key` dari Auth Service (Project 7). | Must |
| FR-7 | **Usage Analytics Log:** Setiap pemanggilan API compiler mencatat latensi eksekusi, ID token, dan estimasi jumlah kata (token) ke database `prompt_analytics`. | Should |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Security | API Key luar menggunakan prefiks (`prompt_live_...`). Kunci asli diverifikasi dengan konversi hash SHA-256. |
| Performance | Proses kompilasi prompt via API Key harus memiliki throughput tinggi dengan latency $<15\text{ms}$. |
| offline auth | Downstream Booking/Wallet/Notification memverifikasi JWT dari Project 7 tanpa koneksi jaringan ke Auth Service. |

## 7. Constraints

- **Teknologi:** Go, PostgreSQL (v15), GORM, Gin, Docker & Compose.
- **Auth Key Dependency:** Membaca berkas `public.key` RSA dari Auth Service untuk memvalidasi otentikasi admin dashboard.

## 8. Success criteria

- User dapat masuk ke Dashboard menggunakan token JWT dari Auth Service (diverifikasi offline).
- Client aplikasi dapat memanggil REST API menggunakan API Key untuk mengambil/mengompilasi prompt template.
- Modifikasi prompt secara otomatis menaikkan nomor versi (v1 -> v2) dengan snapshot lengkap yang tercatat di DB.
- Grafik analitik pemakaian prompt terdata secara terperinci.

## 9. Open questions

- **Versioning Strategy (Deltas vs Snapshots):** Memilih **Full Snapshots**. Setiap perubahan pada prompt template diduplikasi secara utuh sebagai baris versi baru di database, guna mempercepat kueri baca runtime ($O(1)$) tanpa overhead rekonstruksi diff.
- **API Key Caching:** Memilih **Redis Cache**. Hasil hash API Key yang sukses diverifikasi disimpan di Redis dengan format `apikey:<hash>` beserta metadata workspace ID agar validasi runtime server-to-server berada di bawah $2\text{ms}$.


---

## Revision history

| Date | Change |
|---|---|
| 2026-06-29 | Draft awal dibuat oleh Antigravity |
