# Testing Report: AI Prompt Management API

---

## 1. Testing Strategy

Pengujian dilakukan untuk memvalidasi integritas pipeline kompilasi prompt, efisiensi caching, dan keandalan skema multi-tenancy:

1. **Regex Compiler Test (`internal/utils/compiler_test.go`):**
   - Menguji parser double curly braces `{{var}}` untuk berbagai format teks.
   - Memvalidasi parameter dinamis yang hilang otomatis dikosongkan.
   - Memeriksa akurasi formula token length estimation.
2. **API Key Caching Test (`internal/middleware/apikey_test.go`):**
   - Memanfaatkan `miniredis` in-memory redis.
   - Mensimulasikan cache miss (lookup ke database relasional, menyimpan hasil hash ke Redis).
   - Mensimulasikan cache hit (verifikasi API key instan via Redis).
   - Memverifikasi penolakan request token invalid.
3. **Offline JWT Verification Test (`internal/middleware/jwt_test.go`):**
   - Men-generate RSA private/public key dinamis temporer.
   - Menandatangani token JWT RS256 dan memverifikasinya secara offline via public key file.
4. **Relational Service Integration Test (`internal/service/prompt_test.go`):**
   - Menggunakan database **SQLite in-memory** untuk simulasi query relasional.
   - Menguji isolasi data per Workspace (multi-tenancy) dan otentikasi API Key terenkripsi.
   - Menguji asinkron worker log analytics daemon.

---

## 2. Test Execution Command

```bash
go test -v ./...
```

### Hasil Test Suites (PASS):
```
=== RUN   TestAPIKeyMiddleware
--- PASS: TestAPIKeyMiddleware (0.00s)
=== RUN   TestJWTMiddleware_OfflineVerification
--- PASS: TestJWTMiddleware_OfflineVerification (0.08s)
PASS
ok  	github.com/timurdian/prompt-management/internal/middleware	1.075s
=== RUN   TestPromptService_Workflow
--- PASS: TestPromptService_Workflow (0.12s)
PASS
ok  	github.com/timurdian/prompt-management/internal/service	1.356s
=== RUN   TestCompilePrompt
--- PASS: TestCompilePrompt (0.00s)
PASS
ok  	github.com/timurdian/prompt-management/internal/utils	0.453s
```
All tests passed successfully!
