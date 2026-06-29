# API Specification: Inventory Management API

**Base URL:** `http://localhost:8080`
**Auth:** `None`

---

## Conventions

- Semua parameter input dan output dikomunikasikan dalam bentuk format JSON.
- Ekspor data produk dikirim langsung sebagai stream file teks CSV (`text/csv`).
- Respons error default:
  ```json
  {
    "error": "Deskripsi kesalahan teknis atau bisnis"
  }
  ```

---

## Endpoints

### 1. Categories API

#### `POST /categories`
- **Description:** Membuat kategori baru.
- **Request:** `{"name": "Elektronik"}`
- **Response — 201 Created:** `{"id": 1, "name": "Elektronik", "created_at": "...", "updated_at": "..."}`

#### `GET /categories`
- **Description:** Mengambil daftar seluruh kategori.
- **Response — 200 OK:** `[{"id": 1, "name": "Elektronik"}]`

#### `GET /categories/:id`
- **Description:** Mengambil detail kategori berdasarkan ID.
- **Response — 200 OK:** `{"id": 1, "name": "Elektronik"}`

#### `PUT /categories/:id`
- **Description:** Mengubah nama kategori.
- **Request:** `{"name": "Gadget & Elektronik"}`
- **Response — 200 OK:** `{"id": 1, "name": "Gadget & Elektronik"}`

#### `DELETE /categories/:id`
- **Description:** Menghapus kategori. Menghasilkan error jika masih ada produk yang menggunakan kategori ini.
- **Response — 200 OK:** `{"message": "category deleted successfully"}`

---

### 2. Suppliers API

#### `POST /suppliers`
- **Request:** `{"name": "PT. Vendor Jaya", "contact_name": "Budi", "email": "sales@vendor.com", "phone": "0812"}`
- **Response — 201 Created:** `{"id": 1, "name": "PT. Vendor Jaya", ...}`

#### `GET /suppliers`
- **Response — 200 OK:** `[{"id": 1, "name": "PT. Vendor Jaya"}]`

#### `DELETE /suppliers/:id`
- **Description:** Menghapus supplier. Menghasilkan error jika masih ada produk yang merujuk supplier ini.
- **Response — 200 OK:** `{"message": "supplier deleted successfully"}`

---

### 3. Products API

#### `POST /products`
- **Request:**
  ```json
  {
    "sku": "PROD-LAP-001",
    "name": "Laptop Lenovo ThinkPad",
    "description": "Lenovo T490 Core i5",
    "price": 12500000.00,
    "category_id": 1,
    "supplier_id": 1
  }
  ```
- **Response — 201 Created:** Mengembalikan objek produk lengkap beserta relasi `category` dan `supplier` (preloaded), dengan `stock_quantity` bernilai `0`.

#### `GET /products`
- **Description:** Membaca daftar produk berpaginasi.
- **Query Params:**
  - `search` (opsional) - mencari produk berdasarkan nama/SKU.
  - `category_id` (opsional) - memfilter kategori produk.
  - `page` (opsional, default: `1`).
  - `limit` (opsional, default: `10`).
- **Response — 200 OK:**
  ```json
  {
    "data": [
      {
        "id": 1,
        "sku": "PROD-LAP-001",
        "name": "Laptop Lenovo ThinkPad",
        "price": 12500000.00,
        "stock_quantity": 0,
        "category_id": 1,
        "supplier_id": 1,
        "category": { "id": 1, "name": "Elektronik" },
        "supplier": { "id": 1, "name": "PT. Vendor Jaya" }
      }
    ],
    "meta": {
      "total": 1,
      "page": 1,
      "limit": 10
    }
  }
  ```

#### `GET /products/export`
- **Description:** Mengunduh berkas laporan CSV persediaan produk saat ini secara penuh (all products).
- **Response — 200 OK:**
  - Header: `Content-Type: text/csv`
  - Header: `Content-Disposition: attachment; filename=inventory_report.csv`
  - Body:
    ```csv
    ID,SKU,Name,Category,Supplier,Price,Stock Quantity
    1,PROD-LAP-001,Laptop Lenovo ThinkPad,Elektronik,PT. Vendor Jaya,12500000.00,15
    ```

---

### 4. Stock Mutations API

#### `POST /products/:id/stock-in`
- **Description:** Melakukan penambahan kuantitas stok barang masuk.
- **Request:**
  ```json
  {
    "quantity": 15,
    "reference": "PO-2026-001"
  }
  ```
- **Response — 200 OK:** Mengembalikan objek Stock Movement log.
  ```json
  {
    "id": 1,
    "product_id": 1,
    "type": "IN",
    "quantity": 15,
    "reference": "PO-2026-001",
    "created_at": "2026-06-29T13:30:00Z"
  }
  ```

#### `POST /products/:id/stock-out`
- **Description:** Melakukan pengurangan kuantitas stok barang keluar. Menghasilkan error jika kuantitas stok produk saat ini tidak mencukupi (stok kurang).
- **Request:**
  ```json
  {
    "quantity": 5,
    "reference": "SO-2026-099"
  }
  ```
- **Response — 200 OK:** Mengembalikan objek Stock Movement log.
  ```json
  {
    "id": 2,
    "product_id": 1,
    "type": "OUT",
    "quantity": 5,
    "reference": "SO-2026-099",
    "created_at": "2026-06-29T13:35:00Z"
  }
  ```
- **Response — 400 Bad Request:**
  ```json
  {
    "error": "insufficient stock quantity"
  }
  ```

#### `GET /stock-movements`
- **Description:** Melihat log histori seluruh mutasi pergerakan barang.
- **Query Params:**
  - `product_id` (opsional) - memfilter log produk tertentu.
  - `type` (opsional) - memfilter tipe mutasi (`IN` atau `OUT`).
  - `page`, `limit` (opsional).
- **Response — 200 OK:** Mengembalikan daftar log berpaginasi.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi spesifikasi REST API untuk Inventory API (master CRUD, mutasi stok, ekspor CSV) |
