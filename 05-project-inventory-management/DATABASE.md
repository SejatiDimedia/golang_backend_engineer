# Database Design: Inventory Management API

**Status:** `Implemented`
**Engine:** PostgreSQL 15

---

## 1. Entity-relationship overview

Proyek ini memiliki skema relasional 4 tabel utama: `categories`, `suppliers`, `products`, dan `stock_movements`.

```
  ┌──────────────┐             ┌──────────────┐
  │  categories  │             │  suppliers   │
  └──────┬───────┘             └──────┬───────┘
         │ 1                          │ 1
         │                            │
         │ *                          │ *
   ┌─────▼────────────────────────────▼─────┐
   │                products                │
   └────────────────────┬───────────────────┘
                        │ 1
                        │
                        │ *
              ┌─────────▼─────────┐
              │  stock_movements  │
              └───────────────────┘
```

---

## 2. Schema

### Table: `categories`
| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | BIGINT | PRIMARY KEY, Auto-increment | |
| name | VARCHAR(100) | UNIQUE, NOT NULL | Nama kategori (misal: "Elektronik"). |
| created_at | TIMESTAMP | NOT NULL | |
| updated_at | TIMESTAMP | NOT NULL | |

### Table: `suppliers`
| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | BIGINT | PRIMARY KEY, Auto-increment | |
| name | VARCHAR(150) | NOT NULL | Nama vendor/supplier. |
| contact_name| VARCHAR(100) | | Nama perwakilan sales. |
| email | VARCHAR(100) | | |
| phone | VARCHAR(30) | | |
| created_at | TIMESTAMP | NOT NULL | |
| updated_at | TIMESTAMP | NOT NULL | |

### Table: `products`
| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | BIGINT | PRIMARY KEY, Auto-increment | |
| sku | VARCHAR(50) | UNIQUE, NOT NULL | Stock Keeping Unit (kode unik produk). |
| name | VARCHAR(150) | NOT NULL | Nama produk. |
| description | TEXT | | Deskripsi produk. |
| price | DECIMAL(12,2) | NOT NULL, CHECK (price > 0) | Harga produk. |
| stock_quantity| BIGINT | NOT NULL, DEFAULT 0 | Jumlah stok barang saat ini di gudang. |
| category_id | BIGINT | FOREIGN KEY, NOT NULL | Merujuk ke `categories.id` (RESTRICT on delete). |
| supplier_id | BIGINT | FOREIGN KEY, NOT NULL | Merujuk ke `suppliers.id` (RESTRICT on delete). |
| created_at | TIMESTAMP | NOT NULL | |
| updated_at | TIMESTAMP | NOT NULL | |

### Table: `stock_movements`
| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | BIGINT | PRIMARY KEY, Auto-increment | |
| product_id | BIGINT | FOREIGN KEY, NOT NULL | Merujuk ke `products.id` (CASCADE on delete). |
| type | VARCHAR(10) | NOT NULL | "IN" atau "OUT". |
| quantity | BIGINT | NOT NULL, CHECK (quantity > 0) | Jumlah barang mutasi. |
| reference | VARCHAR(100) | | Nomor faktur, PO, atau referensi penjualan. |
| created_at | TIMESTAMP | NOT NULL | Waktu mutasi terjadi. |

---

## 3. Relationships

| From Table | To Table | Relationship | On Delete | Reason |
|---|---|---|---|---|
| `products` | `categories` | Many-to-One | RESTRICT | Kategori tidak boleh dihapus jika masih ada produk yang menggunakan kategori tersebut. |
| `products` | `suppliers` | Many-to-One | RESTRICT | Supplier tidak boleh dihapus jika masih dirujuk oleh produk aktif. |
| `stock_movements`| `products` | Many-to-One | CASCADE | Jika data master produk dihapus, seluruh riwayat mutasi stok produk tersebut ikut terhapus. |

---

## 4. Indexes

| Table | Index Name | Columns | Type | Reason |
|---|---|---|---|---|
| `products` | `idx_products_sku` | `sku` | UNIQUE | Untuk pencarian cepat berdasarkan kode SKU produk. |
| `products` | `idx_products_category_id` | `category_id` | Standard | Mempercepat pencarian filter produk berdasarkan kategori. |
| `stock_movements`| `idx_stock_movements_product_id`| `product_id`| Standard | Mempercepat pembacaan riwayat pergerakan stok per produk. |

---

## 5. Transactions and consistency

Konsistensi data kuantitas stok barang (`stock_quantity`) dilindungi di bawah level isolasi default **Read Committed** PostgreSQL dengan menerapkan **Pessimistic Locking** (`SELECT ... FOR UPDATE`) saat proses pengurangan stok (Stock Out).

| Operation | Database Transaction Query Sequence | Reason for Locking |
|---|---|---|
| **Stock Out Mutation** | 1. `BEGIN;`<br>2. `SELECT stock_quantity FROM products WHERE id = 1 FOR UPDATE;`<br>3. Check if stock is sufficient. If yes:<br>4. `UPDATE products SET stock_quantity = stock_quantity - 5 WHERE id = 1;`<br>5. `INSERT INTO stock_movements ...;`<br>6. `COMMIT;` | Penggunaan `FOR UPDATE` mengunci baris produk tersebut di tingkat database. Jika ada request pengurangan stok simultan (paralel), request kedua akan dipaksa menunggu hingga request pertama melakukan `COMMIT`. Ini mencegah race condition di mana stok produk dikurangi melebihi persediaan. |

---

## 6. Migrations strategy

Sama seperti Project 1, kami memanfaatkan fitur **AutoMigration GORM** di berkas [main.go](file:///Users/timurdianradhasejati/Programming/Code/Golang/golang-backend-roadmap/05-project-inventory-management/cmd/server/main.go) dengan urutan pembuatan tabel: `Category` -> `Supplier` -> `Product` -> `StockMovement` untuk menjamin tidak ada kendala pembuatan foreign key constraint saat database diinisialisasi.

---

## 7. Sample queries

### Transactional Stock In Update
```sql
BEGIN;
UPDATE products SET stock_quantity = stock_quantity + 10, updated_at = NOW() WHERE id = 1;
INSERT INTO stock_movements (product_id, type, quantity, reference, created_at)
VALUES (1, 'IN', 10, 'PO-12345', NOW());
COMMIT;
```

### Paginated Search with Filters
```sql
SELECT p.*, c.name as category_name 
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.name ILIKE '%laptop%' AND p.category_id = 1
ORDER BY p.id DESC
LIMIT 10 OFFSET 0;
```

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi desain database relasional 4 tabel PostgreSQL 15 untuk inventory system |
