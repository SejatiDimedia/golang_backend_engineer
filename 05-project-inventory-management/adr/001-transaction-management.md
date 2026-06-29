# ADR-001: Lokasi dan Pola Manajemen Transaksi Database

**Status:** `Accepted`
**Date:** 2026-06-29

---

## Context

Pada Project 2 (Inventory Management), beberapa operasi bisnis seperti mutasi stok (stok masuk/keluar) membutuhkan lebih dari satu kueri database yang harus dijalankan secara atomis (All-or-Nothing). Misalnya, saat melakukan Stock In, kami harus memperbarui kuantitas produk di tabel `products` dan menambahkan log histori mutasi di tabel `stock_movements`. Jika salah satu gagal, seluruh operasi harus di-rollback.

Kami memerlukan cara untuk mengelola transaksi database ini tanpa melanggar prinsip *Clean Architecture* (yaitu tidak membocorkan detail pustaka database seperti `*gorm.DB` atau `*sql.Tx` ke dalam Service Layer).

## Decision

Kami memutuskan untuk mengelola transaksi database menggunakan **Pola Transaction Manager berbasis Go Context** yang diimplementasikan di tingkat Repository/Infrastructure layer, sementara inisiasinya dipicu dari Service layer.

Implementasi akan menggunakan interface berikut:
```go
package repository

import "context"

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

## Alternatives considered

| Option | Pros | Cons |
|---|---|---|
| **A. Transaksi di Service Layer (Langsung)** | Sangat mudah dibaca secara langsung alurnya di kode service. | Mengimpor tipe data spesifik ORM (seperti `*gorm.DB`) ke Service Layer, merusak portabilitas Clean Architecture. |
| **B. Transaksi di Repository Wrapper (Callback via Context) [Chosen]** | - Service layer tetap bersih dari detail database relasional.<br>- Logika koordinasi transaksi tetap berada di level bisnis (Service).<br>- Transaksi dipropagasi secara implisit via context. | - Membutuhkan boilerplate tambahan untuk ekstraksi koneksi database dari context.<br>- Sedikit lebih sulit dipahami oleh pemula di awal. |
| **C. Transaksi di Repository Layer (Kueri Gabungan)** | Service layer tidak tahu ada transaksi. Cukup panggil `repo.RecordStockMovement(...)`. | Mengurangi fleksibilitas service layer jika ingin merangkai beberapa aksi repository yang bervariasi dalam satu transaksi. |

## Reasoning

Opsi B dipilih karena menjaga kesucian Service Layer. Logika bisnis di Service Layer hanya perlu memanggil:
```go
err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
    // 1. Update stok produk
    if err := s.productRepo.UpdateStock(txCtx, productID, quantity); err != nil {
        return err
    }
    // 2. Tambah log mutasi
    if err := s.movementRepo.Create(txCtx, movement); err != nil {
        return err
    }
    return nil
})
```
Repository Layer bertugas mendeteksi apakah `txCtx` membawa transaksi aktif GORM (melalui context key). Jika ada, repository menggunakan koneksi transaksi tersebut; jika tidak, ia menggunakan koneksi database standar.

Pola ini sangat tangguh dan meniru cara kerja manajemen transaksi deklaratif (seperti `@Transactional` di Spring Framework) namun diimplementasikan secara eksplisit dan idiomatik di Go.

## Consequences

- **Positif:** Dependensi GORM terisolasi penuh di tingkat repository. Pengujian unit (unit testing) pada service layer tetap mudah dilakukan dengan meniru (*mocking*) `TransactionManager` yang hanya menjalankan fungsi callback langsung tanpa interaksi database nyata.
- **Negatif:** Kerumitan kode infrastruktur bertambah karena kami harus mengimplementasikan fungsi utility context parser di tingkat database helper.

## Revisit conditions

Keputusan ini akan ditinjau kembali jika kami beralih dari GORM ke sqlx di proyek mendatang, karena mekanisme penyimpanan state transaksi di context mungkin memerlukan penyesuaian tipe data koneksi.
