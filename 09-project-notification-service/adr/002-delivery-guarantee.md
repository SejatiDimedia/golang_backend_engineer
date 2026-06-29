# ADR-002: Pilihan Jaminan Pengiriman Pesan

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Dalam perancangan layanan notifikasi, kegagalan jaringan temporer ke provider eksternal adalah hal yang pasti terjadi. Kami harus menentukan tingkat jaminan pengiriman (*delivery guarantees*) yang disediakan oleh Notification Service dan bagaimana cara memitigasi resiko pengiriman pesan ganda (spam) ke pengguna.

## Decision

Kami memutuskan mengimplementasikan jaminan pengiriman **At-Least-Once Delivery** dengan perlindungan status pengecekan di database relasional (idempotency status check).

Langkah penanganannya:
1. Notifikasi pasti dikirim *minimal satu kali*.
2. Jika provider eksternal mengembalikan error (atau timeout), tugas tersebut akan dimasukkan kembali ke antrean untuk dicoba lagi (*retry*).
3. Untuk mencegah duplikasi pengiriman jika kegagalan terjadi *setelah* provider memproses pengiriman namun sebelum respons diterima, kami melacak status eksklusif notifikasi (`status = 'SENT'`) di database relasional PostgreSQL. Sebelum worker mengirim pesan, ia wajib memeriksa status baris notifikasi tersebut secara atomik.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. At-Least-Once Delivery (Chosen)** | - Menjamin notifikasi tidak akan pernah hilang (pasti sampai ke pengguna).<br>- Logika retry asinkron sederhana dan tangguh terhadap downtime provider. | - Berpotensi mengirimkan pesan ganda jika respons sukses provider hilang di tengah jalan. |
| **B. Exactly-Once Delivery** | - Menghindari pesan ganda secara mutlak. | - Hampir tidak mungkin diimplementasikan tanpa dukungan koordinasi transaksi terdistribusi (2PC) dari pihak provider eksternal (SMTP/SMS gateways).<br>- Kompleksitas sistem melonjak drastis. |

## Reasoning

Di dunia nyata, provider SMS/Email eksternal tidak menyediakan jaminan *exactly-once*. Mereka hanya menerima payload dan mengirimkannya. Oleh karena itu, merancang server backend dengan jaminan Exactly-Once (Opsi B) adalah kesia-siaan arsitektur karena bottleneck tetap ada di sisi provider luar.

Opsi A (At-Least-Once) adalah standar industri terbaik untuk sistem notifikasi. Kami mengimbangi risiko pesan ganda dengan melakukan pengecekan status notifikasi secara ketat di PostgreSQL sebelum memicu pengiriman fisik, meminimalisir kemungkinan duplikasi akibat retry konkuren.

## Consequences

- **Positif:** Keandalan pengiriman notifikasi sangat tinggi, aman dari kegagalan jaringan temporer.
- **Negatif:** Klien harus mendesain aplikasinya untuk toleran terhadap kemungkinan kecil duplikasi pesan yang sangat jarang terjadi.
