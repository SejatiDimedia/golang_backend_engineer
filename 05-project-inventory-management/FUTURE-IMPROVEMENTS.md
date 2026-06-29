# Future Improvements: Inventory Management API

Dokumen ini mencatat batasan cakupan teknis yang kami batasi secara sengaja di Project 2 demi tujuan pembelajaran yang terfokus, serta memetakan rencana perbaikan di masa depan.

---

## Deferred by design (in scope for a future project)

| Item | Why deferred | Where it's actually addressed |
|---|---|---|
| **Manajemen Hak Akses Gudang (RBAC)** | URL mutasi stok saat ini bebas diakses tanpa login. Kami menunda autentikasi agar fokus pada transaksi SQL. | [Project 3 (Booking System)](../06-project-booking-system/) dan [Project 7 (Auth Service)](../10-project-auth-service/) |
| **Audit Logs Sistem** | Pencatatan mutasi stok hanya merekam kuantitas barang, bukan identitas operator yang melakukan aksi (karena tidak ada login). | [Project 4 (Digital Wallet API)](../07-project-digital-wallet/) |
| **Pemberitahuan Stok Menipis (Notifications)** | Sistem belum mengirim email atau notifikasi otomatis ketika stok produk berada di bawah batas minimum (*low stock alert*). | [Project 6 (Notification Service)](../09-project-notification-service/) |

## Deferred due to scope (not a future project's job, just not done)

| Item | Why deferred | Would require |
|---|---|---|
| **Database Migration Versioning** | Penggunaan `AutoMigrate` GORM sangat cepat saat memulai, tetapi berisiko tinggi saat deploy paralel. | Pengenalan alat migrasi seperti `golang-migrate` dan penulisan berkas DDL SQL up/down manual. |
| **Optimasi Batch Insert Mutasi** | Mutasi stok saat ini dilakukan satu per satu. Pembelian dalam jumlah besar dari supplier (bulk import) akan sangat lambat. | Penulisan method repository baru yang memanfaatkan fitur batch insert GORM (`CreateInBatches`). |
| **CSV Import** | Gudang biasanya memerlukan fitur import data master produk dari berkas CSV, bukan sekadar ekspor. | Penulisan parser CSV upload menggunakan multipart form data reader dan validasi data baris demi baris sebelum disimpan. |

## Known weaknesses worth revisiting

| Weakness | Risk if unaddressed | Candidate trigger to fix |
|---|---|---|
| **Pessimistic Locking Overhead** | Penggunaan `FOR UPDATE` saat Stock Out mengunci baris database produk. Jika produk yang sama dipotong stoknya oleh ratusan staff paralel, terjadi bottleneck kueri antrian (*locking queue*). | Saat traffic mutasi harian meningkat tajam atau saat kita melatih integrasi Redis mutex locking. |
| **Paginasi Offset Lambat** | Penggunaan kueri paginasi `LIMIT` dan `OFFSET` pada riwayat mutasi akan mengalami degradasi performa yang signifikan jika data mutasi mencapai jutaan baris. | Saat data log mutasi melebihi 100.000 entri, migrasikan ke kueri paginasi berbasis kursor (*keyset pagination*). |

## Ideas considered and explicitly rejected

| Idea | Why rejected |
|---|---|
| **Melakukan Update Stok Produk secara Asinkron (Goroutines)** | *Rejected.* Kami menolak memperbarui stok produk di background thread karena jika pencatatan histori mutasi sukses sementara update stok produk gagal (atau sebaliknya), data fisik barang menjadi tidak akurat. Integritas data mutasi mutlak harus dilindungi di bawah transaksi database sinkron (ACID). |

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen perbaikan masa depan untuk sistem inventaris |
