# Lessons Learned: Notification Service

Retrospektif pembelajaran dari arsitektur asinkron Notification Service.

---

## 1. Keuntungan Blocking Pop (`BRPOP`) Redis dibanding Polling
- **Temuan:** Di awal perancangan worker queue, alternatif polling berkala (`LPOP` setiap 50ms) dipertimbangkan untuk mengambil tugas. Namun, ini memicu pemborosan CPU cycle server Go dan membanjiri Redis dengan jutaan query kosong per menit saat sepi.
- **Pembelajaran:** Menggunakan perintah blocking pop `BRPOP` di Redis memecahkan masalah ini dengan elegan. Utas worker Go akan secara pasif tertidur (suspensi runtime scheduler) tanpa mengonsumsi resource, dan langsung dibangunkan secara instan oleh Redis saat tugas baru di-`LPUSH` ke antrean.

## 2. Pentingnya Lua Scripting untuk Atomic Scheduler
- **Temuan:** Saat scheduled task poller mengambil data jatuh tempo dari Redis Sorted Set, proses ini terdiri dari dua langkah: mengambil data (`ZRANGEBYSCORE`) lalu menghapusnya (`ZREMRANGEBYSCORE`). Di lingkungan horizontal scaling (banyak server pod berjalan bersamaan), dua pod poller dapat mengambil tugas yang sama secara bersamaan sebelum sempat dihapus, memicu duplikasi antrean (*double-enqueue*).
- **Pembelajaran:** Memindahkan seluruh logika polling (Range, Delete, LPUSH) ke dalam **Lua Script** menjamin atomisitas. Redis mengeksekusi Lua script secara *single-threaded*, memastikan tidak ada interupsi di tengah jalan dari poller node lain, menjaga keselamatan transaksi terdistribusi secara utuh.

## 3. Pentingnya Pengecekan Status DB Sebelum Eksekusi Task (Idempotensi)
- **Temuan:** Retry exponential backoff dapat mengirimkan notifikasi ganda jika respons provider hilang setelah sukses dikirim namun sebelum worker sempat mencatat sukses di database.
- **Pembelajaran:** Pengecekan status notifikasi (`Status == 'SENT'`) secara real-time di PostgreSQL tepat sebelum worker memanggil API provider eksternal adalah benteng pertahanan terbaik. Ini mencegah duplikasi eksekusi pada transaksi asinkron yang terlanjur terdistribusi.
