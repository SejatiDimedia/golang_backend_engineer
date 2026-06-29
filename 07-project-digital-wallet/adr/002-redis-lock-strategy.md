# ADR-002: Strategi Distributed Lock Berbasis Redis Hand-Rolled

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Untuk mencegah double-spending di dompet digital, kami harus memastikan bahwa tidak ada transaksi paralel yang memanipulasi saldo dari wallet yang sama secara bersamaan. Meskipun transaksi database PostgreSQL (`FOR UPDATE`) dapat mengamankan satu database, dalam lingkungan terdistribusi (di mana backend server berjalan di beberapa container Docker paralel), kami memerlukan sistem pengunci terdistribusi (*distributed locking*) sebelum query database dimulai untuk meminimalkan beban database.

Kami perlu memilih apakah menggunakan pustaka distributed lock standar (seperti `go-redsync/redsync`) atau menulis utility pengunci sendiri menggunakan Redis command primitive.

## Decision

Kami memutuskan untuk mengimplementasikan **Hand-Rolled Redis Distributed Lock** secara mandiri menggunakan redis client dasar.

Mekanisme locking hand-rolled:
1. **Acquire Lock:** Menggunakan perintah Redis `SET lock_key unique_token NX PX 10000` (Set value jika kunci belum ada, dengan masa kadaluarsa 10 detik). Token unik acak dihasilkan menggunakan UUID/Timestamp untuk setiap request lock.
2. **Release Lock:** Menggunakan skrip **Lua** di Redis untuk memastikan pelepasan lock aman secara atomik (hanya melepas lock jika token unik di Redis cocok dengan token unik milik proses tersebut). Hal ini mencegah proses A menghapus kunci milik proses B yang secara tidak sengaja ter-extend.

Skrip Lua Pelepasan Lock:
```lua
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
```

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Pustaka Redsync (Standard Library)** | - Mengimplementasikan Redlock algorithm resmi yang teruji di multi-node Redis cluster. | - Library blackbox yang menyembunyikan detail operasional Redis dasar. |
| **B. Hand-Rolled Redis Lock (Chosen)** | - Memberikan pembelajaran mendalam mengenai command Redis `SET NX PX` dan eksekusi Lua Script secara atomik.<br>- Mengurangi ketergantungan library pihak ketiga eksternal yang tidak diperlukan di single-node redis. | - Kurang aman jika diaplikasikan di multi-node Redis cluster tanpa sinkronisasi kompleks (namun sangat memadai untuk single-node Redis deployment staging/local). |

## Reasoning

Karena Project 4 ditujukan untuk melatih keahlian core backend dan integrasi Redis dasar, menulis sistem distributed locking sendiri (Opsi B) akan jauh lebih mendidik. Pengguna akan belajar mengapa token acak (*random value*) diperlukan saat acquire lock dan mengapa Lua script mutlak dibutuhkan saat release lock agar terhindar dari *race condition lock deletion*.

## Consequences

- **Positif:** Pengguna mengerti detail atomisitas Redis dan Lua script, kode utility terisolasi bersih di folder `internal/utils/lock.go`.
- **Negatif:** Harus menulis mekanisme percobaan ulang (*retry mechanism/backoff*) secara manual di tingkat kode Go jika lock gagal didapatkan pada percobaan pertama.

## Revisit conditions

Jika sistem ini bermigrasi ke arsitektur multi-node Redis Cluster terdistribusi yang aktif-aktif di produksi nyata, kami akan meninjau ulang dan mengganti hand-rolled lock ini dengan pustaka Redlock mapan seperti `go-redsync/redsync`.
