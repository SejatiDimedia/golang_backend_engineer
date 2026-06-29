# Testing Report: Notification Service

---

## 1. Testing Strategy

Pengujian difokuskan pada keandalan logika antrean pesan asinkron Redis dan ketepatan perhitungan retry eksponensial:

1. **Queue Manager Test (`internal/queue/manager_test.go`):**
   - Menggunakan in-memory Redis mock (`miniredis/v2`) untuk mengisolasi unit testing dari Redis fisik.
   - Menguji keandalan blocking `Dequeue` (BRPOP) dan instant `Enqueue` (LPUSH).
   - Memverifikasi Lua script `MoveScheduledToReady` dengan cara fast-forward waktu di `miniredis` untuk melihat apakah scheduled task dipromosikan ke list instan tepat setelah jatuh tempo.
2. **Service & Backoff Test (`internal/service/notification_test.go`):**
   - Menggunakan repositori memori mock (`mockNotificationRepository`) untuk PostgreSQL.
   - Menguji apakah parameter `send_at` dinamis memicu masuk ke sorted set.
   - Memverifikasi perhitungan matematika exponential backoff $2^{\text{attempt}} \times 2$ detik.
3. **Middleware Test (`internal/middleware/auth_test.go`):**
   - Menguji parsing dan validasi JWT token (missing header, invalid token, dan success claim mapping).

---

## 2. Test Execution Command

Untuk mengeksekusi suite unit testing:
```bash
go test -v ./...
```

### Hasil Test Suites (PASS):
```
=== RUN   TestAuthMiddleware_MissingHeader
--- PASS: TestAuthMiddleware_MissingHeader (0.00s)
=== RUN   TestAuthMiddleware_InvalidToken
--- PASS: TestAuthMiddleware_InvalidToken (0.00s)
=== RUN   TestAuthMiddleware_Success
--- PASS: TestAuthMiddleware_Success (0.00s)
=== RUN   TestRedisQueueManager_EnqueueAndDequeue
--- PASS: TestRedisQueueManager_EnqueueAndDequeue (0.05s)
=== RUN   TestRedisQueueManager_MoveScheduledToReady
--- PASS: TestRedisQueueManager_MoveScheduledToReady (0.02s)
=== RUN   TestNotificationService_Create_Instant
--- PASS: TestNotificationService_Create_Instant (0.00s)
=== RUN   TestNotificationService_Create_Scheduled
--- PASS: TestNotificationService_Create_Scheduled (0.00s)
=== RUN   TestExponentialBackoffCalculation
--- PASS: TestExponentialBackoffCalculation (0.00s)
PASS
```

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Inisiasi dokumen laporan penulisan unit testing miniredis dan backoff math |
