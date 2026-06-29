# Setup Guide: AI Prompt Management API

---

## 1. Local Scaffolding & Run

1. **Inisiasi Environment:**
   ```bash
   cp .env.example .env
   ```

2. **Salin Public Key (PENTING):**
   Layanan ini membutuhkan public key dari Auth Service (Project 7) untuk validasi offline JWT token:
   ```bash
   mkdir -p certs
   cp ../10-project-auth-service/certs/public.key certs/public.key
   ```

3. **Nyalakan Database & Redis:**
   ```bash
   docker-compose up -d
   ```

4. **Jalankan Aplikasi:**
   ```bash
   go run cmd/server/main.go
   ```
   *Server akan berjalan di port `8082`.*

---

## 2. API Verification Walkthrough (cURL)

### 1. Dapatkan Token JWT RS256
Pertama, jalankan login/register di **Auth Service (Project 7)** pada port 8081 untuk memperoleh token JWT.
```bash
# Simpan JWT di env shell
export JWT_TOKEN="eyJhbGciOiJSUzI1NiIs..."
```

### 2. Membuat Workspace
```bash
curl -i -X POST http://localhost:8082/api/v1/workspaces \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Engineering Team"}'
```
*Respons mengembalikan data workspace ID (contoh: `1`).*

### 3. Generate API Key Workspace
```bash
curl -i -X POST http://localhost:8082/api/v1/workspaces/1/api-keys \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Live Server Key"}'
```
*Salin nilai `"api_key"` (mentah, diawali `prompt_live_...`). Simpan di env shell:*
```bash
export API_KEY="prompt_live_..."
```

### 4. Membuat Prompt Template & Versi
```bash
# 1. Create Prompt
curl -i -X POST http://localhost:8082/api/v1/prompts \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"workspace_id": 1, "name": "Translator", "description": "Translate to Spanish"}'

# 2. Create Version 1 (Draft)
curl -i -X POST http://localhost:8082/api/v1/prompts/1/versions \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"prompt_text": "Translate this sentence to Spanish: {{sentence}}"}'

# 3. Activate Version 1
curl -i -X PUT http://localhost:8082/api/v1/prompts/1/versions/1/activate \
  -H "Authorization: Bearer $JWT_TOKEN"
```

### 5. Kompilasi Prompt via API Key (Server-to-Server)
Sekarang gunakan API Key untuk memanggil prompt compiler:
```bash
curl -i -X POST http://localhost:8082/api/v1/client/prompts/1/compile \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"variables": {"sentence": "Welcome home, friend."}}'
```
*Respons mengembalikan teks terkompilasi: `"Translate this sentence to Spanish: Welcome home, friend."` beserta estimasi token.*
*Kueri validasi API Key ini berjalan $<2\text{ms}$ karena disangga Redis cache-aside.*

---

## 3. Jalankan Automated Tests
```bash
go test -v ./...
```
