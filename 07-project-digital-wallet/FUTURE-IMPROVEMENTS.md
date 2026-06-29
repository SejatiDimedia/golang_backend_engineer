# Future Improvements: Digital Wallet API

---

## Deferred by design (in scope for a future project)

| Item | Why deferred | Where it's actually addressed |
|---|---|---|
| **Shared Auth Service Retrofit** | Logika registrasi, login, dan middleware JWT saat ini ditulis lokal secara ad-hoc (duplikasi kode). | [Project 7 (Auth Service)](../10-project-auth-service/) |
| **Notification Integration** | Konfirmasi transaksi topup/transfer saat ini belum memicu notifikasi. | [Project 6 (Notification Service)](../09-project-notification-service/) |
| **Distributed Multi-Node Redlock** | Hand-rolled lock saat ini menggunakan single-node Redis command yang tidak aman jika Redis dikonfigurasi cluster aktif-aktif. | [Project 7 / Produksi Cloud Staging](../10-project-auth-service/) |

## Deferred due to scope (not a future project's job, just not done)

| Item | Why deferred | Would require |
|---|---|---|
| **Transaction Fee & Promo Models** | Pemotongan biaya admin transfer antar bank atau bonus promo cashback saat top-up. | Penambahan tabel kebijakan tarif (`fees` / `promos`) dan integrasi kalkulasi potongan di `WalletService.Transfer`. |
| **Interactive Recon Cron** | Layanan background job berkala (seperti cron job) untuk mencocokkan total balance ledger dengan running balance dompet dan memberi tanda jika selisih. | Implementasi task runner asinkron di Go dan kueri database `SUM(amount)` agregasi periodik. |
| **Daily Transaction Limits** | Batasan nominal akumulasi transfer harian (misal: maksimal Rp 10.000.000 per hari untuk keamanan). | Tabel limit profil user, dan kueri checking `SUM(amount)` transaksi transfer keluar pada hari tersebut sebelum kueri lock dimulai. |

## Known weaknesses worth revisiting

| Weakness | Risk if unaddressed | Candidate trigger to fix |
|---|---|---|
| **High Latency under Lock Retry** | Client transfer paralel yang diblokir oleh distributed lock melakukan retry dengan backoff konstan. Jika ratusan request paralel antri, thread aplikasi Go bisa terbuang sia-sia menunggu lock. | Ketika response time transfer konkuren tinggi terdeteksi di production metrics. |
| **Redis Memory Growth** | Cache idempotency key disimpan tanpa enkripsi body. Jika request body sangat besar, memori Redis membengkak. | Jika penggunaan ram Redis kontainer melebihi 1 GB. |
