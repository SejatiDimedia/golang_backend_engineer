# Database Design: Notification Service

---

## 1. Entity-Relationship Diagram (ERD)

```mermaid
erDiagram
    USERS {
        uint id PK
        varchar email UK "Not Null"
        text password_hash "Not Null"
        timestamp created_at
        timestamp updated_at
    }
    NOTIFICATIONS {
        uint id PK
        varchar type "Not Null"
        varchar target "Not Null"
        text content "Not Null"
        varchar status "Default: 'PENDING', Not Null"
        int max_retries "Default: 5"
        int attempt_count "Default: 0"
        timestamp send_at "Not Null"
        timestamp created_at
        timestamp updated_at
    }
    NOTIFICATION_LOGS {
        uint id PK
        uint notification_id FK "On Delete CASCADE"
        int attempt "Not Null"
        varchar status "Not Null"
        text error_message
        timestamp created_at
    }

    USERS ||--o{ NOTIFICATIONS : "creates"
    NOTIFICATIONS ||--o{ NOTIFICATION_LOGS : "logs"
```

## 2. Table schemas

### 1. `notifications`
Tabel utama pencatatan metadata notifikasi.
- **Constraints:**
  - `status`: Berisi `PENDING`, `PROCESSING`, `SENT`, atau `FAILED`.
  - `send_at`: Tipe data timestamp UTC untuk merekam waktu eksekusi.

### 2. `notification_logs`
Tabel audit log terperinci untuk merekam setiap percobaan pengiriman asinkron.
- **Constraints:**
  - `notification_id`: Foreign key ke `notifications.id` dengan konfigurasi `ON DELETE CASCADE`. Jika record notifikasi dihapus, data audit logs otomatis dibersihkan.

---

## 3. Database Indexes

Untuk mempercepat kueri pengecekan status notifikasi:

```sql
CREATE INDEX idx_notifications_status ON notifications (status);
```

**Justifikasi Indeks:**
Meskipun audit logs bertambah besar seiring waktu, status notifikasi yang aktif (terutama `PENDING` atau `PROCESSING`) selalu dapat dicari dengan cepat melalui indeks status tunggal ini.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi ERD skema audit trail notifikasi dan status indexes |
