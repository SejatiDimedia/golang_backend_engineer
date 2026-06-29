# API Documentation: AI Prompt Management API

---

## 1. Dashboard Admin Endpoints (JWT RS256 Protected)
*Seluruh endpoint di bawah ini membutuhkan header `Authorization: Bearer <JWT_TOKEN_PROJECT_7>`.*

### 1. Create Workspace
- **HTTP Method:** `POST`
- **Path:** `/api/v1/workspaces`
- **Request Body:**
  ```json
  {
    "name": "Generative AI Team"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 1,
    "name": "Generative AI Team",
    "created_at": "2026-06-29T19:40:00Z",
    "updated_at": "2026-06-29T19:40:00Z"
  }
  ```

### 2. Generate API Key
- **HTTP Method:** `POST`
- **Path:** `/api/v1/workspaces/:id/api-keys`
- **Request Body:**
  ```json
  {
    "name": "Production Server Key"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 5,
    "message": "API Key created successfully. Save it now, it won't be shown again.",
    "api_key": "prompt_live_a1b2c3d4e5f6...",
    "masked_key": "prompt_live_xxxx...e5f6",
    "expires_at": "2027-06-29T19:40:00Z"
  }
  ```

### 3. Create Prompt
- **HTTP Method:** `POST`
- **Path:** `/api/v1/prompts`
- **Request Body:**
  ```json
  {
    "workspace_id": 1,
    "name": "Sentiment Classifier",
    "description": "Classifies sentiment of user feedback"
  }
  ```

### 4. Create Prompt Version (Draft)
- **HTTP Method:** `POST`
- **Path:** `/api/v1/prompts/:id/versions`
- **Request Body:**
  ```json
  {
    "prompt_text": "Analyze the sentiment of this feedback as positive/negative: {{text}}"
  }
  ```

### 5. Activate Prompt Version
Menjadikan versi tertentu aktif sebagai acuan compiler.
- **HTTP Method:** `PUT`
- **Path:** `/api/v1/prompts/:id/versions/:version_number/activate`
- **Response (200 OK):**
  ```json
  {
    "message": "Prompt version activated successfully"
  }
  ```

### 6. Get Analytics Logs
- **HTTP Method:** `GET`
- **Path:** `/api/v1/workspaces/:id/analytics`

---

## 2. Client Compiler Endpoints (API Key Protected)
*Endpoint ini membutuhkan header `Authorization: Bearer <API_KEY_MINTAH>`.*

### 1. Compile Active Prompt
- **HTTP Method:** `POST`
- **Path:** `/api/v1/client/prompts/:id/compile`
- **Request Body:**
  ```json
  {
    "variables": {
      "text": "Antigravity coding is awesome!"
    }
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "compiled_prompt": "Analyze the sentiment of this feedback as positive/negative: Antigravity coding is awesome!",
    "token_estimate": 14
  }
  ```
