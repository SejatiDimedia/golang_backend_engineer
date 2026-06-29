# ADR-001: Pilihan Algoritma Penandatanganan Token

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Dalam arsitektur microservices terdistribusi, downstream services (seperti Wallet, Booking, dan Notification Service) harus memverifikasi integritas dan masa aktif token JWT yang dikirimkan oleh pengguna di setiap HTTP request. Jika verifikasi token memerlukan pemanggilan API sinkron ke Auth Service (introspeksi online) untuk setiap request, performa sistem akan terhambat dan Auth Service menjadi single point of failure (SPOF) dengan beban traffic tinggi.

## Decision

Kami memutuskan menggunakan **Asymmetric Encryption (RS256 - RSA Signature dengan SHA-256)** untuk penandatanganan token JWT.

Rancangan Key Management:
1. **Auth Service (Signer):** Memegang private key RSA 2048-bit (`private.key`) dan menggunakannya untuk menandatangani Access Token JWT saat login/refresh.
2. **Downstream Services (Verifier):** Hanya memegang public key RSA (`public.key`). Mereka memverifikasi tanda tangan JWT secara offline menggunakan public key tersebut tanpa perlu memanggil network API ke Auth Service.
3. Kunci RSA digenerate secara otomatis pada server booting pertama kali di folder `config/certs` jika file belum ada.

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Asymmetric RS256 (Chosen)** | - Downstream services dapat memvalidasi token secara offline secara instan ($O(1)$).<br>- Menghilangkan overhead network traffic verifikasi JWT.<br>- Mengurangi resiko kebocoran kunci, karena downstream service tidak perlu tahu private key. | - Penandatanganan token di CPU membutuhkan daya komputasi sedikit lebih tinggi dibanding HMAC.<br>- Harus mengelola siklus hidup key pair secara aman. |
| **B. Symmetric HS256** | - Proses verifikasi dan penandatanganan token di CPU sangat cepat. | - Seluruh service harus memiliki shared secret key yang sama. Jika satu service berhasil diretas, seluruh otentikasi di platform terancam bocor. |

## Reasoning

Opsi A (RS256) dipilih demi mendukung prinsip arsitektur microservices modern: modularitas dan skalabilitas independen. Dengan offline token verification, downstream services tetap dapat melayani otentikasi request meskipun Auth Service sedang mengalami downtime singkat. Selain itu, ini mengisolasi credential private key hanya di Auth Service saja.

## Consequences

- **Positif:** Latensi API di downstream service tetap rendah, keamanan platform tinggi.
- **Negatif:** Harus menambahkan folder `certs/` ke `.gitignore` agar private key tidak terpublikasikan ke GitHub.
